package lucy

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

//百度搜索
func REDIRECT_baidu(router *gin.Engine)  {
	router.GET("/baidu=:要搜索的内容", func(c *gin.Context) {
	some:=c.Param("要搜索的内容")
	ht:="https://www.baidu.com/s?wd="+some
	c.Redirect(http.StatusMovedPermanently, ht)
	})
}
//重定向管理系统
func REDIRECT_manage_system(router *gin.Engine) {
	REDIRECT_baidu(router)
}