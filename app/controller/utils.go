package controller

import (
	"RssReader/pkg/errorhelper"
	"github.com/gin-gonic/gin"
	"net/http"
)

const PaginationDefaultPageSize = 20

type ApiResp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationReq struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

type PaginationResp struct {
	Total   int64       `json:"total"`
	PerPage int64       `json:"per_page"`
	Page    int64       `json:"page"`
	Data    interface{} `json:"data"`
}

func returnErrorMsg(c *gin.Context, err error) {
	var errMsg string
	errCode, errMsg := errorhelper.DecodeErr(err)
	c.JSON(http.StatusOK, ApiResp{
		Code:    int32(errCode),
		Message: errMsg,
	})
	c.Abort()
}

func returnSuccessMsg(c *gin.Context, message string, data interface{}) {
	if len(message) == 0 {
		message = "OK"
	}
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, ApiResp{
		Code:    200,
		Message: message,
		Data:    data,
	})
	c.Abort()
}
