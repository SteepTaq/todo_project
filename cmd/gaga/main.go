package main

import (
	"log"

	todov1 "github.com/SteepTaq/todo_project/pkg/proto/gen/todo"
)

func main() {

	log.Printf("%T", todov1.TaskStatus_TASK_STATUS_IN_PROGRESS)

}
