package middleware

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	m "github.com/pergamenum/go-utils-gin/messages"
	"go.uber.org/zap"
)

func AddRequestLogger(r *gin.Engine, logger *zap.Logger) {

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
}

func AddRecovery(r *gin.Engine, logger *zap.Logger) {

	recovery := func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			m.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		m.ErrorResponse(c, http.StatusInternalServerError)
	}
	r.Use(ginzap.CustomRecoveryWithZap(logger, true, recovery))
}

func AddAuth(r *gin.Engine, logger *zap.Logger) {
	// TODO: Auth Middleware
	panic("NOT IMPLEMENTED")
}
