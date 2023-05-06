package http

import (
	"elegant-tw-api/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SentenceHandler struct {
	SentenceUsecase domain.SentenceUsecase
}

func NewSentenceHandler(e *gin.Engine, sentenceUsecase domain.SentenceUsecase) {
	handler := &SentenceHandler{
		SentenceUsecase: sentenceUsecase,
	}

	e.GET("/", handler.GetRandomSentence)
	e.GET("/all", handler.GetRandomSentenceWithToxic)
	e.GET("/test", handler.Benchmark)
}

func (s *SentenceHandler) Benchmark(c *gin.Context) {
	c.JSON(200, gin.H{
		"health": "OK",
	})
}

func (s *SentenceHandler) GetRandomSentence(c *gin.Context) {
	// category := c.Param("c")

	aSentence, err := s.SentenceUsecase.GetRandomSentence(c)

	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"status":  500,
			"message": "Internal error. Please try again later.",
		})
		return
	}

	c.JSON(200, aSentence)
}

func (s *SentenceHandler) GetRandomSentenceWithToxic(c *gin.Context) {
	// category := c.Param("c")

	aSentence, err := s.SentenceUsecase.GetRandomSentenceWithToxic(c)

	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"status":  500,
			"message": "Internal error. Please try again later.",
		})
		return
	}

	c.JSON(200, aSentence)
}
