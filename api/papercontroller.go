package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"marking/model"
	"marking/response"
	"sync"
)

var once sync.Once

// NewPaper
//
//	@Summary	添加试卷
//	@Tags		Paper
//	@Param		papers			body	[]model.PaperJson	true	"试卷信息列表"
//	@Param		Authorization	header	string				true	"token"
//	@Router		/api/paper/new [post]
func NewPaper(c *gin.Context) {
	var paper []model.PaperJson
	err := c.ShouldBindJSON(&paper)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	err = model.SavePapers(paper)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}

	once.Do(func() {
		go model.OcrAndJudge()
	})

	response.Success(c.Writer, "success", nil)
}

// TotalScore
//
//	@Summary	更新试卷的总分
//	@Tags		Paper
//	@Param		paperIDs		body	[]int	true	"试卷ID"
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/paper/sum [post]
func TotalScore(c *gin.Context) {
	var paperID []int
	err := c.ShouldBind(&paperID)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	err = model.UpdateTotalScore(paperID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// DeletePaper
//
//	@Summary	删除试卷
//	@Tags		Paper
//	@Param		paperID			formData	int		true	"试卷ID"
//	@Param		Authorization	header		string	true	"token"
//	@Router		/api/paper/delete [post]
func DeletePaper(c *gin.Context) {
	paperID := c.PostForm("paperID")

	err := model.DeletePaper(paperID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// PaperInfos
//
//	@Summary	获取试卷信息
//	@Tags		Paper
//	@Param		paperIDs		body	[]int	true	"试卷ID列表"
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/paper/info [post]
func PaperInfos(c *gin.Context) {
	var paperID []int
	err := c.ShouldBind(&paperID)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	papers, err := model.PaperInfo(paperID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", papers)
}

// NewAnswer
//
//	@Summary	添加试卷答案
//	@Tags		Paper
//	@Param		answers			body	model.AnswerJson	true	"试卷信息列表"
//	@Param		Authorization	header	string				true	"token"
//	@Router		/api/paper/na [post]
func NewAnswer(c *gin.Context) {
	var answers model.AnswerJson
	err := c.ShouldBind(&answers)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	err = model.SaveAnswer(&answers)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// AllPaper
//
//	@Summary	获取所有试卷的ID
//	@Tags		Paper
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/paper/all [post]
func AllPaper(c *gin.Context) {
	ids, err := model.PaperAll()
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", ids)
}
