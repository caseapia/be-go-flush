package auth

import (
	"context"
)

func (s *Service) CheckMultiAccount(ctx context.Context, userID int64, fingerprint, ip, userAgent string) (bool, error) {
	_, _, err := s.repository.CheckMultiAccountByFingerprint(ctx, userID, fingerprint)
	if err != nil {
		return false, err
	}

	_, _, err = s.repository.CheckMultiAccountByIP(ctx, userID, ip)
	if err != nil {
		return false, err
	}

	_, _, err = s.repository.CheckMultiAccountByUA(ctx, userID, userAgent)
	if err != nil {
		return false, err
	}

	return false, nil
}
