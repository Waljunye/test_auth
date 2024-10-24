package v1

import (
	"auth/internal/contracts/open_api/v1/dto"
	"auth/internal/metrics"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type BaseContract struct {
	logger *slog.Logger
}

func NewBaseContract() *BaseContract {
	return &BaseContract{}
}

func (bc *BaseContract) Ping(ginC *gin.Context) {
	var (
		payload *dto.PingResponse
	)

	payload = &dto.PingResponse{
		Data: "pong",
	}
	metrics.HttpRequestTotal.WithLabelValues(ginC.Request.Method, "200").Inc()

	ginC.JSON(200, payload)
}
