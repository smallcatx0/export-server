package valid

import (
	"export-server/pkg/exception"

	"github.com/gin-gonic/gin"
)

type CustomValidor interface {
	Valid() error
}

func customCheck(param interface{}) error {
	if validor, ok := param.(CustomValidor); ok {
		err := validor.Valid()
		if err != nil {
			return exception.NewException(exception.Fail, err.Error())
		}
	}
	return nil
}

func BindAndCheck(c *gin.Context, param interface{}) error {
	err := c.ShouldBindJSON(param)
	if err != nil {
		return exception.NewException(exception.Fail, err.Error())
	}
	// 自定义验证规则
	return customCheck(param)
}

func BindQAndCheck(c *gin.Context, param interface{}) error {
	err := c.ShouldBindQuery(param)
	if err != nil {
		return exception.NewException(exception.Fail, err.Error())
	}
	// 自定义验证规则
	return customCheck(param)
}
