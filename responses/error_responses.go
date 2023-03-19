package responses

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	e "github.com/pergamenum/go-consensus-standards/ehandler"
)

func ErrorResponse(ctx *gin.Context, status int, message ...any) {

	var str string
	if len(message) == 1 {
		m := message[0]
		if v, ok := m.(string); ok {
			str = v
		}
		if v, ok := m.(error); ok {
			var valErrs validator.ValidationErrors
			if errors.As(v, &valErrs) {
				var sb strings.Builder
				sb.WriteString("required: ")
				for _, ev := range valErrs {
					sb.WriteString(fmt.Sprintf("(%s) ", ev.Field()))
				}
				str = strings.TrimSpace(sb.String())
			} else {
				str = v.Error()
			}
		}
	}

	jo := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  status,
		Message: str,
	}
	ctx.AbortWithStatusJSON(status, jo)
}

func StandardResponses(ctx *gin.Context, err error) {

	if errors.Is(err, e.ErrConflict) {
		ErrorResponse(ctx, http.StatusConflict, err)
		return
	}

	if errors.Is(err, e.ErrNotFound) {
		ErrorResponse(ctx, http.StatusNotFound, err)
		return
	}

	if errors.Is(err, e.ErrBadRequest) {
		ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// Any error below this point is to be logged.
	// This defers the logging responsibility to the gin middleware request/response logger.
	_ = ctx.Error(err)

	if errors.Is(err, e.ErrBadGateway) {
		ErrorResponse(ctx, http.StatusBadGateway, err)
		return
	}

	ErrorResponse(ctx, http.StatusInternalServerError, err)
}
