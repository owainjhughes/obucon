package models

import (
	"time"
)

// Reference: https://gorm.io/docs/models.html
type User struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	Email        string      `gorm:"uniqueIndex;not null" json:"email"`
	Username     string      `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash string      `gorm:"not null" json:"-"` // Never expose in JSON
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	KnownWords   []KnownWord `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"known_words,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// KnownWord represents a user-scoped known word entry.
type KnownWord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Language   string    `gorm:"size:10;not null" json:"language"`
	Lemma      string    `gorm:"not null;index" json:"lemma"`
	GradeLevel *int      `gorm:"index" json:"grade_level"`
	Status     string    `gorm:"default:known" json:"status"`
	Metadata   []byte    `gorm:"type:jsonb" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
}

func (KnownWord) TableName() string {
	return "known_words"
}

// JapaneseDictionary represents an entry in the Japanese word dictionary
// This is for the PRE POPULATED non-user dictionary; where the data for a langage is stored
// Vocabulary words link to this for grade/difficulty lookup.
type JapaneseDictionary struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Kanji     string    `gorm:"not null;index" json:"kanji"`       // Kanji form (e.g., "日本") or hiragana if no kanji
	Hiragana  string    `gorm:"not null" json:"hiragana"`          // Hiragana reading (e.g., "にほん")
	Meaning   string    `gorm:"not null;type:text" json:"meaning"` // English translation
	JLPTLevel *int      `gorm:"index" json:"jlpt_level"`           // JLPT N5-N1 (nullable)
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for GORM
func (JapaneseDictionary) TableName() string {
	return "japanese_dictionary"
}

// Reference: https://gorm.io/docs/models.html
