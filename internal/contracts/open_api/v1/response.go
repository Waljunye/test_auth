package v1

import "github.com/gin-gonic/gin"

func responseOk(ginCtx *gin.Context, payload interface{}) {
	ginCtx.JSON(200, payload)
}
