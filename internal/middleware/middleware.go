package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/ranks"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func LoadRank(rankSrv *ranks.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Locals("user")
		user, ok := val.(*models.User)
		if !ok {
			return c.Next()
		}

		ranksList := make([]*models.Rank, 0)

		if user.StaffRank > 0 {
			if r, err := rankSrv.SearchRankByID(c, user.StaffRank); err == nil {
				ranksList = append(ranksList, r)
			}
		}
		if user.DeveloperRank > 0 {
			if r, err := rankSrv.SearchRankByID(c, user.DeveloperRank); err == nil {
				ranksList = append(ranksList, r)
			}
		}

		c.Locals("rank", ranksList)
		return c.Next()
	}
}

func RequireFlag(flags ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Locals("rank")
		userVal := c.Locals("user")

		ranks, ok := val.([]*models.Rank)
		if !ok || len(ranks) == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		user, ok := userVal.(*models.User)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		for _, requiredFlag := range flags {
			for _, rank := range ranks {
				if rank != nil && rank.HasFlag(requiredFlag) {
					return c.Next()
				}
			}

			if user.Flags != nil {
				for _, userFlag := range *user.Flags {
					if userFlag == requiredFlag || userFlag == "MANAGER" {
						return c.Next()
					}
				}
			}
		}

		slog.WithData(slog.M{
			"required_flags": flags,
			"rank":           ranks,
			"user":           user,
		}).Errorf("action stopped because it must have flags: %v", flags)

		return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("forbidden. required flags: %s", flags))
	}
}

func UpdateLastLogin(repo *mysql.Repository) fiber.Handler {
	type cacheEntry struct {
		updatedAt time.Time
	}

	var (
		mu    sync.Mutex
		cache = make(map[uint64]cacheEntry)
	)

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			cutoff := time.Now().Add(-5 * time.Minute)
			mu.Lock()

			for id, e := range cache {
				if e.updatedAt.Before(cutoff) {
					delete(cache, id)
				}
			}

			mu.Unlock()
		}
	}()

	return func(c *fiber.Ctx) error {
		err := c.Next()

		u := account.GetUserFromContext(c)
		if u == nil {
			return err
		}

		mu.Lock()
		shouldUpdate := time.Since(cache[u.ID].updatedAt) > time.Minute
		if shouldUpdate {
			cache[u.ID] = cacheEntry{updatedAt: time.Now()}
		}
		mu.Unlock()

		if shouldUpdate {
			userID := u.ID
			ip := c.IP()

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if updateErr := repo.UpdateLastLogin(ctx, repo.DB, userID, ip); updateErr != nil {
					slog.WithData(slog.M{
						"userID": userID,
						"error":  updateErr,
					}).Warn("Failed to update last_login")
				}
			}()
		}

		return err
	}
}
