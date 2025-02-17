package taskservice

import "github.com/Jesuloba-world/koodle-server/model"

type subtaskInput struct {
	ID   string `json:"id,omitempty" example:"sub_123" doc:"ID of subtask (empty for new subtasks)"`
	Name string `json:"name" example:"Make coffee" doc:"name of subtask"`
}

type taskInput struct {
	Title       string         `json:"title" example:"Take Coffee Break" doc:"title of task"`
	Description string         `json:"description" example:"It's always good to take a break. This 15 minute break will recharge the batteries a little." doc:"description of task"`
	SubTasks    []subtaskInput `json:"subtasks" doc:"Subtasks of task"`
}

type addTaskReq struct {
	Body struct {
		BoardId  string     `json:"boardId" example:"136789874673893" doc:"ID of board task will be added to"`
		Task     *taskInput `json:"task" doc:"Task to be added"`
		ColumnId string     `json:"columnId" doc:"column the task should belong to" example:"col_123"`
	}
}

type addTaskResp struct {
	Body struct {
		Message string      `json:"message" example:"Task created successfully"`
		Task    *model.Task `json:"task" doc:"Task added"`
	}
}
