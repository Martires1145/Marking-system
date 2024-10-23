package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"marking/model"
	"marking/response"
	"marking/util"
)

// Register
//
//	@Summary	用户注册
//	@Tags		User
//	@Param		ud	body	model.UserJson	true	"用户信息"
//	@Router		/api/user/new [post]
func Register(c *gin.Context) {
	var userJ model.UserJson
	err := c.ShouldBindJSON(&userJ)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}
	user := userJ.ToUser()

	err = user.Check(userJ.Captcha)
	if err != nil {
		response.Fail(c.Writer, "wrong captcha", 403)
		return
	}

	id, err := user.Save()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("error when save the user data, err:%s", err.Error()), 500)
		return
	}
	user.Uid = uint(id)

	err = user.SetRole()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("role error, err:%s", err.Error()), 403)
		return
	}

	response.Success(c.Writer, "success", nil)
}

// Verify
//
//	@Summary	用户验证
//	@Tags		User
//	@Param		userName	formData	string	true	"用户名"
//	@Param		email		formData	string	true	"邮箱"
//	@Router		/api/user/verify [post]
func Verify(c *gin.Context) {
	var user model.User
	user.UserName = c.PostForm("userName")
	user.Email = c.PostForm("email")
	captcha := util.RandToken()

	err := util.SendMessage(captcha, user.Email)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("error when send email, err:%s", err.Error()), 500)
		return
	}
	err = user.SaveCaptcha(captcha)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("error when save the captcha, err:%s", err.Error()), 500)
		return
	} else {
		response.Fail(c.Writer, "success", 200)
	}
}

// Login
//
//	@Summary	登录
//	@Tags		User
//	@Param		userName	formData	string	true	"用户名"
//	@Param		passWord	formData	string	true	"密码"
//	@Router		/api/user/login [post]
func Login(c *gin.Context) {
	var user model.User
	user.UserName = c.PostForm("userName")
	user.PassWord = c.PostForm("passWord")

	err, u := user.Login()
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = util.NoSuchUserError
		response.Fail(c.Writer, err.Error(), 403)
		return
	}
	if err != nil {
		response.Fail(c.Writer, err.Error(), 403)
		return
	}
	token, err := util.CreatToken(u.Uid, u.Role)
	if err != nil {
		response.Fail(c.Writer, err.Error(), 500)
		return
	}
	data := gin.H{
		"token":     token,
		"user_info": u,
	}
	response.Success(c.Writer, "success", data)
}

// CheckUserName
//
//	@Summary	检验用户名
//	@Tags		User
//	@Param		userName	formData	string	true	"用户名"
//	@Router		/api/user/check [post]
func CheckUserName(c *gin.Context) {
	var user model.User
	user.UserName = c.PostForm("userName")
	if !user.CheckUserName() {
		response.Fail(c.Writer, "user already exist", 400)
		return
	}
	response.Success(c.Writer, "ok", nil)
}

// UserInfo
//
//	@Summary	用户信息
//	@Tags		User
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/user/info [post]
func UserInfo(c *gin.Context) {
	var u model.User
	token := c.GetHeader("Authorization")
	t, claim, err := util.ParseToken(token)
	if err != nil {
		response.Fail(c.Writer, "failed to parse token", 400)
		return
	}
	if !t.Valid {
		response.Fail(c.Writer, "failed", 400)
		return
	}
	u.Uid = claim.UID
	nu, err := u.Info()
	if err != nil {
		response.Fail(c.Writer, "failed", 500)
	}
	response.Success(c.Writer, "success", nu)
}

// UserRevise
//
//	@Summary	用户修改信息
//	@Tags		User
//	@Param		ud	body	model.UserJson	true	"用户信息"
//	@Router		/api/user/rv [post]
func UserRevise(c *gin.Context) {
	var userJ model.UserJson
	err := c.ShouldBindJSON(&userJ)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	err = userJ.Save()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("error when save the user data, err:%s", err.Error()), 500)
		return
	}

	response.Success(c.Writer, "success", nil)
}

// UserClass
//
//	@Summary	获取用户所在的班级
//	@Tags		User
//	@Param		Authorization	header	string	true	"token"
//	@Router		/api/user/c [post]
func UserClass(c *gin.Context) {
	var u model.User
	token := c.GetHeader("Authorization")
	_, claim, _ := util.ParseToken(token)
	u.Uid = claim.UID
	u.Role = claim.Role

	class, err := u.GetClass()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("err:%s", err.Error()), 500)
		return
	}

	response.Success(c.Writer, "success", class)
}

// UserInfoList
//
//	@Summary	用户信息列表
//	@Tags		User
//	@Param		IDs	body	[]string	true	"用户ID列表"
//	@Router		/api/user/list [post]
func UserInfoList(c *gin.Context) {
	var u model.User
	var list []string
	err := c.ShouldBind(&list)
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	nu, err := u.InfoList(list)
	if err != nil {
		response.Fail(c.Writer, "failed", 500)
	}
	response.Success(c.Writer, "success", nu)
}

// UserAll
//
//	@Summary	获取用户ID列表
//	@Tags		User
//	@Router		/api/user/all [post]
func UserAll(c *gin.Context) {
	var u model.User
	ids, err := u.GetAll()
	if err != nil {
		response.Fail(c.Writer, fmt.Sprintf("data bind error, err:%s", err.Error()), 400)
		return
	}

	response.Success(c.Writer, "success", ids)
}
