package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sonkiee/ehhh-api/internal/domain"
	"github.com/sonkiee/ehhh-api/internal/service"
)

type DilemmaHandler struct {
	svc *service.DilemmaService
}

func NewDilemmaHandler(svc *service.DilemmaService) *DilemmaHandler {
	return &DilemmaHandler{svc: svc}
}

type createReq struct {
	Question string `json:"question"`
	OptionA  string `json:"optionA"`
	OptionB  string `json:"optionB"`
}

func (h *DilemmaHandler) Create(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	d, err := h.svc.Create(c.Request.Context(), req.Question, req.OptionA, req.OptionB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *DilemmaHandler) List(c *gin.Context) {
	limit := 20
	offset := 0

	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	items, err := h.svc.List(c.Request.Context(), limit, offset) // you’ll add List() in service + repo
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "limit": limit, "offset": offset})
}

type voteReq struct {
	Choice string `json:"choice"` // "A" or "B"
}

func (h *DilemmaHandler) Vote(c *gin.Context) {
	id := c.Param("id")

	var req voteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	var choice domain.VoteChoice
	switch req.Choice {
	case "A":
		choice = domain.VoteA
	case "B":
		choice = domain.VoteB
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "choice must be A or B"})
		return
	}

	a, b, err := h.svc.Vote(c.Request.Context(), id, choice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dilemmaId": id, "countA": a, "countB": b})
}

func (h *DilemmaHandler) Get(c *gin.Context) {
	id := c.Param("id")

	d, a, b, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dilemma": d,
		"votes": gin.H{
			"A": a,
			"B": b,
		},
	})
}
