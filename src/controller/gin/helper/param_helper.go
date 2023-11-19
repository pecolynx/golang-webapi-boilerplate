package helper

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleFunction(c *gin.Context, fn func() error, errorHandle func(c *gin.Context, err error) bool) {
	if err := fn(); err != nil {
		if handled := errorHandle(c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}

func GetUintFromPath(c *gin.Context, param string) (uint, error) {
	idS := c.Param(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func GetIntFromPath(c *gin.Context, param string) (int, error) {
	idS := c.Param(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetStringFromPath(c *gin.Context, param string) string {
	return c.Param(param)
}

func GetIntFromQuery(c *gin.Context, param string) (int, error) {
	idS := c.Query(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetStringFromQuery(c *gin.Context, param string) string {
	return c.Query(param)
}

func GetIntFromForm(c *gin.Context, param string) (int, error) {
	idS := c.PostForm(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}

	return id, nil
}
