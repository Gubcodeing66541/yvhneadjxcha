package Common

import (
	"github.com/gin-gonic/gin"
	Common2 "server/App/Common"
	"server/App/Sdk"
)

type Qiniu struct{}

func (Qiniu) QiniuToken(c *gin.Context) {
	token := Sdk.QiNiu{}.GetUpToken()
	Common2.ApiResponse{}.Success(c, "OK", gin.H{"token": token})
}
