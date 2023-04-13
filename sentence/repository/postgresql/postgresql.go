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
	row := p.db.QueryRow("SELECT id, sentence, category_id, cite, author FROM sentences")
	d := &domain.Sentence{}
	if err := row.Scan(&d.ID, &d.Sentence, &d.Category, &d.Cite, &d.Author); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return d, nil
}