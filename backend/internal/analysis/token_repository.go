package analysis

import (
	"context"
	"fmt"
	"obucon/internal/models"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	fmt.Print("Analysis NewTokenRepository Function Reached\n")
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, token *models.AnalysisToken) error {
	fmt.Print("Analysis TokenRepository Create Function Reached\n")
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *TokenRepository) BulkCreate(ctx context.Context, tokens []models.AnalysisToken) error {
	fmt.Print("Analysis TokenRepository BulkCreate Function Reached\n")
	return r.db.WithContext(ctx).CreateInBatches(tokens, 100).Error
}

func (r *TokenRepository) GetByAnalysisID(ctx context.Context, analysisID uint) ([]models.AnalysisToken, error) {
	fmt.Print("Analysis TokenRepository GetByAnalysisID Function Reached\n")
	var tokens []models.AnalysisToken
	err := r.db.WithContext(ctx).Where("analysis_id = ?", analysisID).Find(&tokens).Error
	return tokens, err
}

func (r *TokenRepository) GetUnknownByAnalysisID(ctx context.Context, analysisID uint) ([]models.AnalysisToken, error) {
	fmt.Print("Analysis TokenRepository GetUnknownByAnalysisID Function Reached\n")
	var tokens []models.AnalysisToken
	err := r.db.WithContext(ctx).Where("analysis_id = ? AND is_known = ?", analysisID, false).Order("grade_level ASC").Find(&tokens).Error
	return tokens, err
}

func (r *TokenRepository) Delete(ctx context.Context, id uint) error {
	fmt.Print("Analysis TokenRepository Delete Function Reached\n")
	return r.db.WithContext(ctx).Delete(&models.AnalysisToken{}, id).Error
}
