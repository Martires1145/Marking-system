package util

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var path = map[string]string{
	"paper":  "./static/paper/",
	"avatar": "./static/avatar/",
	"answer": "./static/answer/",
}

func GetFileSavePath(name string, typ string) (p string, err error) {
	if typ != "paper" && typ != "avatar" && typ != "answer" {
		return "", WrongTypeError
	}
	p = path[typ]

	p = p + fmt.Sprintf("%d", time.Now().Unix()%1000000) + name

	return p, nil
}

func GetFileUrl(path string) string {
	p := viper.GetString("web.protocol")
	host := viper.GetString("web.host")
	return fmt.Sprintf("%s://%s%s", p, host, path[1:])
}
