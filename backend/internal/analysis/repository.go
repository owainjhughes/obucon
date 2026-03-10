package analysis

import (
	"context"
	"database/sql"
	"fmt"
	"obucon/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *Repository) UpsertKnownWord(ctx context.Context, knownWord *models.KnownWord) error {
	fmt.Print("Analysis Repository UpsertKnownWord Function Reached\n")
	if err := r.populateKnownWordGradeLevel(ctx, knownWord); err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "language"}, {Name: "lemma"}},
			DoUpdates: clause.AssignmentColumns([]string{"grade_level", "status", "metadata"}),
		}).
		Create(knownWord).Error
}

func (r *Repository) populateKnownWordGradeLevel(ctx context.Context, knownWord *models.KnownWord) error {
	if knownWord == nil || knownWord.GradeLevel != nil || knownWord.Lemma == "" {
		return nil
	}

	switch knownWord.Language {
	case "ja":
		var jlptLevel sql.NullInt32
		err := r.db.WithContext(ctx).
			Model(&models.JapaneseDictionary{}).
			Select("jlpt_level").
			Where("kanji = ? OR hiragana = ?", knownWord.Lemma, knownWord.Lemma).
			Order("jlpt_level ASC").
			Limit(1).
			Scan(&jlptLevel).Error
		if err != nil {
			return err
		}

		if jlptLevel.Valid {
			level := int(jlptLevel.Int32)
			knownWord.GradeLevel = &level
		}
	}

	return nil
}

func (r *Repository) ListKnownWordsByUser(ctx context.Context, userID uint, language string) ([]models.KnownWord, error) {
	fmt.Print("Analysis Repository ListKnownWordsByUser Function Reached\n")
	var knownWords []models.KnownWord
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if language != "" {
		query = query.Where("language = ?", language)
	}
	err := query.Order("created_at DESC").Find(&knownWords).Error
	return knownWords, err
}

func (r *Repository) DeleteKnownWordByID(ctx context.Context, userID, knownWordID uint) error {
	fmt.Print("Analysis Repository DeleteKnownWordByID Function Reached\n")
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.KnownWord{}, knownWordID).Error
}

func (r *Repository) IsKnownLemma(ctx context.Context, userID uint, language, lemma string) (bool, error) {
	fmt.Print("Analysis Repository IsKnownLemma Function Reached\n")
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.KnownWord{}).
		Where("user_id = ? AND language = ? AND lemma = ?", userID, language, lemma).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) GetKnownLemmas(ctx context.Context, userID uint, language string, lemmas []string) (map[string]bool, error) {
	fmt.Print("Analysis Repository GetKnownLemmas Function Reached\n")
	known := make(map[string]bool)
	if len(lemmas) == 0 {
		return known, nil
	}

	var rows []models.KnownWord
	err := r.db.WithContext(ctx).
		Select("lemma").
		Where("user_id = ? AND language = ? AND lemma IN ?", userID, language, lemmas).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		known[row.Lemma] = true
	}

	return known, nil
}

// VocabEntry represents a known word with optional dictionary meaning.
type VocabEntry struct {
	Lemma      string `json:"lemma"`
	GradeLevel *int   `json:"grade_level"`
	Meaning    string `json:"meaning"`
}

func (r *Repository) ListKnownWordsWithMeaning(ctx context.Context, userID uint, language string) ([]VocabEntry, error) {
	fmt.Print("Analysis Repository ListKnownWordsWithMeaning Function Reached\n")
	var entries []VocabEntry

	switch language {
	case "ja":
		err := r.db.WithContext(ctx).
			Table("known_words").
			Select("known_words.lemma, known_words.grade_level, coalesce(japanese_dictionary.meaning, '') AS meaning").
			Joins("LEFT JOIN japanese_dictionary ON known_words.lemma = japanese_dictionary.kanji OR known_words.lemma = japanese_dictionary.hiragana").
			Where("known_words.user_id = ? AND known_words.language = ?", userID, language).
			Order("known_words.created_at desc").
			Scan(&entries).Error
		if err != nil {
			return nil, err
		}
	default:
		// Fallback: just return known words without meaning
		var rows []models.KnownWord
		err := r.db.WithContext(ctx).
			Where("user_id = ? AND language = ?", userID, language).
			Order("created_at desc").
			Find(&rows).Error
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			entries = append(entries, VocabEntry{Lemma: row.Lemma, GradeLevel: row.GradeLevel})
		}
	}

	return entries, nil
}

func (r *Repository) BulkAddKnownWordsByJLPT(ctx context.Context, userID uint, language string, jlptLevel int) (int64, error) {
	fmt.Print("Analysis Repository BulkAddKnownWordsByJLPT Function Reached\n")

	// Insert JLPT words for the user, avoiding duplicates.
	// Uses PostgreSQL INSERT ... SELECT ... ON CONFLICT DO NOTHING.
	res := r.db.WithContext(ctx).Exec(
		`INSERT INTO known_words (user_id, language, lemma, grade_level, status, created_at)
		 SELECT ?, ?, japanese_dictionary.kanji, ?, 'known', CURRENT_TIMESTAMP
		 FROM japanese_dictionary
		 WHERE jlpt_level = ?
		 ON CONFLICT (user_id, language, lemma) DO NOTHING`,
		userID, language, jlptLevel, jlptLevel,
	)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, nil
}

func (r *Repository) GetDictionaryGradeLevels(ctx context.Context, language string, lemmas []string) (map[string]int, error) {
	fmt.Print("Analysis Repository GetDictionaryGradeLevels Function Reached\n")
	levels := make(map[string]int)
	if len(lemmas) == 0 {
		return levels, nil
	}

	switch language {
	case "ja":
		type dictionaryRow struct {
			Kanji     string `gorm:"column:kanji"`
			Hiragana  string `gorm:"column:hiragana"`
			JLPTLevel *int   `gorm:"column:jlpt_level"`
		}

		var rows []dictionaryRow
		err := r.db.WithContext(ctx).
			Model(&models.JapaneseDictionary{}).
			Select("kanji, hiragana, jlpt_level").
			Where("jlpt_level IS NOT NULL AND (kanji IN ? OR hiragana IN ?)", lemmas, lemmas).
			Find(&rows).Error
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			if row.JLPTLevel == nil {
				continue
			}

			level := *row.JLPTLevel

			if row.Kanji != "" {
				if existing, ok := levels[row.Kanji]; !ok || level < existing {
					levels[row.Kanji] = level
				}
			}

			if row.Hiragana != "" {
				if existing, ok := levels[row.Hiragana]; !ok || level < existing {
					levels[row.Hiragana] = level
				}
			}
		}
	}

	return levels, nil
}
