package domain

import "context"

type Sentence struct {
	ID       int64   `json:"id"`
	Sentence string  `json:"sentence"`
	Category int64   `json:"category"`
	Cite     *string `json:"cite"`
	Author   *string `json:"author"`
}

type SentenceRepository interface {
	GetRandomSentence(ctx context.Context) (*Sentence, error)
}

type SentenceUsecase interface {
	GetRandomSentence(ctx context.Context) (*Sentence, error)
}
