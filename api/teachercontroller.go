package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"marking/model"
	"marking/response"
	"marking/util"
	"strconv"
)

// GetTaskMission
//
//	@Summary	获取改卷任务
//	@Tags		Teacher
//	@Param		Authorization	header		string	true	"token"
//	@Param		examID			formData	int		true	"考试ID"
//	@Router		/api/seisei/mission [post]
func GetTaskMission(c *gin.Context) {
	token := c.GetHeader("Authorization")
	examID := c.PostForm("examID")
	_, claim, _ := util.ParseToken(token)
	uid := claim.UID

	mission, err := model.GetMission(uid, examID)
	if err == sql.ErrNoRows {
		response.Success(c.Writer, "success", nil)
		return
	}

	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", mission)
}

func GetClassAnalysis(c *gin.Context) {

}

// MarkPaperPart
//
//	@Summary	为单题打分
//	@Tags		Teacher
//	@Param		partID			formData	int		true	"部分试卷ID"
//	@Param		mark			formData	int		true	"分数"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/seisei/mark [post]
func MarkPaperPart(c *gin.Context) {
	partID := c.PostForm("partID")
	mark, _ := strconv.Atoi(c.PostForm("mark"))

	err := model.MarkPaperPart(partID, mark)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// GetAnswerOfPaper
//
//	@Summary	获取试卷答案
//	@Tags		Teacher
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/seisei/as [post]
func GetAnswerOfPaper(c *gin.Context) {
	examID := c.PostForm("examID")

	answer, err := model.GetAnswer(examID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", answer)
}

// GetMissionPart
//
//	@Summary	按任务阅卷
//	@Tags		Teacher
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/seisei/ms/part [post]
func GetMissionPart(c *gin.Context) {
	examID := c.PostForm("examID")
	token := c.GetHeader("Authorization")
	_, claim, _ := util.ParseToken(token)
	uid := claim.UID

	paperPart, err := model.GetPaperPart(examID, uid)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", paperPart)
}

// FinishMission
//
//	@Summary	完成阅卷任务
//	@Tags		Teacher
//	@Param		examID			formData	int		true	"考试ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/seisei/ms/fm [post]
func FinishMission(c *gin.Context) {
	token := c.GetHeader("Authorization")
	examID := c.PostForm("examID")
	_, claim, _ := util.ParseToken(token)
	uid := claim.UID

	err := model.FinishMission(examID, uid)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// GetExamIn
//
//	@Summary	查看所在考试
//	@Tags		Teacher
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/seisei/ms/exam [post]
func GetExamIn(c *gin.Context) {
	token := c.GetHeader("Authorization")
	_, claim, _ := util.ParseToken(token)
	uid := claim.UID

	e, err := model.GetExamIn(uid)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", e)
}
