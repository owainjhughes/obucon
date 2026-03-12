package analysis

import (
	"net/http"
	"obucon/internal/helpers"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	analysisService *Service
}

func NewAnalysisHandler(analysisService *Service) *AnalysisHandler {
	return &AnalysisHandler{analysisService: analysisService}
}

type AnalyzeRequest struct {
	Text     string `json:"text" binding:"required"`
	Language string `json:"language" binding:"required,len=2"`
}

type addKnownWordRequest struct {
	Lemma    string `json:"lemma" binding:"required"`
	Language string `json:"language" binding:"required,len=2"`
}

type bulkVocabRequest struct {
	JLPTLevel string `json:"jlpt_level" binding:"required,oneof=N5 N4 N3 N2 N1"`
	Language  string `json:"language" binding:"required,len=2"`
}

func (h *AnalysisHandler) AnalyzeText(c *gin.Context) {
	userID, ok := helpers.UserIDFromContext(c)
	if !ok {
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
	userID, ok := helpers.UserIDFromContext(c)
	if !ok {
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
	userID, ok := helpers.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req bulkVocabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	levelMap := map[string]int{"N5": 5, "N4": 4, "N3": 3, "N2": 2, "N1": 1}
	jlptLevel := levelMap[req.JLPTLevel]

	count, err := h.analysisService.AddBulkKnownVocabulary(c.Request.Context(), userID, req.Language, jlptLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"added": count})
}

func (h *AnalysisHandler) AddKnownWord(c *gin.Context) {
	userID, ok := helpers.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req addKnownWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.analysisService.AddKnownWord(c.Request.Context(), userID, req.Language, req.Lemma)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
