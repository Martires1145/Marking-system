package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"marking/model"
	"marking/response"
	"time"
)

// CreateExam
//
//	@Summary	创建考试
//	@Tags		Exam
//	@Param		exam			body	model.Exam	true	"考试信息"
//	@Param		Authorization	header	string		true	"token"
//	@Router		/api/exam/new [post]
func CreateExam(c *gin.Context) {
	var exam model.Exam
	err := c.ShouldBindJSON(&exam)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}
	exam.CreateTime = uint(time.Now().Unix())

	err = exam.Save()
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// DeleteExam
//
//	@Summary	删除考试
//	@Tags		Exam
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/exam/delete [post]
func DeleteExam(c *gin.Context) {
	examID := c.PostForm("examID")

	err := model.DeleteExam(examID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// ReviseExam
//
//	@Summary	修改考试信息
//	@Tags		Exam
//	@Param		exam			body	model.Exam	true	"考试信息"
//	@Param		Authorization	header	string		true	"token"
//	@Router		/api/exam/re [post]
func ReviseExam(c *gin.Context) {
	var exam model.Exam
	err := c.ShouldBindJSON(&exam)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	err = exam.Update()
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// ExamInfo
//
//	@Summary	获取考试信息
//	@Tags		Exam
//	@Param		examIDs			body	[]uint	true	"考试ID列表"
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/exam/info [post]
func ExamInfo(c *gin.Context) {
	var examIDs []uint
	err := c.ShouldBind(&examIDs)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	eInfos, err := model.GetExams(examIDs)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 400)
		return
	}

	response.Success(c.Writer, "success", eInfos)
}

// AllExam
//
//	@Summary	获取所有考试ID
//	@Tags		Exam
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/exam/all [post]
func AllExam(c *gin.Context) {
	ids, err := model.GetAllExam()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 400)
		return
	}

	response.Success(c.Writer, "success", ids)
}
