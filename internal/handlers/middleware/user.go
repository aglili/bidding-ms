package middleware

import (
	"net/http"

	"github.com/aglili/auction-app/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireUserAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		userID := session.Get("user_id")
		if userID == nil {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("unauthorized", nil))
			ctx.Abort()
			return
		}

		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("invalid session id", err))
			ctx.Abort()
			return
		}

		ctx.Set("user_id", uid.String())
		ctx.Next()
	}
}
