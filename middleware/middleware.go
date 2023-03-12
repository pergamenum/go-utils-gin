package middleware

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	er "github.com/pergamenum/go-utils-gin/responses"
	"go.uber.org/zap"
)

func AddRequestLogger(r *gin.Engine, logger *zap.Logger) {

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
}

func AddRecovery(r *gin.Engine, logger *zap.Logger) {

	recovery := func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			er.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		er.ErrorResponse(c, http.StatusInternalServerError)
	}
	r.Use(ginzap.CustomRecoveryWithZap(logger, true, recovery))
}

func AddAuth(_ *gin.Engine, _ *zap.Logger) {
	// TODO: Auth Middleware
	panic("NOT IMPLEMENTED")
}
