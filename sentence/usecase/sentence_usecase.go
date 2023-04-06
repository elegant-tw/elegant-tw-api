package usecase

import (
	"context"
	"elegant-tw-api/domain"

	"github.com/sirupsen/logrus"
)

type sentenceUsecase struct {
	sentenceRepo domain.SentenceRepository
}

func NewSentenceUsecase(sentenceRepo domain.SentenceRepository) domain.SentenceUsecase {
	return &sentenceUsecase{
		sentenceRepo: sentenceRepo,
	}
}

func (su *sentenceUsecase) GetRandomSentence(ctx context.Context) (*domain.Sentence, error) {
	aSentence, err := su.sentenceRepo.GetRandomSentence(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return aSentence, nil
}
