package model

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"marking/common"
	"marking/util"
	"time"
)

type Paper struct {
	ID      uint `db:"id"`
	OwnerID uint `db:"owner_id"`
	ExamID  uint `db:"exam_id"`
	State   int  `db:"state"`
	Mark    int  `db:"mark"`
}

type PaperPart struct {
	OffSet   int    `db:"offset"`
	Mark     int    `db:"mark"`
	MaxMark  int    `db:"max_mark"`
	State    int    `db:"state"`
	ID       uint   `db:"id"`
	ExamID   uint   `db:"exam_id"`
	PaperID  uint   `db:"paper_id"`
	Img      string `db:"img"`
	Text     string `db:"text"`
	Comments string `db:"comments"`
}

type PaperJson struct {
	OwnerID uint     `json:"ownerID"`
	ExamID  uint     `json:"examID"`
	Imgs    []string `json:"imgs"`
	MaxMark []int    `json:"maxMark"`
}

type Answer struct {
	ID     uint   `json:"ID" db:"id"`
	ExamID uint   `json:"examID" db:"exam_id"`
	Offset int    `json:"offset" db:"offset"`
	Img    string `json:"img" db:"img"`
}

type AnswerJson struct {
	ExamID uint     `json:"examID" db:"exam_id"`
	Imgs   []string `json:"img" db:"imgs"`
}

type PaperScore struct {
	PaperID uint `json:"paperID" db:"paper_id"`
	Mark    int  `json:"mark" db:"mark"`
}

var (
	Unmodified = 0
	Revised    = 1
	Judge      = 2
	Done       = 3
)

func GetAnswer(examID string) (a []Answer, err error) {
	db := common.DB
	sqlStr := "SELECT * FROM paper_answer WHERE exam_id = ?"
	err = db.Select(&a, sqlStr, examID)
	return
}

func GetPaperPart(examID string, uid uint) (p []PaperPart, err error) {
	db := common.DB
	var mission Mission
	sqlStr := "SELECT start, end FROM paper_mission WHERE exam_id = ? AND teacher_id = ? AND state != ?"
	err = db.Get(&mission, sqlStr, examID, uid, Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	sqlStr = "SELECT * FROM paper_part WHERE paper_id in (SELECT id FROM paper WHERE exam_id = ?) AND offset BETWEEN ? AND ?"
	err = db.Select(&p, sqlStr, examID, mission.Start, mission.End)
	return
}

func MarkPaperPart(partID string, mark int) error {
	db := common.DB
	sqlStr := "UPDATE paper_part SET mark = ?, state = 1 WHERE id = ?"
	_, err := db.Exec(sqlStr, mark, partID)
	return err
}

func SavePapers(papers []PaperJson) error {
	db := common.DB
	sqlStr := "INSERT INTO paper(owner_id, exam_id, state) VALUE (?, ?, 1)"

	for _, paper := range papers {
		exec, err := db.Exec(sqlStr, paper.OwnerID, paper.ExamID)
		if err != nil {
			return err
		}

		id, _ := exec.LastInsertId()
		sqlStr1 := "INSERT INTO paper_part(paper_id, offset, img, max_mark, comments, text) VALUES (:paper_id, :offset, :img, :max_mark, :comments, :text)"
		var paperPart []PaperPart
		for i := range paper.Imgs {
			paperPart = append(paperPart, PaperPart{
				OffSet:   i,
				MaxMark:  paper.MaxMark[i],
				PaperID:  uint(id),
				Img:      paper.Imgs[i],
				Text:     "",
				Comments: "",
			})
		}
		_, err = db.NamedExec(sqlStr1, paperPart)
		if err != nil {
			return err
		}

		sqlStr = "UPDATE exam SET state = 2 WHERE id = ?"
		_, err = db.Exec(sqlStr, paper.ExamID)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateTotalScore(paperIDs []int) error {
	db := common.DB
	sqlStr := "UPDATE paper SET mark = (SELECT SUM(mark) FROM paper_part WHERE paper_id = paper.id) WHERE id IN (?)"
	query, args, err := sqlx.In(sqlStr, paperIDs)
	if err != nil {
		return err
	}
	_, err = db.Exec(query, args...)
	return err
}

func DeletePaper(paperID string) error {
	db := common.DB
	sqlStr := "DELETE FROM paper WHERE id = ?"
	_, err := db.Exec(sqlStr, paperID)
	if err != nil {
		return err
	}
	sqlStr = "DELETE FROM paper_part WHERE paper_id = ?"
	_, err = db.Exec(sqlStr, paperID)
	if err != nil {
		return err
	}
	return nil
}

func PaperInfo(paperIDs []int) (papers []Paper, err error) {
	db := common.DB
	sqlStr := "SELECT * FROM paper WHERE id IN (?)"
	query, args, err := sqlx.In(sqlStr, paperIDs)
	if err != nil {
		return nil, err
	}

	err = db.Select(&papers, query, args...)
	return
}

func SaveAnswer(answerJson *AnswerJson) error {
	db := common.DB

	var answer []Answer
	for i := range answerJson.Imgs {
		answer = append(answer, Answer{
			ExamID: answerJson.ExamID,
			Offset: i,
			Img:    answerJson.Imgs[i],
		})
	}

	sqlStr := "INSERT INTO paper_answer(exam_id, offset, img) VALUES (:exam_id, :offset, :img)"
	_, err := db.NamedExec(sqlStr, answer)
	return err
}

func FinishMission(examID string, teacherID uint) error {
	db := common.DB
	sqlStr := "UPDATE paper_mission SET state = ? WHERE exam_id = ? AND teacher_id = ?"
	_, err := db.Exec(sqlStr, Done, examID, teacherID)
	if err != nil {
		return err
	}

	sqlStr = "SELECT COUNT(1) FROM paper_mission WHERE state != ?"
	cnt := 0
	err = db.Get(&cnt, sqlStr, Done)
	if err != nil {
		return err
	}

	if cnt == 0 {
		sqlStr = "UPDATE exam SET state = ? WHERE id = ?"
		_, err = db.Exec(sqlStr, Ended, examID)
	}

	return err
}

func PaperAll() (ids []int, err error) {
	db := common.DB
	sqlStr := "SELECT id FROM paper"
	err = db.Select(&ids, sqlStr)
	return
}

type IDAndMaxScore struct {
	ID      int `json:"ID" db:"id"`
	MaxMark int `json:"maxMark" db:"max_mark"`
}

func GetOrcRequest() (orc []util.OcrRequest, err error) {
	db := common.DB
	sqlStr := "SELECT id, img FROM paper_part WHERE state = ?"
	err = db.Select(&orc, sqlStr, Unmodified)
	return
}

func SaveOrcResponse(ocrResponse []util.OcrResponse) error {
	db := common.DB
	sqlStr := "UPDATE paper_part SET text = ?, state = ? WHERE id = ?"
	for _, rsp := range ocrResponse {
		_, err := db.Exec(sqlStr, rsp.Text, Revised, rsp.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetJudgeRequest() (judge []util.ModelJudgeRequest, err error) {
	db := common.DB
	sqlStr := "SELECT id, text, max_mark FROM paper_part WHERE state = ?"
	err = db.Select(&judge, sqlStr, Revised)
	return
}

func SaveJudgeResponse(judgeResponse []util.ModelJudgeResponse) error {
	db := common.DB
	sqlStr := "UPDATE paper_part SET comments = ?, mark = ?, state = ? WHERE id = ?"
	for _, rsp := range judgeResponse {
		_, err := db.Exec(sqlStr, rsp.Comments, rsp.Score, Judge, rsp.PartID)
		if err != nil {
			return err
		}
	}
	return nil
}

func OcrAndJudge() {
	for {
		ocrData, err := GetOrcRequest()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if len(ocrData) == 0 {
			continue
		}

		text, err := util.Ocr(ocrData)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		err = SaveOrcResponse(text)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		judge, err := GetJudgeRequest()

		modelResp, err := util.ModelJudge(judge)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		err = SaveJudgeResponse(modelResp)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		time.Sleep(time.Second * 60 * 2)
	}
}
