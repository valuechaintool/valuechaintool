package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Error(code int, message string) (int, gin.H) {
	err := gin.H{
		"code": strconv.Itoa(code),
		"msg":  message,
	}
	return code, err
}
