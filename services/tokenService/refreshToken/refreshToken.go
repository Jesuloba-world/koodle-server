package refreshtoken

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/model"
)

type RefreshTokenRepo struct {
	db *bun.DB
}

func NewRefreshTokenRepo(db *bun.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) StoreRefreshToken(refreshToken, userId string, expiresAt time.Time) error {
	token := &model.RefreshToken{
		RefreshToken: refreshToken,
		UserID:       userId,
		ExpiresAt:    expiresAt,
	}

	_, err := r.db.NewInsert().Model(token).Exec(context.Background())
	return err
}

func (r *RefreshTokenRepo) RevokeRefreshToken(refreshtoken string) error {
	_, err := r.db.NewUpdate().Model(&model.RefreshToken{}).Set("revoked = true").Where("refresh_token = ?", refreshtoken).Exec(context.Background())

	return err
}

func (r *RefreshTokenRepo) IsRefreshTokenRevoked(refreshtoken string) (bool, error) {
	token := new(model.RefreshToken)
	err := r.db.NewSelect().Model(token).Where("refresh_token = ?", refreshtoken).Scan(context.Background())

	if errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return token.Revoked, nil
}
