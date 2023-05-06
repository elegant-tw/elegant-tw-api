package postgresql_test

import (
	"context"
	"elegant-tw-api/domain"
	"testing"

	sentencePostgresqlRepo "elegant-tw-api/sentence/repository/postgresql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/magiconair/properties/assert"
)

func TestGetRandomSentence(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockCite := "Counter-Strike: Global Offensive"
	mockSentence := &domain.Sentence{
		ID:       0,
		Sentence: "Go! Go! Go!",
		Category: 0,
		Cite:     &mockCite,
		Author:   nil,
	}

	rows := sqlmock.NewRows([]string{"id", "sentence", "category_id", "cite", "author"}).
		AddRow(mockSentence.ID, mockSentence.Sentence, mockSentence.Category, mockSentence.Cite, mockSentence.Author)

	query := `SELECT sentences.id, sentences.sentence, sentences.category_id, sentences.cite, sentences.author
	FROM sentences, categories
	WHERE categories.is_toxic = False AND sentences.category_id = categories.category_id`

	mock.ExpectQuery(query).WillReturnRows(rows)
	d := sentencePostgresqlRepo.NewpostgresqlSentenceRepoistory(db)
	aSentence, _ := d.GetRandomSentence(context.TODO())
	assert.Equal(t, mockSentence, aSentence)
}

func TestGetRandomSentenceWithToxic(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockCite := "Counter-Strike: Global Offensive"
	mockSentence := &domain.Sentence{
		ID:       0,
		Sentence: "Go! Go! Go!",
		Category: 0,
		Cite:     &mockCite,
		Author:   nil,
	}

	rows := sqlmock.NewRows([]string{"id", "sentence", "category_id", "cite", "author"}).
		AddRow(mockSentence.ID, mockSentence.Sentence, mockSentence.Category, mockSentence.Cite, mockSentence.Author)

	query := `SELECT id, sentence, category_id, cite, author FROM sentences`

	mock.ExpectQuery(query).WillReturnRows(rows)
	d := sentencePostgresqlRepo.NewpostgresqlSentenceRepoistory(db)
	aSentence, _ := d.GetRandomSentenceWithToxic(context.TODO())
	assert.Equal(t, mockSentence, aSentence)
}
