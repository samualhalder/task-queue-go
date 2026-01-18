package taskexecutor

import taskhandler "github.com/samualhalder/task-queue-go/internal/task_handler"


type Registry map[string]taskhandler.TaskHandler

func NewRegistry() Registry {
	return make(Registry)
}

func(r Registry) Register(taskType string, h taskhandler.TaskHandler) {
	
	if _,exists:=r[taskType];exists{
		panic("duplicate task handler" + taskType)
	}
	r[taskType]=h
}

