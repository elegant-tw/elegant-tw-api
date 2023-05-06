package postgresql

import (
	"context"
	"database/sql"
	"elegant-tw-api/domain"

	"github.com/sirupsen/logrus"
)

type postgresqlSentenceRepository struct {
	db *sql.DB
}

func NewpostgresqlSentenceRepoistory(db *sql.DB) domain.SentenceRepository {
	return &postgresqlSentenceRepository{db}
}

func (p *postgresqlSentenceRepository) GetRandomSentence(ctx context.Context) (*domain.Sentence, error) {
	row := p.db.QueryRow(
		`SELECT sentences.id, sentences.sentence, sentences.category_id, sentences.cite, sentences.author
		FROM sentences, categories
		WHERE categories.is_toxic = False AND sentences.category_id = categories.category_id
		ORDER BY RANDOM()
		LIMIT 1;`,
	)
	d := &domain.Sentence{}
	if err := row.Scan(&d.ID, &d.Sentence, &d.Category, &d.Cite, &d.Author); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return d, nil
}

func (p *postgresqlSentenceRepository) GetRandomSentenceWithToxic(ctx context.Context) (*domain.Sentence, error) {
	row := p.db.QueryRow(
		`SELECT id, sentence, category_id, cite, author
		FROM sentences
		ORDER BY RANDOM()
		LIMIT 1;`,
	)
	d := &domain.Sentence{}
	if err := row.Scan(&d.ID, &d.Sentence, &d.Category, &d.Cite, &d.Author); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return d, nil
}
