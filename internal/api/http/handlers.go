package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"testAnalyticService/internal"
	"testAnalyticService/internal/worker"
)

type handlers struct {
	worker worker.Worker
	logger *zap.Logger
}

func NewHandlers(worker worker.Worker, logger *zap.Logger) *handlers {
	return &handlers{
		worker: worker,
		logger: logger,
	}
}

func (h *handlers) analitycsHandler(w http.ResponseWriter, r *http.Request) {
	type statusResponse struct {
		Status string
	}

	switch r.Method {
	case http.MethodPost:
		var body internal.AnalyticBody

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&body)
		if err != nil {
			h.logger.Error("request decode error", zap.Error(err))
		} else {
			h.worker.AddTask(worker.TaskTypeAnalytics, internal.AnalyticData{
				Headers: r.Header,
				Body:    body,
			})
		}

		response := statusResponse{Status: http.StatusText(http.StatusOK)}
		responseData, err := json.Marshal(response)
		if err != nil {
			h.logger.Error("response marshal error", zap.Error(err))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(responseData)
		if err != nil {
			h.logger.Error("response write error", zap.Error(err))
			return
		}
	default:
		http.NotFound(w, r)
	}
}
