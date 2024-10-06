package boardservice

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"

	humagroup "github.com/Jesuloba-world/koodle-server/lib/humaGroup"
	custommiddleware "github.com/Jesuloba-world/koodle-server/middleware"
	userrepo "github.com/Jesuloba-world/koodle-server/repo/user"
)

type BoardService struct {
	api      *humagroup.HumaGroup
	userRepo *userrepo.UserRepo
}

func NewBoardService(api huma.API, middleware *custommiddleware.Middleware, userRepo *userrepo.UserRepo) *BoardService {
	return &BoardService{
		api:      humagroup.NewHumaGroup(api, "/boards", []string{"Board"}, middleware.Auth),
		userRepo: userRepo,
	}
}

func (s *BoardService) RegisterRoutes() {
	humagroup.Post(s.api, "", s.createBoard, "Create Board")
	humagroup.Get(s.api, "", s.getAllBoards, "Get All Boards")
	humagroup.Get(s.api, "/{boardId}", s.getBoard, "Get Board")
	humagroup.Delete(s.api, "/{boardId}", s.deleteBoard, "Delete Board")
	humagroup.Put(s.api, "/{boardId}", s.updateBoard, "Update Board")
}

func (s *BoardService) createBoard(ctx context.Context, req *createBoardReq) (*createBoardResp, error) {
	user, err := s.userRepo.GetUserByCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Unauthorized", err)
	}
	slog.Info("User", "user_id", user.Email)

	resp := &createBoardResp{}
	resp.Body.Message = "Board created successfully"
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

func (s *BoardService) getAllBoards(ctx context.Context, req *getAllBoardsReq) (*getAllBoardsResp, error) {
	resp := &getAllBoardsResp{}
	return resp, nil
}
