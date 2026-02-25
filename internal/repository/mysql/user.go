package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (r *Repository) populateUserBadges(ctx context.Context, u *models.User) error {
	if len(u.BadgeIDs) == 0 {
		u.Badges = make([]models.Badge, 0)
		return nil
	}

	err := r.db.NewSelect().
		Model(&u.Badges).
		Where("id IN (?)", bun.In(u.BadgeIDs)).
		Scan(ctx)

	return err
}
func (r *Repository) SearchUserByID(ctx context.Context, id uint64) (*models.User, error) {
	u := &models.User{ID: id}

	err := r.db.NewSelect().
		Model(u).
		WherePK().
		Relation("ActiveBan").
		Relation("ActiveBan.Admin").
		Relation("ActiveBan.Target").
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	u.Badges = make([]models.Badge, 0)
	if u.BadgeIDs == nil {
		u.BadgeIDs = make([]uint64, 0)
	}

	if len(u.BadgeIDs) > 0 {
		_ = r.populateUserBadges(ctx, u)
	}

	return u, nil
}

func (r *Repository) SearchUserByName(ctx context.Context, name string) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("user.name = ?", name).
		Relation("ActiveBan").
		Relation("ActiveBan.Admin").
		Relation("ActiveBan.Target").
		Limit(1).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	u.Badges = make([]models.Badge, 0)
	if u.BadgeIDs == nil {
		u.BadgeIDs = make([]uint64, 0)
	}

	if len(u.BadgeIDs) > 0 {
		_ = r.populateUserBadges(ctx, u)
	}

	return u, nil
}

func (r *Repository) SearchAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.NewSelect().
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

		if len(users[i].BadgeIDs) > 0 {
			if err := r.populateUserBadges(ctx, &users[i]); err != nil {
				slog.Errorf("failed to populate badges for user %d: %v", users[i].ID, err)
			}
		}
	}

	return users, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User, columns ...string) (*models.User, error) {
	query := r.db.NewUpdate().Model(user).WherePK()

	if len(columns) > 0 {
		query.Column(columns...)
	} else {
		query.ExcludeColumn("created_at")
	}

	_, err := query.Exec(ctx)

	updatedUser, err := r.SearchUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return updatedUser, err
}

// ! Admin actions
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}

func (r *Repository) SoftDelete(ctx context.Context, u *models.User) error {
	u.Name = u.Name + "_old"

	_, err := r.db.NewUpdate().
		Model(u).
		Column("name").
		WherePK().
		Exec(ctx)
	return err
}

func (r *Repository) LookupByDiscordID(ctx context.Context, discordID string) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("user.discord_id = ?", discordID).
		Relation("ActiveBan").
		Relation("ActiveBan.Admin").
		Relation("ActiveBan.Target").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	u.Badges = make([]models.Badge, 0)
	if u.BadgeIDs == nil {
		u.BadgeIDs = make([]uint64, 0)
	}

	if len(u.BadgeIDs) > 0 {
		_ = r.populateUserBadges(ctx, u)
	}

	return u, nil
}

func (r *Repository) ChangeUserData(ctx context.Context, u *models.User, updateName, updateEmail, updatePassword bool) error {
	var conditions []string
	if updateName {
		conditions = append(conditions, "name")
	}
	if updateEmail {
		conditions = append(conditions, "email")
	}
	if updatePassword {
		conditions = append(conditions, "password")
	}

	for _, col := range conditions {
		val := ""
		if col == "name" {
			val = u.Name
		}
		if col == "email" {
			val = u.Email
		}
		if col == "password" {
			val = u.Password
		}

		exists, err := r.db.NewSelect().
			Model((*models.User)(nil)).
			Where("? = ?", bun.Ident(col), val).
			Where("id != ?", u.ID).
			Exists(ctx)

		if err != nil {
			return err
		}
		if exists && col != "password" {
			return fiber.NewError(fiber.StatusConflict, "User with this "+col+" already exists")
		}
	}

	query := r.db.NewUpdate().Model(u).WherePK()

	if updateName {
		query.Column("name")
	}
	if updateEmail {
		query.Column("email")
	}
	if updatePassword {
		query.Column("password")
	}

	if !updateName && !updateEmail && !updatePassword {
		return nil
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *Repository) HardDelete(ctx context.Context, id uint64) error {
	_, err := r.db.NewDelete().
		Model((*models.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *Repository) Restore(ctx context.Context, user *models.User) error {
	user.Name = strings.ReplaceAll(user.Name, "_old", "")

	_, err := r.db.NewUpdate().
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

func (r *Repository) SetStaffRank(ctx context.Context, userID uint64, rankID int) (*models.User, error) {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("staff_rank = ?", rankID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) SetDeveloperRank(ctx context.Context, userID uint64, rankID int) (*models.User, error) {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("developer_rank = ?", rankID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) EditUserFlags(ctx context.Context, userID uint64, flags []string) (*models.User, error) {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("staff_flags = ?", flags).
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) EditUserBadges(ctx context.Context, userID uint64, badgeIDs []uint64) (*models.User, error) {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("badges = ?", badgeIDs).
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) CreateBan(ctx context.Context, ban *models.BanModel) error {
	_, err := r.db.NewInsert().
		Model(ban).
		Returning("id").
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.NewUpdate().
		Model(&models.User{}).
		Set("active_ban = ?", ban.ID).
		Where("id = ?", ban.IssuedTo).
		Exec(ctx)
	return err
}

func (r *Repository) DeleteBan(ctx context.Context, userID uint64) error {
	_, err := r.db.NewDelete().
		Model(&models.BanModel{}).
		Where("issued_to = ?", userID).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{
			"err": err,
		}).Error("err")
		return err
	}

	_, err = r.db.NewUpdate().
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

func (r *Repository) UpdateLastLogin(ctx *fiber.Ctx, userID uint64) error {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("last_login = ?", time.Now()).
		Set("last_ip = ?", ctx.IP()).
		Set("register_ip = IF(register_ip IS NULL, ?, register_ip)", ctx.IP()).
		Where("id = ?", userID).
		Exec(ctx.UserContext())
	return err
}

func (r *Repository) ResetUserSensitiveData(ctx *fiber.Ctx, userID uint64) error {
	_, err := r.db.NewUpdate().
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
