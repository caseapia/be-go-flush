package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/pkg/utils/models/enums"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (r *Repository) populateUserBadges(ctx context.Context, u *models.User) error {
	if len(*u.BadgeIDs) == 0 {
		u.Badges = make([]models.Badge, 0)
		return nil
	}

	err := r.DB.NewSelect().
		Model(&u.Badges).
		Where("id IN (?)", bun.In(*u.BadgeIDs)).
		Scan(ctx)

	return err
}
func (r *Repository) SearchUserByID(ctx context.Context, id uint64) (*models.User, error) {
	u := &models.User{ID: id}
	err := r.DB.NewSelect().
		Model(u).
		WherePK().
		Relation("ActiveBan").
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &fiber.Error{Code: 404, Message: "user not found"}
		}
		return nil, err
	}

	if u.ActiveBan != nil {
		if err := r.DB.NewSelect().
			Model(u.ActiveBan).
			Relation("Admin").
			Relation("Target").
			WherePK().
			Scan(ctx); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	u.Badges = make([]models.Badge, 0)
	if u.BadgeIDs == nil {
		u.BadgeIDs = &[]uint64{}
	}
	if len(*u.BadgeIDs) > 0 {
		_ = r.populateUserBadges(ctx, u)
	}
	return u, nil
}

func (r *Repository) SearchUserByName(ctx context.Context, name string) (*models.User, error) {
	u := new(models.User)

	err := r.DB.NewSelect().
		Model(u).
		Where("user.name = ?", name).
		Relation("ActiveBan").
		Relation("ActiveBan.Admin").
		Relation("ActiveBan.Target").
		Limit(1).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &fiber.Error{Code: 404, Message: "user not found"}
		}
		return nil, err
	}

	u.Badges = make([]models.Badge, 0)
	if u.BadgeIDs == nil {
		u.BadgeIDs = &[]uint64{}
	}

	if len(*u.BadgeIDs) > 0 {
		_ = r.populateUserBadges(ctx, u)
	}

	return u, nil
}

func (r *Repository) SearchAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.DB.NewSelect().
		Model(&users).
		Relation("ActiveBan").
		Relation("ActiveBan.Admin").
		Relation("ActiveBan.Target").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []models.User{}, nil
	}

	for i := range users {
		users[i].Badges = make([]models.Badge, 0)

		if users[i].BadgeIDs == nil {
			users[i].BadgeIDs = &[]uint64{}
		}

		if len(*users[i].BadgeIDs) > 0 {
			if err := r.populateUserBadges(ctx, &users[i]); err != nil {
				slog.Errorf("failed to populate badges for user %d: %v", users[i].ID, err)
			}
		}
	}

	return users, nil
}

func (r *Repository) UpdateUser(ctx context.Context, tx bun.IDB, user *models.User, columns ...string) (*models.User, error) {
	query := tx.NewUpdate().
		Model(user).
		WherePK()

	if len(columns) > 0 {
		query.Column(columns...)
	} else {
		query.ExcludeColumn("created_at")
	}

	_, err := query.Exec(ctx)
	if err != nil {
		return user, err
	}

	updatedUser, err := r.SearchUserByID(ctx, user.ID)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, err
}

// ! Admin actions
func (r *Repository) CreateUser(ctx context.Context, tx bun.IDB, user *models.User) error {
	_, err := tx.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}

func (r *Repository) SoftDelete(ctx context.Context, tx bun.IDB, u *models.User) error {
	u.Name = u.Name + "_old"

	_, err := tx.NewUpdate().
		Model(u).
		Column("name").
		Set("status = ?", enums.UserStatusDeleted).
		WherePK().
		Exec(ctx)
	return err
}

func (r *Repository) LookupByDiscordID(ctx context.Context, discordID string) (*models.User, error) {
	u := new(models.User)

	err := r.DB.NewSelect().
		Model(u).
		Where("user.discord_id = ?", discordID).
		Relation("ActiveBan").
		Relation("ActiveBan.Admin").
		Relation("ActiveBan.Target").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &fiber.Error{Code: 404, Message: "user not found"}
		}
		return nil, err
	}

	u.Badges = make([]models.Badge, 0)
	if u.BadgeIDs == nil {
		u.BadgeIDs = &[]uint64{}
	}

	if len(*u.BadgeIDs) > 0 {
		_ = r.populateUserBadges(ctx, u)
	}

	return u, nil
}

func (r *Repository) HardDelete(ctx context.Context, tx bun.IDB, id uint64) error {
	_, err := tx.NewDelete().
		Model((*models.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *Repository) Restore(ctx context.Context, tx bun.IDB, user *models.User) error {
	user.Name = strings.ReplaceAll(user.Name, "_old", "")

	_, err := tx.NewUpdate().
		Model(user).
		Column("name").
		WherePK().
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			return &fiber.Error{Code: 403, Message: "cannot restore: nickname is already taken"}
		}
	}
	return err
}

func (r *Repository) SetStaffRank(ctx context.Context, tx bun.IDB, userID uint64, rankID int) (*models.User, error) {
	_, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Set("staff_rank = ?", rankID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) SetDeveloperRank(ctx context.Context, tx bun.IDB, userID uint64, rankID int) (*models.User, error) {
	_, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Set("developer_rank = ?", rankID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) EditUserFlags(ctx context.Context, tx bun.IDB, userID uint64, flags []string) (*models.User, error) {
	_, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Set("staff_flags = ?", flags).
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) EditUserBadges(ctx context.Context, tx bun.IDB, userID uint64, badgeIDs []uint64) (*models.User, error) {
	_, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Set("badges = ?", badgeIDs).
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) CreateBan(ctx context.Context, tx bun.IDB, ban *models.Ban) error {
	_, err := tx.NewInsert().
		Model(ban).
		Returning("id").
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = tx.NewUpdate().
		Model(&models.User{}).
		Set("active_ban = ?", ban.ID).
		Where("id = ?", ban.IssuedTo).
		Exec(ctx)
	return err
}

func (r *Repository) LiftBan(ctx context.Context, tx bun.IDB, userID uint64) error {
	_, err := tx.NewUpdate().
		Model(&models.Ban{}).
		Where("issued_to = ?", userID).
		Set("status = ?", enums.BanRemoved).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{
			"err": err,
		}).Error("err")
		return err
	}

	_, err = tx.NewUpdate().
		Table("users").
		Set("active_ban = ?", nil).
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{
			"err": err,
		}).Error("err")
		return err
	}

	return err
}

func (r *Repository) UpdateLastLogin(ctx context.Context, tx bun.IDB, userID uint64, ip string) error {
	_, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Set("last_login = ?", time.Now()).
		Set("last_ip = ?", ip).
		Set("register_ip = IF(register_ip IS NULL, ?, register_ip)", ip).
		Where("id = ?", userID).
		Exec(ctx)
	return err
}

func (r *Repository) ResetUserSensitiveData(ctx *fiber.Ctx, tx bun.IDB, userID uint64) error {
	_, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Set("last_login = ?", nil).
		Set("last_ip = ?", nil).
		Set("register_ip = ?", nil).
		Where("id = ?", userID).
		Exec(ctx.UserContext())

	if err != nil {
		return err
	}

	return err
}

func (r *Repository) PopulateBanList(ctx context.Context) ([]models.Ban, error) {
	var bans []models.Ban

	err := r.DB.NewSelect().
		Model(&bans).
		Relation("Admin").
		Relation("Target").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	if len(bans) == 0 {
		bans = make([]models.Ban, 0)
	}

	return bans, nil
}
