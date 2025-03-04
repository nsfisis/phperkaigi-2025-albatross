package taskqueue

import (
	"github.com/hibiken/asynq"
)

type WorkerServer struct {
	server    *asynq.Server
	processor *processorWrapper
}

func NewWorkerServer(redisAddr string) *WorkerServer {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		asynq.Config{},
	)
	processor := newProcessorWrapper(newProcessor())
	return &WorkerServer{
		server:    server,
		processor: processor,
	}
}

func (s *WorkerServer) Run() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(string(TaskTypeRunTestcase), s.processor.processTaskRunTestcase)

	return s.server.Run(mux)
}

func (s *WorkerServer) Results() chan TaskResult {
	return s.processor.results
}
