package handler

import (
	"net/http"
	"strconv"
	"task-scheduler/internal/dto"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, dto.Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, dto.Response{
		Code:    code,
		Message: message,
	})
}

func parseUint64Param(c *gin.Context, name string) (uint64, error) {
	val := c.Param(name)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
