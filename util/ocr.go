package util

import (
	"context"
	"encoding/json"
	"fmt"
	"marking/common"
)

type OcrRequest struct {
	ID  int    `json:"ID" db:"id"`
	Img string `json:"img" db:"img"`
}

type TextData struct {
	Transcription string  `json:"transcription"`
	Points        [][]int `json:"points"`
}

type OcrResponse struct {
	ID   int    `json:"ID"`
	Text string `json:"text"`
}

func Ocr(ocrData []OcrRequest) (text []OcrResponse, err error) {
	topicR, topicW := GetTopic()
	common.WriteTopicID(context.Background(), topicR, topicW)

	//common.CreateTopic(context.Background(), topicR, topicW)
	//
	//ack := common.ReadTopicAck(context.Background(), topicW)
	//if !ack {
	//	fmt.Println("time out")
	//	return nil, err
	//}
	//
	//common.WriteTopicAck(context.Background(), topicW)

	err = writeMsg(context.Background(), topicW, ocrData)
	if err != nil {
		return nil, err
	}

	return readMsg(context.Background(), topicR, len(ocrData))
}

func writeMsg(ctx context.Context, topic string, ocrRequest []OcrRequest) error {
	jsonData, _ := json.Marshal(ocrRequest)
	err := common.WriteMsg(ctx, topic, jsonData)
	return err
}

func readMsg(ctx context.Context, topic string, dataLen int) (resps []OcrResponse, err error) {
	resps = make([]OcrResponse, 0)
	dataChan := make(chan []byte)
	go common.ReadMsg(ctx, topic, &dataChan)
	fmt.Println("len:", dataLen)
	for i := 0; i < dataLen; i++ {
		data := <-dataChan
		fmt.Println("data:", i)
		var resp OcrResponse
		err := json.Unmarshal(data, &resp)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		resps = append(resps, resp)
	}
	return
}
