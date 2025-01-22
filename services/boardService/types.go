package boardservice

import (
	"github.com/Jesuloba-world/koodle-server/model"
)

type boardInput struct {
	Name    string   `json:"name" example:"Web design" doc:"Name of board"`
	Columns []string `json:"columns" example:"To Do,In Progress,Done" doc:"Columns of board"`
}

type createBoardReq struct {
	Body struct {
		Board *boardInput `json:"board" doc:"Board to create"`
	}
}

type createBoardResp struct {
	Body struct {
		Message string             `json:"message" example:"Board created successfully"`
		Board   *model.BoardObject `json:"board" doc:"Created board"`
	}
}

type getBoardReq struct {
	BoardId string `path:"boardId" example:"136789874673893" doc:"ID of board"`
}

type getBoardResp struct {
	Body struct {
		Board *model.Board `json:"board" doc:"The full Board object"`
	}
}

type getAllBoardsReq struct{}

type getAllBoardsResp struct {
	Body struct {
		Boards []*model.BoardObject `json:"boards" doc:"List of boards"`
	}
}

type updateBoardReq struct {
	BoardId string `path:"boardId" example:"136789874673893" doc:"ID of board"`
	Body    struct {
		Board *boardInput `json:"board" doc:"Board to update"`
	}
}

type updateBoardResp struct {
	Body struct {
		Message string             `json:"message" example:"Board updated successfully"`
		Board   *model.BoardObject `json:"board" doc:"Updated board"`
	}
}

type deleteBoardReq struct {
	BoardId string `path:"boardId" example:"136789874673893" doc:"ID of board"`
}

type deleteBoardResp struct {
	Body struct {
		Message string `json:"message" example:"Board deleted successfully"`
	}
}
