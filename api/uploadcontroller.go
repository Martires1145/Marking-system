package api

import (
	"github.com/gin-gonic/gin"
	"marking/response"
	"marking/util"
)

// Upload
//
//	@Tags		Upload
//	@Summary	上传文件
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		data			formData	file	true	"文件数据"
//	@Param		type			formData	string	true	"文件类型 paper avatar answer"
//	@Param		Authorization	header		string	true	"token"
//	@Success	200				{object}	response.Response
//	@Router		/api/upload [post]
func Upload(c *gin.Context) {
	tye := c.PostForm("type")
	data, err := c.FormFile("data")
	if err != nil {
		response.Fail(c.Writer, err.Error(), 400)
		return
	}

	p, err := util.GetFileSavePath(data.Filename, tye)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}

	err = c.SaveUploadedFile(data, p)

	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}

	name := util.GetFileUrl(p)

	response.Success(c.Writer, "success", name)
}
