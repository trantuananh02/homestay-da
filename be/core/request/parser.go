package http_request

import "github.com/gin-gonic/gin"

func BindFormData(c *gin.Context, b interface{}) error {
	err := c.ShouldBind(b)
	if err != nil {
		return err
	}
	return nil
}

func BindQueryString(c *gin.Context, b interface{}) error {
	err := c.ShouldBind(b)
	if err != nil {
		return err
	}
	return nil
}

func BindUri(c *gin.Context, b interface{}) error {
	err := c.ShouldBindUri(b)
	if err != nil {
		return err
	}
	return nil
}

func BindBodyJson(c *gin.Context, b interface{}) error {
	err := c.ShouldBindJSON(b)
	if err != nil {
		return err
	}
	return nil
}
