package handlers

import (
	"fmt"
	"github.com/devfeel/dottask"
	"time"
)

func Task_Print(context *task.TaskContext) error {
	fmt.Println(time.Now(), "Task_Print", context.TaskID)
	return nil
}
