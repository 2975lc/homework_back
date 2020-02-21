package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"lucy"
)
func main (){
	router := gin.Default()
	SQL_name:="mysql"
	database_name:="lucy"
	db, _ := lucy.DB_open_database_lucy(SQL_name, database_name)
	lucy.REDIRECT_manage_system(router)
	lucy.ACCOUNT_manage_system(router,db)
	lucy.BOOK_manage_system(router,db)
	lucy.QA_manage_system(router,db)
	_ = router.Run(":2975")
}