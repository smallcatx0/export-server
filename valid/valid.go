package valid

import (
	"export-server/pkg/exception"

	"github.com/gin-gonic/gin"
)

type CustomValidor interface {
	Valid() error
}

func BindAndCheck(c *gin.Context, param interface{}) error {
	err := c.ShouldBindJSON(param)
	if err != nil {
		return exception.NewException(exception.Fail, err.Error())
	}
	// 自定义验证规则
	if validor, ok := param.(CustomValidor); ok {
		err = validor.Valid()
		if err != nil {
			return exception.NewException(exception.Fail, err.Error())
		}
	}
	return nil
}
