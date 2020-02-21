package lucy

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

//查找新问题序号
func QA_find_sequence(db *sql.DB,type_bd string)(r int) {
	tmp, err := db.Query("select * from "+type_bd+";")
	r=0
	if err != nil {
		fmt.Println(err.Error())
		return -1                                 //异常
	}
	for ;tmp.Next(); {
		var quizzer string
		var quiz string
		var quiz_describe string
		var sequence int
		var attention int
		var exist int
		err:= tmp.Scan(&quizzer,&quiz, &quiz_describe, &sequence,&attention,&exist)
		if err != nil {
			fmt.Println(err.Error())
			return -1                            //异常
		}
		r=sequence
	}
	return r+1
}
//提问
func QA_quiz(router *gin.Engine,db *sql.DB){
	router.GET("/quiz/:account/:quiz/:message", func(c *gin.Context) {
		account:=c.Param("account")
		quiz:=c.Param("quiz")
		message:=c.Param("message")
		r:= QA_find_sequence(db,"issue")
		account_s:="account"
		if r!=-1{
			statement:="INSERT INTO issue (quizzer,quiz,quiz_describe,sequence,attention,exist) VALUES ('"+account+"','"+quiz+"','"+message+"','"+strconv.Itoa(r)+"',0,1);"
			if DB_updateDB(db,statement){
				state:=1//成功
				c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
			}else{
				state:=2//异常
				c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
			}
		}else {
			state:=0//失败
			c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
		}
	})
}
//判断提问是否存在
func QA_judge_quiz_exist(db *sql.DB,type_bd string,found_sequence int)(exist int)  {
	tmp, err := db.Query("select * from "+type_bd+";")
	if err != nil {
		fmt.Println(err.Error())
		return -1                                 //异常
	}
	for ;tmp.Next(); {
		var quizzer string
		var quiz string
		var quiz_describe string
		var sequence int
		var attention int
		var exist int
		err:= tmp.Scan(&quizzer,&quiz, &quiz_describe, &sequence,&attention,&exist)
		if err != nil {
			fmt.Println(err.Error())
			return -1                            //异常
		}
		if found_sequence==sequence {
			return exist
		}
	}
	return -1
}
//查找新回答序号
func QA_find_floor(db *sql.DB,type_bd string,found_sequence int)(f int) {
	tmp, err := db.Query("select * from "+type_bd+";")
	f =0
	if err != nil {
		fmt.Println(err.Error())
		return -1                                 //异常
	}
	for ;tmp.Next(); {
		var possessor string
		var answer string
		var sequence int
		var floor int
		var praise int
		var exist int
		err:= tmp.Scan(&possessor, &answer, &sequence, &floor,&praise,&exist)
		if err != nil {
			fmt.Println(err.Error())
			return -1                            //异常
		}
		if found_sequence==sequence {
			if f<floor {
				f=floor
			}
		}
	}
	return f +1
}
//解答
func QA_resolve(router *gin.Engine,db *sql.DB){
	router.GET("/resolve/:account/:answer/:sequence", func(c *gin.Context) {
		account:=c.Param("account")
		answer:=c.Param("answer")
		sequence, _ :=strconv.Atoi(c.Param("sequence"))
		exist:=QA_judge_quiz_exist(db,"issue",sequence)
		account_s:="account"
		if exist==1 {
			f := QA_find_floor(db,"solution",sequence)
			if f !=-1{
				statement:="INSERT INTO solution (possessor,answer,sequence,floor,praise,exist) VALUES ('"+account+"','"+answer+"','"+strconv.Itoa(sequence)+"','"+strconv.Itoa(f)+"',0,1);"
				if DB_updateDB(db,statement){
					state:=1//成功
					c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
				}else{
					state:=2//异常
					c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
				}
			}else {
				state:=0//故障
				c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
			}
		}else if exist == -1 {
			state:=0//故障
			c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
		}else if exist == 0 {
			state:=0//失败
			c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
		}
	})
}
//查找新评论序号
func QA_find_tier(db *sql.DB,type_bd string,found_sequence int,found_floor int)(t int) {
	tmp, err := db.Query("select * from "+type_bd+";")
	t =0
	if err != nil {
		fmt.Println(err.Error())
		return -1                                 //异常
	}
	for ;tmp.Next(); {
		var possessor string
		var discuss string
		var sequence int
		var floor int
		var tier int
		var little_tier int
		var star int
		var praise int
		var exist int
		err:= tmp.Scan(&possessor, &discuss, &sequence, &floor,&tier,&little_tier,&star,&praise,&exist)
		if err != nil {
			fmt.Println(err.Error())
			return -1                            //异常
		}
		if found_sequence==sequence {
			if found_floor==floor {
				if t < tier {
					t =tier
				}
			}
		}
	}
	return t +1
}
//评论
func QA_review(router *gin.Engine,db *sql.DB){
	router.GET("/review/:account/:discuss/:sequence/:floor", func(c *gin.Context) {
		account:=c.Param("account")
		discuss:=c.Param("discuss")
		sequence, _ :=strconv.Atoi(c.Param("sequence"))
		floor, _ :=strconv.Atoi(c.Param("floor"))
		t:=QA_find_tier(db,"comment",sequence,floor)
		account_s:="account"
		if t !=-1{
			statement:="INSERT INTO comment (possessor,discuss,sequence,floor,tier,little_tier,star,praise,exist) VALUES ('"+account+"','"+discuss+"','"+strconv.Itoa(sequence)+"','"+strconv.Itoa(floor)+"','"+strconv.Itoa(t)+"',0,0,0,1);"
			if DB_updateDB(db,statement){
				state:=1//成功
				c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
			}else{
				state:=2//异常
				c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
			}
		}else {
			state:=3//故障
			c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
		}
	})
}
//查找新回复序号
func QA_find_little_tier(db *sql.DB,type_bd string,found_sequence int,found_floor int,found_tier int)(lt int) {
	tmp, err := db.Query("select * from "+type_bd+";")
	lt =0
	if err != nil {
		fmt.Println(err.Error())
		return -1                                 //异常
	}
	for ;tmp.Next(); {
		var possessor string
		var discuss string
		var sequence int
		var floor int
		var tier int
		var little_tier int
		var star int
		var praise int
		var exist int
		err:= tmp.Scan(&possessor, &discuss, &sequence, &floor,&tier,&little_tier,&star,&praise,&exist)
		if err != nil {
			fmt.Println(err.Error())
			return -1                            //异常
		}
		if found_sequence==sequence {
			if found_floor==floor {
				if found_tier==tier {
					if lt < little_tier {
						lt =little_tier
					}
				}
			}
		}
	}
	return lt +1
}
//回复
func QA_reply(router *gin.Engine,db *sql.DB){
	router.GET("/reply/:account/:discuss/:sequence/:floor/:tier", func(c *gin.Context) {
		account:=c.Param("account")
		discuss:=c.Param("discuss")
		sequence, _ :=strconv.Atoi(c.Param("sequence"))
		floor, _ :=strconv.Atoi(c.Param("floor"))
		tier, _ :=strconv.Atoi(c.Param("tier"))
		lt :=QA_find_little_tier(db,"comment",sequence,floor,tier)
		account_s:="account"
		if lt !=-1{
			statement:="INSERT INTO comment (possessor,discuss,sequence,floor,tier,little_tier,star,praise,exist) VALUES ('"+account+"','"+discuss+"','"+strconv.Itoa(sequence)+"','"+strconv.Itoa(floor)+"','"+strconv.Itoa(tier)+"','"+strconv.Itoa(lt)+"',0,0,1);"
			if DB_updateDB(db,statement){
				state:=1//成功
				c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
			}else{
				state:=2//异常
				c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
			}
		}else {
			state:=3//故障
			c.String(http.StatusOK, "[{%q:%d}]",account_s,state)
		}
	})
}
//删除问答
func QA_delete_system(router *gin.Engine,db *sql.DB)  {
	router.GET("/delete_QA/:sequence/:floor/:tier/:little_tier", func(c *gin.Context) {
		sequence:=c.Param("sequence")
		floor:=c.Param("floor")
		tier:=c.Param("tier")
		little_tier:=c.Param("little_tier")
		tme_floor,_:=strconv.Atoi(floor)
		tme_tier,_:=strconv.Atoi(tier)
		var statement string
		state_s:="state"
		if tme_floor == 0 {
			statement="UPDATE issue SET exist='0' WHERE sequence = '"+sequence+"';"
		}else if tme_tier == 0 {
			statement="UPDATE solution SET exist='0' WHERE sequence = '"+sequence+"' and floor = '"+floor+"';"
		}else {
			statement="UPDATE comment SET exist='0' WHERE sequence = '"+sequence+"' and floor = '"+floor+"' and tier = '"+tier+"' and little_tier='"+little_tier+"';"
		}
		if DB_change_data(db,statement) {
			state:=1//成功
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else {
			state:=2//异常
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}
	})
}
//查找问题关注度
func QA_find_attention(db *sql.DB,type_bd string,found_sequence int)(attention int) {
	tmp, err := db.Query("select * from "+type_bd+";")
	if err != nil {
		fmt.Println(err.Error())
		return -1                                 //异常
	}
	for ;tmp.Next(); {
		var quizzer string
		var quiz_describe string
		var sequence int
		var attention int
		var exist int
		err:= tmp.Scan(&quizzer, &quiz_describe, &sequence,&attention,&exist)
		if err != nil {
			fmt.Println(err.Error())
			return -1                            //异常
		}
		if sequence==found_sequence {
			return attention
		}
	}
	return -1                                 //未找到该问题
}
//查找回答赞数
func QA_find_resolve_praise(db *sql.DB,type_bd string,found_sequence int,found_floor int)(praise int) {
	tmp, err := db.Query("select * from "+type_bd+";")
	if err != nil {
		fmt.Println(err.Error())
		return -1                                 //异常
	}
	for ;tmp.Next(); {
		var possessor string
		var message string
		var sequence int
		var floor int
		var praise int
		var exist int
		err:= tmp.Scan(&possessor, &message, &sequence, &floor,&praise,&exist)
		if err != nil {
			fmt.Println(err.Error())
			return -1                            //异常
		}
		if sequence==found_sequence {
			if found_floor==floor {
				return praise
			}
		}
	}
	return -1                                 //未找到该提问
}
//查找回复赞数
func QA_find_review_praise(db *sql.DB,type_bd string,found_sequence int,found_floor int,found_tier int,found_little_tier int)(praise int) {
	tmp, err := db.Query("select * from "+type_bd+";")
	if err != nil {
		fmt.Println(err.Error())
		return -1                                 //异常
	}
	for ;tmp.Next(); {
		var possessor string
		var message string
		var sequence int
		var floor int
		var tier int
		var little_tier int
		var star int
		var praise int
		var exist int
		err:= tmp.Scan(&possessor, &message, &sequence, &floor,&tier,&little_tier,&star,&praise,&exist)
		if err != nil {
			fmt.Println(err.Error())
			return -1                            //异常
		}
		if sequence==found_sequence {
			if found_floor==floor {
				if found_tier==tier {
					if found_little_tier==little_tier {
						return praise
					}
				}
			}
		}
	}
	return -1                                 //未找到该回复
}
//点赞or关注系统
func QA_praise_system(router *gin.Engine,db *sql.DB)  {
	router.GET("/praise_QA/:add_or_minus/:sequence/:floor/:tier/:little_tier", func(c *gin.Context) {
		add_or_minus:=c.Param("add_or_minus")
		change:=1
		if add_or_minus=="-" {
			change=-1
		}
		sequence:=c.Param("sequence")
		floor:=c.Param("floor")
		tier:=c.Param("tier")
		little_tier:=c.Param("little_tier")
		tme_sequence, _ :=strconv.Atoi(sequence)
		tme_floor,_:=strconv.Atoi(floor)
		tme_tier,_:=strconv.Atoi(tier)
		tme_little_tier,_:=strconv.Atoi(little_tier)
		var statement string
		state_s:="state"
		if tme_floor == 0 {
			attention:=QA_find_attention(db,"issue",tme_sequence)
			if attention!=-1 {
				statement="UPDATE issue SET attention='"+strconv.Itoa(attention+change)+"' WHERE sequence = '"+sequence+"';"
			}else {
				state:=0//失败
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}
		}else if tme_tier == 0 {
			praise:=QA_find_resolve_praise(db,"solution",tme_sequence,tme_floor)
			if praise != -1 {
				statement="UPDATE solution SET praise='"+strconv.Itoa(praise+change)+"' WHERE sequence = '"+sequence+"' and floor = '"+floor+"';"
			}else {
				state:=0//失败
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}
		}else {
			praise:=QA_find_review_praise(db,"comment",tme_sequence,tme_floor,tme_tier,tme_little_tier)
			if praise != -1 {
				statement="UPDATE comment SET praise='"+strconv.Itoa(praise+change)+"' WHERE sequence = '"+sequence+"' and floor = '"+floor+"' and tier = '"+tier+"' and little_tier='"+little_tier+"';"
			}else {
				state:=0//失败
				c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
			}
		}
		if DB_change_data(db,statement) {
			state:=1//成功
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}else {
			state:=-1//异常
			c.String(http.StatusOK, "[{%q:%d}]",state_s,state)
		}
	})
}
//展示问题列表
func QA_show_quiz(router *gin.Engine,db *sql.DB)  {
	router.GET("/show_quiz", func(c *gin.Context) {
		type_bd:="issue"
		tmp, err := db.Query("select * from "+type_bd+";")
		if err != nil {
			fmt.Println(err.Error())
			c.String(http.StatusOK, "异常\n")
		}else {
			tm:=0
			c.String(http.StatusOK, "[\n")
			for ;tmp.Next(); {
				var quizzer string
				var quiz string
				var quiz_describe string
				var sequence int
				var attention int
				var exist int
				err:= tmp.Scan(&quizzer,&quiz,&quiz_describe, &sequence,&attention,&exist)
				if err != nil {
					fmt.Println(err.Error())
					c.String(http.StatusOK, "异常\n")
				}else {
					if tm != 0 {
						c.String(http.StatusOK,",")
					}
					tm=1
					c.String(http.StatusOK,"{\n")
					quizzer_s:="issue_quizzer"
					c.String(http.StatusOK,"%q:%q,\n",quizzer_s,quizzer)
					quiz_s:="issue_quiz"
					c.String(http.StatusOK,"%q:%q,\n",quiz_s,quiz)
					quiz_describe_s:="issue_quiz_describe"
					c.String(http.StatusOK,"%q:%q,\n",quiz_describe_s,quiz_describe)
					sequence_s:="issue_sequence"
					c.String(http.StatusOK,"%q:%d,\n",sequence_s,sequence)
					attention_s:="issue_attention"
					c.String(http.StatusOK,"%q:%d,\n",attention_s,attention)
					exist_s:="issue_exist"
					c.String(http.StatusOK,"%q:%d\n",exist_s,exist)
					c.String(http.StatusOK,"}\n")
				}
			}
			c.String(http.StatusOK, "]")
		}
	})
}
//展示解答列表
func QA_show_resolve(router *gin.Engine,db *sql.DB)  {
	router.GET("/show_resolve/:sequence", func(c *gin.Context) {
		found_sequence, _ :=strconv.Atoi(c.Param("sequence"))
		type_bd:="solution"
		tmp, err := db.Query("select * from "+type_bd+";")
		tm:=0
		if err != nil {
			fmt.Println(err.Error())
			c.String(http.StatusOK, "异常\n")
		}else {
			c.String(http.StatusOK,"[\n")
			for ;tmp.Next(); {
				var possessor string
				var answer string
				var sequence int
				var floor int
				var praise int
				var exist int
				err:= tmp.Scan(&possessor, &answer, &sequence, &floor,&praise,&exist)
				if err != nil {
					fmt.Println(err.Error())
					c.String(http.StatusOK, "异常\n")
				}else {
					if found_sequence==sequence {
						if tm != 0 {
							c.String(http.StatusOK,",")
						}
						tm=1
						c.String(http.StatusOK,"{\n")
						possessor_s:="solution_possessor"
						c.String(http.StatusOK,"%q:%q,\n",possessor_s,possessor)
						answer_s:="solution_answer"
						c.String(http.StatusOK,"%q:%q,\n",answer_s,answer)
						sequence_s:="solution_sequence"
						c.String(http.StatusOK,"%q:%d,\n",sequence_s,sequence)
						floor_s:="solution_floor"
						c.String(http.StatusOK,"%q:%d,\n",floor_s,floor)
						praise_s:="solution_praise"
						c.String(http.StatusOK,"%q:%d,\n",praise_s,praise)
						exist_s:="solution_exist"
						c.String(http.StatusOK,"%q:%d\n",exist_s,exist)
						c.String(http.StatusOK,"}\n")
					}
				}
			}
			c.String(http.StatusOK,"]\n")
		}
	})
}
//展示评论列表
func QA_show_review(router *gin.Engine,db *sql.DB)  {
	router.GET("/show_review/:sequence/:floor", func(c *gin.Context) {
		found_sequence, _ :=strconv.Atoi(c.Param("sequence"))
		found_floor, _ :=strconv.Atoi(c.Param("floor"))
		type_bd:="comment"
		tmp, err := db.Query("select * from "+type_bd+";")
		tm:=0
		if err != nil {
			fmt.Println(err.Error())
			c.String(http.StatusOK, "异常\n")
		}else {
			c.String(http.StatusOK,"[\n")
			for ;tmp.Next(); {
				var possessor string
				var discuss string
				var sequence int
				var floor int
				var tier int
				var little_tier int
				var star int
				var praise int
				var exist int
				err:= tmp.Scan(&possessor, &discuss, &sequence, &floor,&tier,&little_tier,&star,&praise,&exist)
				if err != nil {
					fmt.Println(err.Error())
					c.String(http.StatusOK, "异常\n")
				}else{
					if found_sequence==sequence {
						if found_floor==floor {
							if little_tier==0{
								if tm != 0 {
									c.String(http.StatusOK,",")
								}
								tm=1
								c.String(http.StatusOK,"{\n")
								possessor_s:="comment_possessor"
								c.String(http.StatusOK,"%q:%q,\n",possessor_s,possessor)
								discuss_s:="comment_discuss"
								c.String(http.StatusOK,"%q:%q,\n",discuss_s,discuss)
								sequence_s:="comment_sequence"
								c.String(http.StatusOK,"%q:%d,\n",sequence_s,sequence)
								floor_s:="comment_floor"
								c.String(http.StatusOK,"%q:%d,\n",floor_s,floor)
								tier_s:="comment_tier"
								c.String(http.StatusOK,"%q:%d,\n",tier_s,tier)
								little_tier_s:="comment_little_tier"
								c.String(http.StatusOK,"%q:%d,\n",little_tier_s,little_tier)
								star_s:="comment_star"
								c.String(http.StatusOK,"%q:%d,\n",star_s,star)
								praise_s:="comment_praise"
								c.String(http.StatusOK,"%q:%d,\n",praise_s,praise)
								exist_s:="comment_exist"
								c.String(http.StatusOK,"%q:%d\n",exist_s,exist)
								c.String(http.StatusOK,"}\n")
							}
						}
					}
				}
			}
			c.String(http.StatusOK,"]\n")
		}
	})
}
//问答管理系统
func QA_manage_system(router *gin.Engine,db *sql.DB){
	QA_quiz(router,db)
	QA_resolve(router,db)
	QA_review(router,db)
	QA_reply(router,db)
	QA_delete_system(router,db)
	QA_praise_system(router,db)
	QA_show_quiz(router,db)
	QA_show_resolve(router,db)
	QA_show_review(router,db)
}