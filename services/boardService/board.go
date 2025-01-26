package boardservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"

	humagroup "github.com/Jesuloba-world/koodle-server/lib/humaGroup"
	custommiddleware "github.com/Jesuloba-world/koodle-server/middleware"
	"github.com/Jesuloba-world/koodle-server/model"
	boardrepo "github.com/Jesuloba-world/koodle-server/repo/board"
	userrepo "github.com/Jesuloba-world/koodle-server/repo/user"
)

type BoardService struct {
	api       *humagroup.HumaGroup
	userRepo  *userrepo.UserRepo
	boardRepo *boardrepo.BoardRepo
}

func NewBoardService(api huma.API, middleware *custommiddleware.Middleware, userRepo *userrepo.UserRepo, boardRepo *boardrepo.BoardRepo) *BoardService {
	return &BoardService{
		api:       humagroup.NewHumaGroup(api, "/boards", []string{"Board"}, middleware.Auth),
		userRepo:  userRepo,
		boardRepo: boardRepo,
	}
}

func (s *BoardService) RegisterRoutes() {
	humagroup.Post(s.api, "", s.createBoard, "Create Board")
	humagroup.Get(s.api, "/mine", s.getAllMyBoards, "Get All My Boards")
	humagroup.Get(s.api, "/{boardId}", s.getBoard, "Get Board")
	humagroup.Delete(s.api, "/{boardId}", s.deleteBoard, "Delete Board")
	humagroup.Patch(s.api, "/{boardId}", s.updateBoard, "Update Board")
}

func (s *BoardService) createBoard(ctx context.Context, req *createBoardReq) (*createBoardResp, error) {
	user, err := s.userRepo.GetUserByCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Unauthorized", err)
	}

	board := &model.Board{
		Name:   req.Body.Board.Name,
		UserId: user.ID,
	}

	var columns []*model.Column
	for i, colInput := range req.Body.Board.Columns {
		columns = append(columns, &model.Column{
			Name:     colInput.Name,
			Position: i,
		})
	}

	if err := s.boardRepo.CreateBoardWithColumn(ctx, board, columns); err != nil {
		slog.Error("Failed to create board", "error", err)
		return nil, huma.Error500InternalServerError("Failed to create board", err)
	}

	resp := &createBoardResp{}
	resp.Body.Message = "Board created successfully"
	resp.Body.Board = board.MapBoardToResponse()
	return resp, nil
}

func (s *BoardService) getBoard(ctx context.Context, req *getBoardReq) (*getBoardResp, error) {
	user, err := s.userRepo.GetUserByCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Unauthorized", err)
	}

	board, err := s.boardRepo.GetBoardWithColumns(ctx, req.BoardId)
	if err != nil {
		slog.Error("Failed to get board", "error", err, "boardId", req.BoardId)
		return nil, huma.Error404NotFound("board not found", err)
	}

	// check access
	if board.UserId != user.ID {
		return nil, huma.Error401Unauthorized("you don't have access", fmt.Errorf("user is not authorized to this board"))
	}

	resp := &getBoardResp{}
	resp.Body.Board = board
	return resp, nil
}

func (s *BoardService) updateBoard(ctx context.Context, req *updateBoardReq) (*updateBoardResp, error) {
	user, err := s.userRepo.GetUserByCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Unauthorized", err)
	}

	board, err := s.boardRepo.GetBoardWithColumns(ctx, req.BoardId)
	if err != nil {
		return nil, huma.Error404NotFound("board not found", err)
	}

	// check access
	if board.UserId != user.ID {
		return nil, huma.Error401Unauthorized("you don't have access", fmt.Errorf("user is not authorized to this board"))
	}

	board.Name = req.Body.Board.Name

	existingColumns := make(map[string]*model.Column)
	for _, col := range board.Columns {
		existingColumns[col.ID] = &col
	}

	// track which columns are updated
	processedColumns := make(map[string]bool)
	var updatedColumns []*model.Column

	// process the columns in the req
	for i, colInput := range req.Body.Board.Columns {
		if colInput.ID != "" {
			// update existing column
			existingCol, exists := existingColumns[colInput.ID]
			if !exists {
				return nil, huma.Error400BadRequest("Invalid column ID", fmt.Errorf("column with ID %s not found", colInput.ID))
			}
			existingCol.Name = colInput.Name
			existingCol.Position = i
			updatedColumns = append(updatedColumns, existingCol)
			processedColumns[colInput.ID] = true
		} else {
			// create a new column
			updatedColumns = append(updatedColumns, &model.Column{
				Name:     colInput.Name,
				Position: i,
				BoardID:  board.ID,
			})
		}
	}

	// find coluumns to delete (columns not in the request)
	var columnsToDelete []string
	for id, col := range existingColumns {
		if !processedColumns[id] {
			columnsToDelete = append(columnsToDelete, col.ID)
		}
	}

	// update the board and columns
	if err := s.boardRepo.UpdateBoardWithColumns(ctx, board, updatedColumns, columnsToDelete); err != nil {
		slog.Error("Failed to update board", "error", err)
		return nil, huma.Error500InternalServerError("failed to update board", err)
	}

	resp := &updateBoardResp{}
	resp.Body.Message = "Board updated Successfully"
	resp.Body.Board = board.MapBoardToResponse()
	return resp, nil
}

func (s *BoardService) deleteBoard(ctx context.Context, req *deleteBoardReq) (*deleteBoardResp, error) {
	resp := &deleteBoardResp{}
	return resp, nil
}

func (s *BoardService) getAllMyBoards(ctx context.Context, req *getAllBoardsReq) (*getAllBoardsResp, error) {
	user, err := s.userRepo.GetUserByCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Unauthorized", err)
	}

	boards, err := s.boardRepo.GetBoardsByUser(ctx, user.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to get boards", err)
	}

	boardsResponse := make([]*model.BoardObject, 0, len(boards))
	for _, board := range boards {
		boardsResponse = append(boardsResponse, board.MapBoardToResponse())
	}

	resp := &getAllBoardsResp{}
	resp.Body.Boards = boardsResponse
	return resp, nil
}
