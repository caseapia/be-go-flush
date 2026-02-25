package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func (s *Service) generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *Service) GetOAuthURL(c *fiber.Ctx) (string, string, error) {
	state, err := s.generateState()
	if err != nil {
		return "", "", err
	}

	params := url.Values{}
	params.Add("client_id", s.cfg.DiscordClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", s.cfg.DiscordRedirectURI())
	params.Add("scope", "identify email")
	params.Add("state", state)

	authURL := "https://discord.com/api/oauth2/authorize?" + params.Encode()

	return authURL, state, nil
}

func (s *Service) validateState(c *fiber.Ctx, state string) bool {
	cookieState := c.Cookies("discord_oauth_state")

	slog.WithData(slog.M{
		"cookieState": cookieState,
		"state":       cookieState != "" && cookieState == state,
	}).Debug("Discord OAUTH state validation")

	return cookieState != "" && cookieState == state
}

func (s *Service) LinkDiscord(ctx *fiber.Ctx, userID uint64, code, state string) (*models.User, error) {
	if !s.validateState(ctx, state) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid state")
	}

	token, err := s.discord.ExchangeCode(ctx.UserContext(), code)
	if err != nil {
		return nil, err
	}

	discordUser, err := s.discord.GetDiscordUser(ctx.UserContext(), token.AccessToken)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.repository.LookupByDiscordID(ctx.UserContext(), discordUser.ID)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "discord already linked to account "+existingUser.Name)
	}

	user, err := s.repository.UpdateUser(
		ctx.UserContext(),
		&models.User{
			ID:          userID,
			DiscordID:   &discordUser.ID,
			DiscordName: &discordUser.GlobalName,
		},
		"discord_id",
		"discord_name",
	)
	if err != nil {
		return nil, err
	}

	ctx.ClearCookie("discord_oauth_state")

	addInfo := fmt.Sprintf("Discord ID: %v | Discord name: %s | Discord username: %s", discordUser.ID, discordUser.GlobalName, discordUser.Username)
	s.logger.Log(ctx.UserContext(), models.AdminAuthLogger, nil, &userID, models.LinkDiscord, addInfo)

	return user, nil
}

func (s *Service) UnlinkDiscord(ctx context.Context, sender *models.User, userID uint64) (*models.User, error) {
	user, err := s.repository.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.DiscordID == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "discord account is not linked to account "+user.Name)
	}

	oldDiscordID := *user.DiscordID

	user, err = s.repository.UpdateUser(ctx, &models.User{
		ID:          user.ID,
		DiscordID:   nil,
		DiscordName: nil,
	},
		"discord_id",
		"discord_name",
	)
	if err != nil {
		return nil, err
	}

	if sender.ID == user.ID {
		s.logger.Log(ctx, models.AdminAuthLogger, nil, &user.ID, models.UnlinkDiscord, "Discord ID: "+oldDiscordID)
	} else {
		s.logger.Log(ctx, models.StaffCommonLogger, &sender.ID, &user.ID, models.ForceUnlinkUserDiscord, "Discord ID: "+oldDiscordID)
	}

	return user, nil
}
