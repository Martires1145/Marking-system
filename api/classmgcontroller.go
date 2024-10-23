package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"marking/model"
	"marking/response"
	"marking/util"
	"strconv"
)

// QuitClass
//
//	@Summary	离开班级
//	@Tags		Class
//	@Param		classID			formData	int		true	"班级ID"
//	@Param		Authorization	header		string	true	"token"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/c/mg/quit [post]
func QuitClass(c *gin.Context) {
	var u model.User
	classID := c.PostForm("classID")
	token := c.GetHeader("Authorization")
	_, claim, _ := util.ParseToken(token)
	u.Uid = claim.UID
	u.Role = claim.Role

	err := u.QuitClass(classID)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("error when user quit class, err:%s", err.Error()), 500)
		return
	}

	response.Success(c.Writer, "success", nil)
}

// ClearStudent
//
//	@Summary	将学生请出班级
//	@Tags		Class
//	@Param		classID			formData	int		true	"班级ID"
//	@Param		studentID		formData	int		true	"班级ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/c/mg/clear [post]
func ClearStudent(c *gin.Context) {
	classID := c.PostForm("classID")
	studentID := c.PostForm("studentID")
	id, _ := strconv.Atoi(studentID)

	err := model.StudentQuitClass(uint(id), classID)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 500)
		return
	}

	response.Success(c.Writer, "success", nil)
}

// AddClass
//
//	@Summary	加入班级
//	@Tags		Class
//	@Param		classToken		formData	string	true	"班级验证码"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/c/mg/add [post]
func AddClass(c *gin.Context) {
	var u model.User
	token := c.GetHeader("Authorization")
	_, claim, _ := util.ParseToken(token)
	u.Role = claim.Role
	classToken := c.PostForm("classToken")
	u.Uid = claim.UID

	err := u.AddClass(classToken)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 500)
		return
	}

	response.Success(c.Writer, "success", nil)
}

// ClassScore
//
//	@Summary	查询班级某次考试的成绩
//	@Tags		Class
//	@Param		classID			formData	string	true	"班级ID"
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/c/mg/score [post]
func ClassScore(c *gin.Context) {
	classID := c.PostForm("classID")
	examID := c.PostForm("examID")

	classScore, err := model.GetClassScore(classID, examID)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 500)
		return
	}

	response.Success(c.Writer, "success", classScore)
}

// ClassHistoryScore
//
//	@Summary	查询班级最近n次的考试成绩
//	@Tags		Class
//	@Param		classID			formData	string	true	"班级ID"
//	@Param		n				formData	int		true	"查询的数目"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/c/mg/sh [post]
func ClassHistoryScore(c *gin.Context) {
	classID := c.PostForm("classID")
	n := c.PostForm("n")

	classScore, err := model.GetClassScoreHistory(classID, n)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 500)
		return
	}

	response.Success(c.Writer, "success", classScore)
}
