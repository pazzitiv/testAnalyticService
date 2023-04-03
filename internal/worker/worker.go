package worker

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"testAnalyticService/internal"
)

type Task struct {
	Id     uuid.UUID
	Type   TaskType
	Data   interface{}
	Status TaskStatus
}

type TaskType int

const (
	TaskTypeAnalytics TaskType = iota + 1
)

type TaskStatus int

const (
	TaskStatusCreated TaskStatus = iota + 1
	TaskStatusWorking
	TaskStatusCompleted
)

type Worker interface {
	worker(ctx context.Context, tasks chan *Task)
	AddTask(taskType TaskType, taskData interface{})
}

type worker struct {
	tasks         chan *Task
	analyticsRepo internal.AnalyticsRepository
	mutex         sync.Mutex
	logger        *zap.Logger
}

func NewWorker(ctx context.Context, analyticsRepo internal.AnalyticsRepository, logger *zap.Logger) Worker {
	tasks := make(chan *Task)
	newWorker := &worker{
		tasks:         tasks,
		analyticsRepo: analyticsRepo,
		mutex:         sync.Mutex{},
		logger:        logger,
	}
	go newWorker.worker(ctx, tasks)

	return newWorker
}

func (w *worker) worker(ctx context.Context, tasks chan *Task) {
	defer func() {
		if e := recover(); e != nil {
			w.logger.Panic("worker panic", zap.Any("error", e))
		}
	}()

	for {
		select {
		case task := <-tasks:
			task.Status = TaskStatusWorking

			switch task.Type {
			case TaskTypeAnalytics:
				data, ok := task.Data.(internal.AnalyticData)
				if !ok {
					w.logger.Error("can't get analytic data")
				}

				err := w.analyticsRepo.Add(ctx, data)
				if err != nil {
					w.logger.Error("set analytics data error", zap.Error(err))
				}

				task.Status = TaskStatusCompleted
			}
			if task.Status != TaskStatusCompleted {
				w.logger.Warn("exec task error", zap.Any("task", task))
			}
		case <-ctx.Done():
			break
		}
	}
}

func (w *worker) AddTask(taskType TaskType, taskData interface{}) {
	switch taskType {
	case TaskTypeAnalytics:
		data, ok := taskData.(internal.AnalyticData)

		if !ok {
			w.logger.Error("wrong analytic data", zap.Any("data", taskData))
			break
		}

		w.mutex.Lock()
		w.tasks <- &Task{
			Id:     uuid.New(),
			Type:   TaskTypeAnalytics,
			Data:   data,
			Status: TaskStatusCreated,
		}
		w.mutex.Unlock()
	}
}
