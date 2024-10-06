package boardservice

import (
	"time"

	"github.com/Jesuloba-world/koodle-server/model"

)

type boardObject struct {
	ID        string    `json:"id" example:"136789874673893" doc:"ID of board"`
	Name      string    `json:"name" example:"Web design" doc:"Name of board"`
	CreatedAt time.Time `json:"createdAt" doc:"Time board was created"`
	UpdatedAt time.Time `json:"updatedAt" doc:"Time board was last updated"`
}

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
		Message string       `json:"message" example:"Board created successfully"`
		Board   *boardObject `json:"board" doc:"Created board"`
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
		Boards []*boardObject `json:"boards" doc:"List of boards"`
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
		Message string       `json:"message" example:"Board updated successfully"`
		Board   *boardObject `json:"board" doc:"Updated board"`
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
