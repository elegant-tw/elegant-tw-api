package usecase_test

import (
	"elegant-tw-api/domain"
	mocks "elegant-tw-api/mocks/domain"
	ucase "elegant-tw-api/sentence/usecase"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestGetRandomSentence(t *testing.T) {
	mockSentenceRepo := new(mocks.SentenceRepository)
	mockSentence := domain.Sentence{
		ID:       0,
		Sentence: "Go Go Go",
		Category: 0,
		Cite:     "CSGO",
		Author:   "No one",
	}

	t.Run("Success", func(t *testing.T) {
		mockSentenceRepo.
			On("GetRandomSentence", mock.Anything).
			Return(&mockSentence, nil).Once()

		u := ucase.NewSentenceUsecase(mockSentenceRepo)
		aSentence, err := u.GetRandomSentence(context.TODO())

		assert.NoError(t, err)
		assert.NotNil(t, aSentence)

		mockSentenceRepo.AssertExpectations(t)
	})
	t.Run("Fail", func(t *testing.T) {
		mockSentenceRepo.
			On("GetRandomSentence", mock.Anything).
			Return(nil, errors.New("Get error")).Once()

		u := ucase.NewSentenceUsecase(mockSentenceRepo)
		aSentence, err := u.GetRandomSentence(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, aSentence)

		mockSentenceRepo.AssertExpectations(t)
	})
}
