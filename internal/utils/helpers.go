package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetQueryInt(ctx *gin.Context, key string, defaultValue int) int {
	valStr := ctx.Query(key)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

func GetParamStr(ctx *gin.Context, key string, defaultValue string) string {
	valStr := ctx.Param(key)
	if valStr == "" {
		return defaultValue
	}
	return valStr
}
