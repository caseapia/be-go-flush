package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
)

func (r *Repository) CreateSession(ctx context.Context, s *models.Session) error {
	_, err := r.db.NewInsert().
		Model(s).
		Exec(ctx)
	return err
}

func (r *Repository) UpdateSession(ctx context.Context, s *models.Session) error {
	_, err := r.db.NewUpdate().
		Model(s).
		WherePK().
		Exec(ctx)
	return err
}

func (r *Repository) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := r.db.NewUpdate().
		Model((*models.Session)(nil)).
		Set("revoked = ?", true).
		Where("id = ?", sessionID).
		Exec(ctx)
	return err
}

func (r *Repository) GetSessionByHash(ctx context.Context, hash string) (*models.Session, error) {
	session := new(models.Session)
	err := r.db.NewSelect().
		Model(session).
		Where("refresh_hash = ?", hash).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return session, nil
}
