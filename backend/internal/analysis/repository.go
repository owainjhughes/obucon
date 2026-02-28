package analysis

import (
	"context"
	"fmt"
	"obucon/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	fmt.Print("Analysis NewRepository Function Reached\n")
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, analysis *models.Analysis) error {
	fmt.Print("Analysis Repository Create Function Reached\n")
	return r.db.WithContext(ctx).Create(analysis).Error
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*models.Analysis, error) {
	fmt.Print("Analysis Repository GetByID Function Reached\n")
	var analysis models.Analysis
	err := r.db.WithContext(ctx).First(&analysis, id).Error
	return &analysis, err
}

func (r *Repository) GetByUserID(ctx context.Context, userID uint) ([]models.Analysis, error) {
	fmt.Print("Analysis Repository GetByUserID Function Reached\n")
	var analyses []models.Analysis
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&analyses).Error
	return analyses, err
}

func (r *Repository) GetByTextHash(ctx context.Context, userID uint, textHash string) (*models.Analysis, error) {
	fmt.Print("Analysis Repository GetByTextHash Function Reached\n")
	var analysis models.Analysis
	err := r.db.WithContext(ctx).Where("user_id = ? AND text_hash = ?", userID, textHash).First(&analysis).Error
	return &analysis, err
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	fmt.Print("Analysis Repository Delete Function Reached\n")
	return r.db.WithContext(ctx).Delete(&models.Analysis{}, id).Error
}
