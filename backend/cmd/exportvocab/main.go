package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"obucon/internal/analysis"
	"obucon/internal/config"
	"obucon/internal/database"
)

const queryTimeout = 30 * time.Second

type flags struct {
	userID   uint
	language string
	out      string
}

func parseFlags(args []string) (flags, error) {
	fs := flag.NewFlagSet("exportvocab", flag.ContinueOnError)
	var (
		userID   uint
		language string
		out      string
	)
	fs.UintVar(&userID, "user", 0, "user id whose vocabulary to export (required, > 0)")
	fs.StringVar(&language, "language", "ja", "2-letter language code")
	fs.StringVar(&out, "out", "-", "output path, or '-' for stdout")
	if err := fs.Parse(args); err != nil {
		return flags{}, err
	}

	if userID == 0 {
		return flags{}, errors.New("-user must be provided and non-zero")
	}
	if len(language) != 2 {
		return flags{}, errors.New("-language must be a 2-character code")
	}

	return flags{userID: userID, language: language, out: out}, nil
}

func openOutput(path string) (io.WriteCloser, error) {
	if path == "" || path == "-" {
		return nopWriteCloser{os.Stdout}, nil
	}
	return os.Create(path)
}

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }

func writeCSV(w io.Writer, entries []analysis.VocabEntry) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{"lemma", "grade_level", "meaning"}); err != nil {
		return err
	}

	for _, e := range entries {
		grade := ""
		if e.GradeLevel != nil {
			grade = strconv.Itoa(*e.GradeLevel)
		}
		if err := cw.Write([]string{e.Lemma, grade, e.Meaning}); err != nil {
			return err
		}
	}

	cw.Flush()
	return cw.Error()
}

func main() {
	f, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "exportvocab:", err)
		os.Exit(2)
	}

	cfg := config.Load()
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("exportvocab: failed to open database: %v", err)
	}

	repo := analysis.NewRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	entries, err := repo.ListKnownWordsWithMeaning(ctx, f.userID, f.language)
	if err != nil {
		log.Fatalf("exportvocab: query failed: %v", err)
	}

	out, err := openOutput(f.out)
	if err != nil {
		log.Fatalf("exportvocab: cannot open output %q: %v", f.out, err)
	}
	defer out.Close()

	if err := writeCSV(out, entries); err != nil {
		log.Fatalf("exportvocab: write failed: %v", err)
	}

	fmt.Fprintf(os.Stderr, "exportvocab: wrote %d row(s) for user=%d language=%s\n", len(entries), f.userID, f.language)
}
