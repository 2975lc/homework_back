package lucy

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)
//书目查重
func BOOK_find_book(db *sql.DB,type_bd string,this_possessor string,befind_book string,befind_author string) int {
	tmp, err := db.Query("select * from "+type_bd+";")
	if err != nil {
		fmt.Println(err.Error())
		return 0                                 //异常
	}
	for ;tmp.Next(); {
		var possessor string
		var book_name string
		var author string
		var sort string
		var price float32
		err:=tmp.Scan(&possessor,&book_name,&author,&sort,&price)
		if err != nil {
			fmt.Println(err.Error())
			return 0                           //异常
		}
		if possessor ==this_possessor{
			if befind_book ==book_name{
				if befind_author==author {
					tmp.Close()
					return 1                       //存在相同的书
				} else{
					return 2                       //存在作者不同的相同的书
				}
			}
		}
	}
	return 3                                   //该书未在书架中找到
}
//添加书目
func BOOK_add_book(router *gin.Engine,db *sql.DB) {
	router.GET("/add_book/:possessor/:book_name/:author/:sort/:price", func(c *gin.Context) {
		possessor:=c.Param("possessor")
		book_name:=c.Param("book_name")
		author:=c.Param("author")
		sort:=c.Param("sort")
		price:=c.Param("price")
		i:=BOOK_find_book(db,"bookrack",possessor,book_name,author)
		c.String(http.StatusOK, "[\n{\n")
		state_s:="state"
		if i==0{
			state:=0//错误
			c.String(http.StatusOK, "%q:%d,\n",state_s,state)
		}else if i==1{
			state:=1//已存在
			c.String(http.StatusOK, "%q:%d,\n",state_s,state)
		}else{
			statement:="INSERT INTO bookrack(possessor,book_name,author,sort,price) VALUES ('"+possessor+"','"+book_name+"','"+author+"','"+sort+"',"+price+");"
			if DB_updateDB(db,statement){
				state:=2//添加完成
				c.String(http.StatusOK, "%q:%d,\n",state_s,state)
			}else{
				state:=3//添加失败
				c.String(http.StatusOK, "%q:%d,\n",state_s,state)
			}
		}
		possessor_s:="bookrack_possessor"
		c.String(http.StatusOK, "%q:%q,\n",possessor_s,possessor)
		book_name_s:="bookrack_book_name"
		c.String(http.StatusOK, "%q:%q,\n",book_name_s,book_name)
		author_s:="bookrack_author"
		c.String(http.StatusOK, "%q:%q,\n",author_s,author)
		sort_s:="bookrack_sort"
		c.String(http.StatusOK, "%q:%q,\n",sort_s,sort)
		price_s:="bookrack_price"
		c.String(http.StatusOK, "%q:%q\n",price_s,price)
		c.String(http.StatusOK, "}\n]\n")
	})
}
//修改书目
func BOOK_change_book(router *gin.Engine,db *sql.DB)  {
	router.GET("/change_book/:possessor/:book_name/:author/:new_book_name/:new_author/:new_sort/:new_price", func(c *gin.Context) {
		possessor:=c.Param("possessor")
		book_name:=c.Param("book_name")
		author:=c.Param("author")
		new_book_name:=c.Param("new_book_name")
		new_author:=c.Param("new_author")
		new_sort:=c.Param("new_sort")
		new_price:=c.Param("new_price")
		i:=BOOK_find_book(db,"bookrack",possessor,book_name,author)
		c.String(http.StatusOK, "[\n{\n")
		state_s:="state"
		if i==0{
			state:=0//修改异常
			c.String(http.StatusOK, "%q:%d,\n",state_s,state)
		}else if i == 1 {
			statement:="UPDATE bookrack SET book_name='"+new_book_name+"', author='"+new_author+"', sort='"+new_sort+"', price='"+new_price+"' WHERE possessor = '"+possessor+"' and book_name = '"+book_name+"' and author = '"+author+"';"
			if DB_change_data(db,statement) {
				state:=1//修改完成
				c.String(http.StatusOK, "%q:%d,\n",state_s,state)
			}else {
				state:=3//修改错误
				c.String(http.StatusOK, "%q:%d,\n",state_s,state)
			}
		}else{
			state:=2//该书不存在
			c.String(http.StatusOK, "%q:%d,\n",state_s,state)
		}
		possessor_s:="bookrack_possessor"
		c.String(http.StatusOK, "%q:%q,\n",possessor_s,possessor)
		book_name_s:="bookrack_book_name"
		c.String(http.StatusOK, "%q:%q,\n",book_name_s,new_book_name)
		author_s:="bookrack_author"
		c.String(http.StatusOK, "%q:%q,\n",author_s,new_author)
		sort_s:="bookrack_sort"
		c.String(http.StatusOK, "%q:%q,\n",sort_s,new_sort)
		price_s:="bookrack_price"
		c.String(http.StatusOK, "%q:%q\n",price_s,new_price)
		c.String(http.StatusOK, "}\n]\n")
	})
}
//展示所有书目
func BOOK_show_books(router *gin.Engine, db *sql.DB) {
	router.GET("/show_books/:possessor", func(c *gin.Context) {
		this_possessor:=c.Param("possessor")
		tmp, err := db.Query("select * from bookrack;")
		if err != nil {
			fmt.Println(err.Error())
		}else {
			c.String(http.StatusOK, "[\n")
			tm:=0
			for ;tmp.Next(); {
				var possessor string
				var book_name string
				var author string
				var sort string
				var price float64
				err:=tmp.Scan(&possessor,&book_name,&author,&sort,&price)
				if err != nil {
					fmt.Println(err.Error())
				}else {
					if possessor ==this_possessor{
						if tm != 0 {
							c.String(http.StatusOK, ",")
						}
						tm=1
						c.String(http.StatusOK, "{\n")
						possessor_s:="bookrack_possessor"
						c.String(http.StatusOK, "%q:%q,\n",possessor_s,possessor)
						book_name_s:="bookrack_book_name"
						c.String(http.StatusOK, "%q:%q,\n",book_name_s,book_name)
						author_s:="bookrack_author"
						c.String(http.StatusOK, "%q:%q,\n",author_s,author)
						sort_s:="bookrack_sort"
						c.String(http.StatusOK, "%q:%q,\n",sort_s,sort)
						price_s:="bookrack_price"
						c.String(http.StatusOK, "%q:%q,\n",price_s,price)
						c.String(http.StatusOK, "}\n")
					}
				}
			}
			c.String(http.StatusOK, "]\n")
		}
		tmp.Close()
	})
}
//删除书目
func BOOK_delete_book(router *gin.Engine,db *sql.DB) {
	router.GET("/delete_book/:possessor/:book_name/:author", func(c *gin.Context) {
		possessor:=c.Param("possessor")
		book_name:=c.Param("book_name")
		author:=c.Param("author")
		i:=BOOK_find_book(db,"bookrack",possessor,book_name,author)
		c.String(http.StatusOK, "[\n{\n")
		state_s:="state"
		if i==0{
			state:=0//异常
			c.String(http.StatusOK, "%q:%d,\n",state_s,state)
		}else if i==1{
			statement:="DELETE FROM bookrack WHERE possessor = '"+possessor+"' and book_name = '"+book_name+"' and author = '"+author+"';"
			if DB_deleted_data(db,statement) {
				state:=1//成功
				c.String(http.StatusOK, "%q:%d,\n",state_s,state)
			}else {
				state:=3//错误
				c.String(http.StatusOK, "%q:%d,\n",state_s,state)
			}
		}else{
			state:=2//书不存在
			c.String(http.StatusOK, "%q:%d,\n",state_s,state)
		}
		possessor_s:="bookrack_possessor"
		c.String(http.StatusOK, "%q:%q,\n",possessor_s,possessor)
		book_name_s:="bookrack_book_name"
		c.String(http.StatusOK, "%q:%q,\n",book_name_s,book_name)
		author_s:="bookrack_author"
		c.String(http.StatusOK, "%q:%q,\n",author_s,author)
		c.String(http.StatusOK, "}\n]\n")
	})
}
//书架管理系统
func BOOK_manage_system(router *gin.Engine,db *sql.DB) {
	BOOK_add_book(router,db)
	BOOK_change_book(router,db)
	BOOK_show_books(router,db)
	BOOK_delete_book(router,db)
}