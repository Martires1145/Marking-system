package api

import (
	"github.com/gin-gonic/gin"
	"marking/model"
	"marking/response"
	"marking/util"
)

// GetStudentScores
//
//	@Summary	获取学生成绩
//	@Tags		student
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/stu/scores [post]
func GetStudentScores(c *gin.Context) {
	token := c.GetHeader("Authorization")
	_, claim, _ := util.ParseToken(token)

	scores, err := model.GetScore(claim.UID)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	response.Success(c.Writer, "success", scores)
}

func PersonalAnalysis(c *gin.Context) {}
