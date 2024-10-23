package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/spf13/viper"
	"strings"
)

var (
	accessKey = viper.GetString("model.AccessKey")
	secretKey = viper.GetString("model.SecretKey")
	chat      *qianfan.ChatCompletion
)

type ModelJudgeResponse struct {
	PartID   int    `json:"partID"`
	Score    string `json:"score"`
	Comments string `json:"comments"`
}

type ModelJudgeRequest struct {
	PartID   int    `json:"partID" db:"id"`
	Text     string `json:"text" db:"text"`
	MaxScore int    `json:"maxScore" db:"max_mark"`
}

type ScoresResp struct {
	PartID int    `json:"partID"`
	Score  string `json:"score"`
}

type CommentsResp struct {
	PartID   int    `json:"partID"`
	Comments string `json:"comments"`
}

func init() {
	qianfan.GetConfig().AccessKey = accessKey
	qianfan.GetConfig().SecretKey = secretKey
	chat = qianfan.NewChatCompletion(
		qianfan.WithModel("ERNIE-Bot"),
	)
}

func ModelJudge(rqs []ModelJudgeRequest) (scores []ModelJudgeResponse, err error) {
	scores = make([]ModelJudgeResponse, 0)

	for _, rq := range rqs {
		msg := fmt.Sprintf("下面是题目以及学生解答:\n%s\n该题满分为%d分\n请严格按以下形式给出结果:\n评分:xxx\n评语:xxx", rq.Text, rq.MaxScore)
		fmt.Println(msg)
		resp, err := chat.Do(
			context.Background(),
			&qianfan.ChatCompletionRequest{
				Messages: []qianfan.ChatCompletionMessage{
					qianfan.ChatCompletionUserMessage(msg),
				},
			},
		)

		if err != nil {
			return nil, err
		}
		result, err := processTheResults(resp.Result, rq.PartID)
		if err != nil {
			return nil, err
		}
		scores = append(scores, result)
	}

	return
}

func processTheResults(text string, id int) (resp ModelJudgeResponse, err error) {
	splitResult := strings.Split(text, "评语:")
	if len(splitResult) != 2 {
		return ModelJudgeResponse{}, errors.New("wrong text")
	}
	return ModelJudgeResponse{
		PartID:   id,
		Score:    strings.Split(splitResult[0], ":")[1],
		Comments: splitResult[1],
	}, nil
}

func Scores(m []ModelJudgeResponse) []ScoresResp {
	sr := make([]ScoresResp, 0)
	for _, re := range m {
		sr = append(sr, ScoresResp{
			PartID: re.PartID,
			Score:  re.Score,
		})
	}
	return sr
}

func Comments(m []ModelJudgeResponse) []CommentsResp {
	cr := make([]CommentsResp, 0)
	for _, re := range m {
		cr = append(cr, CommentsResp{
			PartID:   re.PartID,
			Comments: re.Comments,
		})
	}
	return cr
}
