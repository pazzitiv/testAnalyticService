package http

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"testAnalyticService/internal"
)

type handlers struct {
	analysticsRepo internal.AnalyticsRepository
	logger         *zap.Logger
}

func NewHandlers(analysticsRepo internal.AnalyticsRepository, logger *zap.Logger) *handlers {
	return &handlers{
		analysticsRepo: analysticsRepo,
		logger:         logger,
	}
}

func (h *handlers) analitycsHandler(w http.ResponseWriter, r *http.Request) {
	type statusResponse struct {
		Status string
	}
	ctx := context.Background()

	switch r.Method {
	case http.MethodPost:
		var body internal.AnalyticBody

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&body)
		if err != nil {
			h.logger.Error("request decode error", zap.Error(err))
		}
		data := internal.AnalyticData{
			Headers: r.Header,
			Body:    body,
		}
		userId := r.Header.Get("X-Tantum-Authorization")
		go func(userId string, data internal.AnalyticData) {
			defer func() {
				if e := recover(); e != nil {
					h.logger.Error("add analytics data error", zap.Any("error", e))
				}
			}()
			err := h.analysticsRepo.Add(ctx, userId, data)
			if err != nil {
				h.logger.Error("set analytics data error", zap.Error(err))
			}
		}(userId, data)

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
