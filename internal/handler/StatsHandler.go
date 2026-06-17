package handler

import (
	"log"
	"net/http"
	"zzz/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHanlder struct {
	service *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHanlder {
	return &StatsHanlder{service: statsService}
}

func (h *StatsHanlder) GetStats(c *gin.Context) {
	resp, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		log.Println("StatsHanlder.GetStats: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
