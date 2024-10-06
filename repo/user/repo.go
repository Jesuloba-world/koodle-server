package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

	custommiddleware "github.com/Jesuloba-world/koodle-server/middleware"
	"github.com/Jesuloba-world/koodle-server/model"
)

type UserRepo struct {
	db *bun.DB
}

func NewUserRepo(db *bun.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (s *UserRepo) CheckUserExists(ctx context.Context, email string) (bool, bool, error) {
	var user model.User
	exists, err := s.db.NewSelect().
		Model(&user).
		Column("email", "password").
		Where("email = ?", email).
		Exists(ctx)

	if err != nil {
		return false, false, fmt.Errorf("error checking user existence: %w", err)
	}

	passwordSet := user.Password != ""

	return exists, passwordSet, nil
}

func (s *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	_, err := s.db.NewInsert().
		Model(user).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (s *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := s.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		if err == sql.ErrNoRows {
			return nil, errors.New("no user found with the provided email")
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return user, nil
}

func (s *UserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	user := new(model.User)
	err := s.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error finding user by ID: %w", err)
	}

	return user, nil
}

func (s *UserRepo) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := s.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (s *UserRepo) GetUserByCtx(ctx context.Context) (*model.User, error) {
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, custommiddleware.ErrInvalidID
	}

	user, err := s.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
