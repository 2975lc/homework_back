package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"lucy"
)
func main (){
	f:=true
	SQL_name:="mysql"
	database_name:="lucy"
	var statement =[5]string{
		//账号表
		"create table account (account varchar(40),password varchar(60),can_use int(1),vip int(1),root int(1))",
		//书架表
		"create table bookrack (possessor varchar(40),book_name varchar(200),author varchar(200),sort varchar(200),price double(9,2))",
		//问题表
		"create table issue (quizzer varchar(40),quiz varchar(200),quiz_describe varchar(20000),sequence int(255),attention int(255),exist int(1))",
		//回答表
		"create table solution (possessor varchar(40),answer varchar(20000),sequence int(255),floor int(255),praise int(255),exist int(1))",
		//评论表
		"create table comment (possessor varchar(40),discuss varchar(20000),sequence int(255),floor int(255),tier int(255),little_tier int(255),star int(1),praise int(255),exist int(1))",
	}
	db,_:=lucy.DB_open_database_lucy(SQL_name,database_name)
	for x:=0;x<5;x++{
		_,err2:=db.Exec(statement[x])
		if err2 != nil {
			f=false
			fmt.Println("fail in create database")
			fmt.Println(err2.Error())
		}
	}
	if f{
		fmt.Println("create tables success!")
	}else {
		fmt.Println("create tables fail")
	}
}