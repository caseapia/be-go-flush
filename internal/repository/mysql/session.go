package mysql

import (
	"context"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/uptrace/bun"
)

func (r *Repository) CreateSession(ctx context.Context, tx bun.IDB, s *models.Session) error {
	_, err := tx.NewInsert().
		Model(s).
		Exec(ctx)
	return err
}

func (r *Repository) UpdateSession(ctx context.Context, tx bun.IDB, s *models.Session) error {
	_, err := tx.NewUpdate().
		Model(s).
		WherePK().
		Exec(ctx)
	return err
}

func (r *Repository) RevokeSession(ctx context.Context, tx bun.IDB, sessionID string) error {
	_, err := tx.NewUpdate().
		Model((*models.Session)(nil)).
		Where("id = ?", sessionID).
		Set("revoked = ?", true).
		Exec(ctx)
	return err
}

func (r *Repository) SearchSessionByRefreshHash(ctx context.Context, hash string) (*models.Session, error) {
	session := new(models.Session)

	err := r.DB.NewSelect().
		Model(session).
		Where("refresh_hash = ?", hash).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *Repository) SearchSessionsByUser(ctx context.Context, tx bun.IDB, userID uint64) ([]models.Session, error) {
	sessions := make([]models.Session, 0)

	err := tx.NewSelect().
		Model(&sessions).
		Where("user_id = ?", userID).
		Limit(10).
		Scan(ctx)
	if err != nil {
		return sessions, err
	}

	return sessions, err
}

func (r *Repository) SearchSessionByID(ctx context.Context, tx bun.IDB, sessionID string) (*models.Session, error) {
	session := new(models.Session)

	err := tx.NewSelect().
		Model(session).
		Where("id = ?", sessionID).
		Scan(ctx)
	if err != nil {
		return session, err
	}

	return session, err
}

func (r *Repository) CleanupExpiredSessions(ctx context.Context, tx bun.IDB, userID uint64) error {
	_, err := tx.NewDelete().
		Model((*models.Session)(nil)).
		Where("user_id = ?", userID).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) LockUser(ctx context.Context, tx bun.IDB, userID uint64) (bool, error) {
	var user models.User

	err := tx.NewSelect().
		Model(&user).
		Where("id = ?", userID).
		For("UPDATE").
		Scan(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) TerminateSession(ctx context.Context, tx bun.IDB, sessionID string) (bool, error) {
	_, err := tx.NewDelete().
		Where("id = ?", sessionID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, err
}

func (r *Repository) TerminateAllSessions(ctx context.Context, tx bun.IDB, userID uint64) (bool, error) {
	_, err := tx.NewDelete().
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, err
}
