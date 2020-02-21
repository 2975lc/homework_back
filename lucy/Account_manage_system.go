package lucy

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)
//账号查重
func ACCOUNT_Find_account(db *sql.DB,type_bd string, befind_account string,befind_password string) (int) {
	tmp, err := db.Query("select * from "+type_bd+";")
	if err != nil {
		fmt.Println(err.Error())
		return 0                                 //异常
	}
	for ;tmp.Next(); {
		var account string
		var password string
		var can_use int
		var vip int
		var root int
		err:=tmp.Scan(&account,&password,&can_use,&vip,&root)
		if err != nil {
			fmt.Println(err.Error())
			return 0                           //异常
		}
		if befind_account ==account{
			if befind_password==password {
				return 1                       //登录成功
			} else{
				return 2                       //密码错误
			}
		}
	}
	return 3                                   //无重复
}
//登录
func ACCOUNT_Register_account(router *gin.Engine,db *sql.DB)  {
	router.GET("/register/:account/:password", func(c *gin.Context) {
		account:=c.Param("account")
		password:=c.Param("password")
		i:=ACCOUNT_Find_account(db,"account",account,password)
		state_s:="state"
		if i == 0{
			state:=0//异常
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i == 1{
			state:=1//成功
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i == 2{
			state:=2//密码错误
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i == 3{
			state:=3//未注册
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}
	})
}
//注册
func ACCOUNT_Login_account(router *gin.Engine,db *sql.DB) {
	router.GET("/login/:account/:password", func(c *gin.Context) {
		account:=c.Param("account")
		password:=c.Param("password")
		i:=ACCOUNT_Find_account(db,"account",account,password)
		state_s:="state"
		if i==0{
			state:=0//异常
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i==1{
			state:=1//已注册
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i==2{
			state:=2//已注册
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i==3{
			statement:="INSERT INTO account(account,password,can_use,vip,root) VALUES ('"+account+"','"+password+"',1,0,0);"
			if DB_updateDB(db,statement){
				state:=3//成功
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}else{
				state:=4//失败
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}
		}
	})
}
//更改密码
func ACCOUNT_Change_password(router *gin.Engine,db *sql.DB)  {
	router.GET("/change_password/:account/:password/:new_password", func(c *gin.Context) {
		account:=c.Param("account")
		password:=c.Param("password")
		new_password:=c.Param("new_password")
		i:=ACCOUNT_Find_account(db,"account",account,password)
		state_s:="state"
		if i==0{
			state:=0//失败
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i == 1 {
			statement:="UPDATE account SET password='"+new_password+"' WHERE account = '"+account+"';"
			if DB_change_data(db,statement) {
				state:=1//成功
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}else {
				state:=4//异常
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}
		}else if i == 2{
			state:=2//原密码错误
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i ==3 {
			state:=3//该账号不存在
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}
	})
}
//注销账号
func ACCOUNT_Logout_account(router *gin.Engine,db *sql.DB)  {
	router.GET("/logout/:account/:password", func(c *gin.Context) {
		account:=c.Param("account")
		password:=c.Param("password")
		i:=ACCOUNT_Find_account(db,"account",account,password)
		state_s:="state"
		if i==0{
			state:=0//失败
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i == 1 {
			statement:="DELETE FROM account WHERE account = '"+account+"';"
			if DB_deleted_data(db,statement) {
				state:=1//成功
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}else {
				state:=4//异常
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}
		}else if i==2{
			state:=2//密码错误
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else if i == 3 {
			state:=3//账号不存在
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}
	})
}
//账号管理系统
func ACCOUNT_manage_system(router *gin.Engine,db *sql.DB){
	ACCOUNT_Login_account(router,db)
	ACCOUNT_Register_account(router,db)
	ACCOUNT_Change_password(router,db)
	ACCOUNT_Logout_account(router,db)
}