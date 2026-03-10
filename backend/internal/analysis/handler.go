package analysis

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	analysisService *Service
}

func NewAnalysisHandler(analysisService *Service) *AnalysisHandler {
	fmt.Print("Analysis NewAnalysisHandler Function Reached\n")
	return &AnalysisHandler{analysisService: analysisService}
}

type AnalyzeRequest struct {
	Text     string `json:"text" binding:"required"`
	Language string `json:"language" binding:"required,len=2"`
}

func (h *AnalysisHandler) AnalyzeText(c *gin.Context) {
	fmt.Print("Analysis Handler AnalyzeText Function Reached\n")

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userID uint
	switch v := userIDInterface.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.analysisService.AnalyzeText(c.Request.Context(), userID, req.Language, req.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AnalysisHandler) ListVocabulary(c *gin.Context) {
	fmt.Print("Analysis Handler ListVocabulary Function Reached\n")

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userID uint
	switch v := userIDInterface.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	language := c.DefaultQuery("language", "ja")

	vocab, err := h.analysisService.ListKnownVocabulary(c.Request.Context(), userID, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vocab": vocab})
}

func (h *AnalysisHandler) BulkAddVocabulary(c *gin.Context) {
	fmt.Print("Analysis Handler BulkAddVocabulary Function Reached\n")

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userID uint
	switch v := userIDInterface.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	type BulkRequest struct {
		JLPTLevel string `json:"jlpt_level" binding:"required,oneof=N5 N4 N3 N2 N1"`
		Language  string `json:"language" binding:"required,len=2"`
	}

	var req BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	levelMap := map[string]int{"N5": 5, "N4": 4, "N3": 3, "N2": 2, "N1": 1}
	jlptLevel, ok := levelMap[req.JLPTLevel]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid jlpt level"})
		return
	}

	count, err := h.analysisService.AddBulkKnownVocabulary(c.Request.Context(), userID, req.Language, jlptLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"added": count})
}
