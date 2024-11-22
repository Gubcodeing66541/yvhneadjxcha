package Common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"server/Base"
	"time"
)

type File struct {
}

// 获取头像图片
func (File) GetImg() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s/%dhead.jpg", Base.AppConfig.HeadImgUrl, r.Intn(4684))
}

func (File) CreateRename(c *gin.Context) {
	ApiResponse{}.Success(c, "获取rename", gin.H{})
}
