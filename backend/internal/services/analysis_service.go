package services

import (
	"context"
	"fmt"
	"obucon/internal/nlp/japanese"
	"obucon/internal/repository"
)

type AnalysisService interface {
	AnalyzeText(ctx context.Context, userID uint, language, text string) (*AnalysisResult, error)
}

type AnalysisResult struct {
	Tokens      []japanese.Token `json:"tokens"`
	TotalTokens int              `json:"total_tokens"`
}

type analysisService struct {
	tokenizer    *japanese.Tokenizer
	analysisRepo repository.AnalysisRepository
	tokenRepo    repository.AnalysisTokenRepository
}

func NewAnalysisService(
	tokenizer *japanese.Tokenizer,
	analysisRepo repository.AnalysisRepository,
	tokenRepo repository.AnalysisTokenRepository,
) AnalysisService {
	return &analysisService{
		tokenizer:    tokenizer,
		analysisRepo: analysisRepo,
		tokenRepo:    tokenRepo,
	}
}

// TODO:
// - Dictionary lookup for word definitions
// - Vocabulary comparison for known/unknown words
// - Coverage percentage calculation
func (s *analysisService) AnalyzeText(ctx context.Context, userID uint, language, text string) (*AnalysisResult, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	tokens, err := s.tokenizer.Tokenize(text)
	if err != nil {
		return nil, fmt.Errorf("tokenization failed: %w", err)
	}

	// Create analysis record (for storage/history)
	// TODO: Later, calculate text_hash and store analysis metadata
	/*
		analysis := &models.Analysis{
			UserID:   userID,
			Language: language,
			TextHash: hashText(text),
		}
		if err := s.analysisRepo.Create(ctx, analysis); err != nil {
			return nil, fmt.Errorf("failed to create analysis record: %w", err)
		}
	*/

	// TODO: Store tokens in analysis_tokens table for history

	return &AnalysisResult{
		Tokens:      tokens,
		TotalTokens: len(tokens),
	}, nil
}
