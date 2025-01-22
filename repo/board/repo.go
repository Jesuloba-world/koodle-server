package boardrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/model"
)

var ErrBoardNotFound = fmt.Errorf("board not found")

type BoardRepo struct {
	db *bun.DB
}

func NewBoardRepo(db *bun.DB) *BoardRepo {
	return &BoardRepo{
		db: db,
	}
}

func (r *BoardRepo) CreateBoard(ctx context.Context, board *model.Board) error {
	_, err := r.db.NewInsert().Model(board).Exec(ctx)
	return err
}

func (r *BoardRepo) CreateBoardWithColumn(ctx context.Context, board *model.Board, columns []*model.Column) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(board).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create board: %w", err)
		}

		if len(columns) > 0 {
			for _, column := range columns {
				column.BoardID = board.ID
			}
			if _, err := tx.NewInsert().Model(&columns).Exec(ctx); err != nil {
				return fmt.Errorf("failed to create columns: %w", err)
			}
		}

		return nil
	})
}

func (r *BoardRepo) GetBoardByID(ctx context.Context, id string) (*model.Board, error) {
	board := new(model.Board)
	err := r.db.NewSelect().
		Model(board).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBoardNotFound
		}
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	return board, nil
}

func (r *BoardRepo) UpdateBoard(ctx context.Context, board *model.Board) error {
	_, err := r.db.NewUpdate().
		Model(board).
		WherePK().
		Exec(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrBoardNotFound
		}
		return fmt.Errorf("failed to update board: %w", err)
	}

	return nil
}

func (r *BoardRepo) DeleteBoard(ctx context.Context, id string) error {
	result, err := r.db.NewDelete().
		Model((*model.Board)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrBoardNotFound
	}

	return nil
}

func (r *BoardRepo) GetAllBoards(ctx context.Context) ([]*model.Board, error) {
	var boards []*model.Board
	err := r.db.NewSelect().
		Model(&boards).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get all boards: %w", err)
	}

	return boards, nil
}

func (r *BoardRepo) GetBoardsByUser(ctx context.Context, userId string) ([]*model.Board, error) {
	var boards []*model.Board
	err := r.db.NewSelect().
		Model(&boards).
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get boards for user %s: %w", userId, err)
	}

	return boards, nil
}
