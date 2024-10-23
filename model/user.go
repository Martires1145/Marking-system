package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"marking/common"
	"marking/util"
	"time"
)

type User struct {
	Uid      uint   `db:"uid" json:"uid,string"`
	Role     int    `db:"role" json:"role"`
	Name     string `db:"name" json:"name"`
	UserName string `db:"user_name" json:"userName"`
	PassWord string `db:"pass_word" json:"passWord"`
	Avatar   string `db:"avatar" json:"avatar"`
	Email    string `db:"email" json:"email"`
}

type UserJson struct {
	Uid      uint   `db:"uid" json:"uid,string"`
	Role     int    `db:"role" json:"role"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	PassWord string `json:"passWord"`
	Avatar   string `json:"avatar"`
	Email    string `db:"email" json:"email"`
	Captcha  string `db:"captcha" json:"captcha"`
}

type Captcha struct {
	Id       uint   `db:"id"`
	Begin    uint   `db:"begin"`
	UserName string `db:"username"`
	Token    string `db:"token"`
}

func (j UserJson) ToUser() User {
	return User{
		Name:     j.Name,
		Avatar:   j.Avatar,
		UserName: j.UserName,
		PassWord: j.PassWord,
		Email:    j.Email,
		Role:     j.Role,
	}
}

type UserInfo struct {
	Uid    uint   `db:"uid" json:"uid,string"`
	Role   int    `db:"role" json:"role"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
	Email  string `db:"email" json:"email"`
}

var (
	Student = 1
	Teacher = 2
	Root    = 3
)

func (u User) GetAll() (ids []string, err error) {
	db := common.DB
	sqlStr := "SELECT uid FROM users"
	err = db.Select(&ids, sqlStr)
	return
}

func (u User) SetRole() error {
	if u.Role > Teacher {
		return util.InsufficientPermissionsError
	}

	var role string
	switch u.Role {
	case Student:
		role = "student"
	case Teacher:
		role = "teacher"
	}
	err := util.AddGroupM(fmt.Sprintf("%d", u.Uid), role)
	return err
}

func (u User) CheckUserName() bool {
	db := common.DB
	sqlStr := "SELECT COUNT(1) FROM users WHERE user_name = ?"
	cnt := 0
	_ = db.Get(&cnt, sqlStr, u.UserName)
	return cnt == 0
}

func (u User) Save() (int, error) {
	db := common.DB
	sqlStr := "INSERT INTO users(name, user_name, pass_word, avatar, email, role) values (:name, :user_name, :pass_word, :avatar, :email, :role)"
	exec, err := db.NamedExec(sqlStr, u)
	if err != nil {
		return -1, err
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return -1, err
	}
	if affected != 1 {
		return -1, util.InsertFailError
	}
	id, _ := exec.LastInsertId()
	return int(id), nil
}

func (j UserJson) Save() error {
	db := common.DB
	sqlStr := "INSERT INTO users(uid, user_name, name, pass_word, avatar) values (:uid, :user_name, :name, :pass_word, :avatar)"
	exec, err := db.NamedExec(sqlStr, j)
	if err != nil {
		return err
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return util.InsertFailError
	}
	return nil
}

func (u User) Login() (err error, user UserInfo) {
	db := common.DB
	userDB := User{}
	sqlStr := "SELECT * FROM users WHERE user_name = ?"
	err = db.Get(&userDB, sqlStr, u.UserName)
	if err != nil {
		return err, UserInfo{}
	}
	if u.PassWord == userDB.PassWord {
		return nil, UserInfo{
			Uid:    userDB.Uid,
			Role:   userDB.Role,
			Name:   userDB.Name,
			Avatar: userDB.Avatar,
			Email:  userDB.Email,
		}
	} else {
		return util.PassWordWrongError, UserInfo{}
	}
}

func (u User) SaveCaptcha(token string) error {
	db := common.DB
	sqlStr := "INSERT INTO captcha(username, token, begin) VALUE(?, ?, ?)"
	_, err := db.Exec(sqlStr, u.UserName, token, time.Now().Unix())
	return err
}

func (u User) Check(token string) error {
	var c Captcha
	db := common.DB
	sqlStr := "SELECT * FROM captcha WHERE username = ? AND begin = (SELECT max(begin) FROM captcha WHERE username = ?)"
	err := db.Get(&c, sqlStr, u.UserName, u.UserName)
	if err != nil {
		return err
	}
	if int64(c.Begin+5*60) < time.Now().Unix() {
		return util.CaptchaExpiredError
	}
	if token != c.Token {
		return util.CaptchaWrongError
	}
	return nil
}

func (u User) Info() (nu UserInfo, err error) {
	db := common.DB
	sqlStr := "SELECT uid, name, role, avatar, email FROM users WHERE uid = ?"
	err = db.Get(&nu, sqlStr, u.Uid)
	return
}

func (u User) GetClass() (id []uint, err error) {
	db := common.DB
	var sqlStr string
	if u.Role == 1 {
		sqlStr = "SELECT class_id FROM class_student_relation WHERE student_id = ?"
		err = db.Select(&id, sqlStr, u.Uid)
	} else if u.Role == 2 {
		sqlStr = "SELECT class_id FROM class_teacher_relation WHERE teacher_id = ?"
		err = db.Select(&id, sqlStr, u.Uid)
	} else {
		sqlStr = "SELECT id FROM class"
		err = db.Select(&id, sqlStr)
	}
	return
}

func (u User) QuitClass(classID string) error {
	if u.Role == 1 {
		return StudentQuitClass(u.Uid, classID)
	} else {
		return TeacherQuitClass(u.Uid, classID)
	}
}

func (u User) AddClass(classToken string) error {
	if u.Role == 1 {
		return StudentAddClass(u.Uid, classToken)
	} else {
		return TeacherAddClass(u.Uid, classToken)
	}
}

func (u User) InfoList(IDs []string) (nu []UserInfo, err error) {
	db := common.DB
	sqlStr := "SELECT uid, name, role, avatar, email FROM users WHERE uid in (?)"
	query, args, err := sqlx.In(sqlStr, IDs)
	if err != nil {
		return nil, err
	}

	err = db.Select(&nu, query, args...)
	return
}
