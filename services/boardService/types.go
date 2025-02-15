package boardservice

import (
	"github.com/Jesuloba-world/koodle-server/model"
)

type BaseRequest struct {
	Authorization string `json:"authorization" header:"Authorization" example:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTksImV4cCI6MTc0MDM1MDUzNX0.MJdLqfXBmSv14PCMZXfRscMdiFOtyBtek5K6l4xZ2XA" doc:"Bearer [token]"`
}

type columnInput struct {
	ID   string `json:"id,omitempty" example:"col_123" doc:"ID of column (empty for new columns)"`
	Name string `json:"name" example:"In Progress" doc:"Name of column"`
}

type boardInput struct {
	Name    string        `json:"name" example:"Web design" doc:"Name of board"`
	Columns []columnInput `json:"columns" doc:"Columns of board"`
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
	BaseRequest
	BoardId     string `path:"boardId" example:"136789874673893" doc:"ID of board"`
	IncludeTask bool   `query:"include_task" example:"false" doc:"Whether to include tasks in the board response"`
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
