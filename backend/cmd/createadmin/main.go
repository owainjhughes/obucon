package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"obucon/internal/auth"
	"obucon/internal/config"
	"obucon/internal/database"
)

const registerTimeout = 30 * time.Second

type flags struct {
	email    string
	username string
	password string
}

func parseFlags(args []string) (flags, error) {
	fs := flag.NewFlagSet("createadmin", flag.ContinueOnError)
	var f flags
	fs.StringVar(&f.email, "email", "", "email address for the new user (required)")
	fs.StringVar(&f.username, "username", "", "username for the new user (required)")
	fs.StringVar(&f.password, "password", "", "plaintext password for the new user (required)")
	if err := fs.Parse(args); err != nil {
		return flags{}, err
	}

	f.email = strings.TrimSpace(f.email)
	f.username = strings.TrimSpace(f.username)

	var missing []string
	if f.email == "" {
		missing = append(missing, "-email")
	}
	if f.username == "" {
		missing = append(missing, "-username")
	}
	if f.password == "" {
		missing = append(missing, "-password")
	}
	if len(missing) > 0 {
		return flags{}, fmt.Errorf("missing required flag(s): %s", strings.Join(missing, ", "))
	}

	if len(f.password) < 4 {
		return flags{}, errors.New("password must be at least 4 characters")
	}
	if len(f.username) < 3 || len(f.username) > 50 {
		return flags{}, errors.New("username must be between 3 and 50 characters")
	}
	if !strings.Contains(f.email, "@") {
		return flags{}, errors.New("email does not look valid")
	}

	return f, nil
}

func main() {
	f, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "createadmin:", err)
		os.Exit(2)
	}

	cfg := config.Load()
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("createadmin: failed to open database: %v", err)
	}

	userRepo := auth.NewRepository(db)
	authService := auth.NewService(userRepo, cfg.JWTSecret)

	ctx, cancel := context.WithTimeout(context.Background(), registerTimeout)
	defer cancel()

	user, err := authService.Register(ctx, f.email, f.username, f.password)
	if err != nil {
		log.Fatalf("createadmin: registration failed: %v", err)
	}

	fmt.Printf("created user id=%d email=%s username=%s\n", user.ID, user.Email, user.Username)
}
