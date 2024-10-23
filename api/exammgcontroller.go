package api

import (
	"github.com/gin-gonic/gin"
	"marking/model"
	"marking/response"
)

// ClassQuitExam
//
//	@Summary	班级退出考试
//	@Tags		Exam
//	@Param		classID			formData	int		true	"班级ID"
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/exam/mg/quit [post]
func ClassQuitExam(c *gin.Context) {
	classID := c.PostForm("classID")
	examID := c.PostForm("examID")

	err := model.QuitExam(classID, examID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// ClassJoinExam
//
//	@Summary	班级加入考试
//	@Tags		Exam
//	@Param		classID			formData	int		true	"班级ID"
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/exam/mg/join [post]
func ClassJoinExam(c *gin.Context) {
	classID := c.PostForm("classID")
	examID := c.PostForm("examID")

	err := model.JoinExam(classID, examID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// GetExamNum
//
//	@Summary	获取考试相关的数据
//	@Tags		Exam
//	@Param		type			formData	string	true	"数据类型 student teacher task"
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/exam/mg/num [post]
func GetExamNum(c *gin.Context) {
	examID := c.PostForm("examID")
	tpe := c.PostForm("type")

	num, err := model.GetExamNum(examID, tpe)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", gin.H{tpe: num})
}

// DistributePaper
//
//	@Summary	分发考试任务
//	@Tags		Exam
//	@Param		cut				body	model.CutMg	true	"改卷任务分段点数组"
//	@Param		Authorization	header	string		true	"token"
//	@Router		/api/exam/mg/plan [post]
func DistributePaper(c *gin.Context) {
	var cut model.CutMg
	err := c.ShouldBind(&cut)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 400)
		return
	}

	err = model.SavePaperMission(cut.ExamID, cut.Cut)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// GetScore
//
//	@Summary	获取某一考试所有人的考试成绩
//	@Tags		Exam
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/exam/mg/score [post]
func GetScore(c *gin.Context) {
	examID := c.PostForm("examID")

	scores, err := model.GetAllScore(examID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", scores)
}
