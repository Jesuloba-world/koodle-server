package taskservice

import (
	"fmt"

	"github.com/Jesuloba-world/koodle-server/model"
)

func checkIfColumnExists(cols []model.Column, columnId string) error {
	for _, col := range cols {
		if col.ID == columnId {
			return nil
		}
	}
	return fmt.Errorf("column %s not found in board", columnId)
}
