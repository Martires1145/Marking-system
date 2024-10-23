package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"marking/model"
	"marking/response"
)

// ClassCreate
//
//	@Summary	创建班级
//	@Tags		Class
//	@Param		classMg			body	model.ClassJson	true	"班级信息"
//	@Param		Authorization	header	string			true	"token"
//	@Router		/api/c/new [post]
func ClassCreate(c *gin.Context) {
	var classJson model.ClassJson
	err := c.ShouldBindJSON(&classJson)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	token, err := classJson.Save()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("error when save the class data, err:%s", err.Error()), 500)
		return
	}

	response.Success(c.Writer, "success", gin.H{"token": token})
}

// ClassRevise
//
//	@Summary	修改班级信息
//	@Tags		Class
//	@Param		classMg			body	model.ClassJson	true	"班级信息"
//	@Param		Authorization	header	string			true	"token"
//	@Router		/api/c/rc [post]
func ClassRevise(c *gin.Context) {
	var classJson model.ClassJson
	err := c.ShouldBindJSON(&classJson)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	err = classJson.Update()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("error when update the class data, err:%s", err.Error()), 500)
		return
	}
	response.Success(c.Writer, "success", nil)
}

// ClassInfo
//
//	@Summary	获取班级信息
//	@Tags		Class
//	@Param		classIDs		body	[]int	true	"班级ID列表"
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/c/info [post]
func ClassInfo(c *gin.Context) {
	var classIDs []uint
	err := c.ShouldBind(&classIDs)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	cInfos, err := model.GetClass(classIDs)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 400)
		return
	}

	response.Success(c.Writer, "success", cInfos)
}

// AllClass
//
//	@Summary	获取所有班级ID
//	@Tags		Class
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/c/all [post]
func AllClass(c *gin.Context) {
	ids, err := model.GetAllClass()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 400)
		return
	}

	response.Success(c.Writer, "success", ids)
}
