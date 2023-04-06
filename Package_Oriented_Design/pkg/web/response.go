package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorRespose struct {
	Status  string `json: "status"`
	Code    int    `json: "status"`
	Message string `json: "status"`
}

type respose struct {
	Data interface{} `json: "status"`
}

func ErrorResp(ctx *gin.Context, err error, code int) {
	ctx.JSON(code, errorRespose{Status: http.StatusText(code), Code: code, Message: err.Error()})
}

func OkResp(ctx *gin.Context, code int, status interface{}) {
	ctx.JSON(code, respose{Data: status})
}
