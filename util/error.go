package util

import "errors"

var (
	NoSuchUserError              = errors.New("no such user")
	InsertFailError              = errors.New("fail to insert")
	PassWordWrongError           = errors.New("password wrong")
	CaptchaExpiredError          = errors.New("all captcha expired")
	CaptchaWrongError            = errors.New("captcha wrong")
	InsufficientPermissionsError = errors.New("insufficient permissions error")
	WrongTypeError               = errors.New("no such type")
	NoMissionError               = errors.New("no mission")
)
