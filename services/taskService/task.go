package taskservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"

	humagroup "github.com/Jesuloba-world/koodle-server/lib/humaGroup"
	custommiddleware "github.com/Jesuloba-world/koodle-server/middleware"
	"github.com/Jesuloba-world/koodle-server/model"
	boardrepo "github.com/Jesuloba-world/koodle-server/repo/board"
	taskrepo "github.com/Jesuloba-world/koodle-server/repo/task"
	userrepo "github.com/Jesuloba-world/koodle-server/repo/user"
)

type TaskService struct {
	api       *humagroup.HumaGroup
	userRepo  *userrepo.UserRepo
	boardRepo *boardrepo.BoardRepo
	taskrepo  *taskrepo.TaskRepo
}

func NewTaskService(
	api huma.API,
	middleware *custommiddleware.Middleware,
	userRepo *userrepo.UserRepo,
	boardRepo *boardrepo.BoardRepo,
	taskrepo *taskrepo.TaskRepo,
) *TaskService {
	return &TaskService{
		api:       humagroup.NewHumaGroup(api, "/tasks", []string{"Task"}, middleware.Auth),
		userRepo:  userRepo,
		boardRepo: boardRepo,
		taskrepo:  taskrepo,
	}
}

func (t *TaskService) RegisterRoutes() {
	humagroup.Post(t.api, "", t.addTask, "Add Task")
}

func (t *TaskService) addTask(ctx context.Context, req *addTaskReq) (*addTaskResp, error) {
	user, err := t.userRepo.GetUserByCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("Unauthorized", err)
	}

	board, err := t.boardRepo.GetBoardWithColumns(ctx, req.Body.BoardId, true)
	if err != nil {
		slog.Error("Failed to get board", "error", err, "boardId", req.Body.BoardId)
		return nil, huma.Error404NotFound("board not found", err)
	}

	// check access
	if board.UserId != user.ID {
		return nil, huma.Error401Unauthorized("you don't have access", fmt.Errorf("user is not authorized to this board"))
	}

	// check if column exists
	err = checkIfColumnExists(board.Columns, req.Body.ColumnId)
	if err != nil {
		return nil, huma.Error404NotFound("column not found", err)
	}

	var maxPosition int
	for _, col := range board.Columns {
		if col.ID == req.Body.ColumnId {
			for _, task := range col.Tasks {
				if task.Position > maxPosition {
					maxPosition = task.Position
				}
			}
			break
		}
	}

	task := &model.Task{
		ColumnID:    req.Body.ColumnId,
		Title:       req.Body.Task.Title,
		Description: req.Body.Task.Description,
		Position:    maxPosition + 1,
	}

	var subtasks []*model.Subtask
	for _, subtaskInput := range req.Body.Task.SubTasks {
		subtasks = append(subtasks, &model.Subtask{
			Name: subtaskInput.Name,
		})
	}

	if err := t.taskrepo.AddTaskToColumn(ctx, task, subtasks); err != nil {
		slog.Error("Failed to create task", "error", err)
		return nil, huma.Error500InternalServerError("Failed to create task", err)
	}

	resp := &addTaskResp{}
	resp.Body.Message = "Task added successfully"
	resp.Body.Task = task
	return resp, nil
}
