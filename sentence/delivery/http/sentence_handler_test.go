package http_test

import (
	"elegant-tw-api/domain"
	mocks "elegant-tw-api/mocks/domain"
	sentenceHandlerHttpDelivery "elegant-tw-api/sentence/delivery/http"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRandomSentence(t *testing.T) {
	mockSentence := domain.Sentence{
		ID:       0,
		Sentence: "Go Go Go",
		Category: 0,
		Cite:     "CSGO",
		Author:   "No one",
	}
	mockSentenceMarshal, _ := json.Marshal(mockSentence)
	mockSentenceCase := new(mocks.SentenceUsecase)

	mockSentenceCase.On("GetRandomSentence", mock.Anything).Return(&mockSentence, nil)

	r := gin.Default()

	handler := sentenceHandlerHttpDelivery.SentenceHandler{
		SentenceUsecase: mockSentenceCase,
	}

	r.GET("/", handler.GetRandomSentence)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(mockSentenceMarshal), w.Body.String())
}
