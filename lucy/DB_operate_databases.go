package lucy

import (
	"database/sql"
	"fmt"
)
//打开数据库
func DB_open_database_lucy(SQL_name string,database_name string)(db *sql.DB,bool1 bool) {
	db,_= sql.Open(SQL_name,"root:@tcp(localhost:3306)/"+database_name)
	db.SetConnMaxLifetime(1000)
	err:=db.Ping()
	if err!=nil{
		return db,false
	}else{
		return db,true
	}
}
//插入数据
func DB_updateDB(db *sql.DB,statement string) (bool) {
	_, err := db.Exec(statement)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	return true
}
//删除数据
func DB_deleted_data(db *sql.DB,statement string) (bool) {
	_, err := db.Exec(statement)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	return true
}
//更改数据
func DB_change_data(db *sql.DB,statement string) (bool) {
	_, err := db.Exec(statement)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	return true
}