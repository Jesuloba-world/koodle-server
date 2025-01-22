package boardservice

import (
	"context"
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
	humagroup.Put(s.api, "/{boardId}", s.updateBoard, "Update Board")
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
	for i, colName := range req.Body.Board.Columns {
		columns = append(columns, &model.Column{
			Name:     colName,
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
	resp := &getBoardResp{}
	return resp, nil
}

func (s *BoardService) updateBoard(ctx context.Context, req *updateBoardReq) (*updateBoardResp, error) {
	resp := &updateBoardResp{}
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
