package model

import (
	"github.com/jmoiron/sqlx"
	"marking/common"
)

type Class struct {
	ID   uint   `db:"id"`
	Name string `db:"name"`
}

type ClassJson struct {
	Name      string `json:"name"`
	TeacherID uint   `json:"teacherID"`
}

type ClassDataReturn struct {
	ID         uint   `db:"id"`
	Name       string `db:"name"`
	Code       string `db:"code"`
	TeacherIDs []uint
	StudentIDs []uint
}

func (j ClassJson) Save() (token string, err error) {
	var t []string
	db := common.DB
	sqlStr := "INSERT INTO class(name) VALUE(?)"
	exec, err := db.Exec(sqlStr, j.Name)
	if err != nil {
		return "", nil
	}
	id, _ := exec.LastInsertId()
	if err != nil {
		return "", nil
	}

	sqlStr = "INSERT INTO class_teacher_relation(teacher_id, class_id) VALUE(?, ?)"
	_, err = db.Exec(sqlStr, j.TeacherID, id)
	if err != nil {
		return "", nil
	}

	sqlStr = "SELECT rand_value FROM dt ORDER BY rand_num LIMIT ?, 1"
	err = db.Select(&t, sqlStr, id)
	if err != nil {
		return "", nil
	}

	token = t[0]
	sqlStr = "UPDATE class SET code = ? WHERE id = ?"
	_, err = db.Exec(sqlStr, token, id)

	return
}

func (j ClassJson) Update() error {
	db := common.DB
	sqlStr := "INSERT INTO class(name) VALUE(?)"
	_, err := db.Exec(sqlStr, j.Name)
	return err
}

func GetClass(ids []uint) (cInfos map[uint]ClassDataReturn, err error) {
	cInfos = map[uint]ClassDataReturn{}
	var c []ClassDataReturn
	db := common.DB
	sqlStr := "SELECT id, name, code FROM class WHERE id in (?)"
	q, args, err := sqlx.In(sqlStr, ids)
	err = db.Select(&c, q, args...)
	if err != nil {
		return nil, err
	}

	for _, class := range c {
		sqlStr = "SELECT student_id FROM class_student_relation WHERE class_id = ?"
		err = db.Select(&class.StudentIDs, sqlStr, class.ID)
		if err != nil {
			return nil, err
		}

		sqlStr = "SELECT teacher_id FROM class_teacher_relation WHERE class_id = ?"
		err = db.Select(&class.TeacherIDs, sqlStr, class.ID)
		if err != nil {
			return nil, err
		}

		cInfos[class.ID] = class
	}
	return
}

func QuitExam(classID, examID string) error {
	db := common.DB
	sqlStr := "DELETE FROM exam_class_relation WHERE class_id = ? AND  exam_id = ?"
	_, err := db.Exec(sqlStr, classID, examID)
	return err
}

func JoinExam(classID, examID string) error {
	db := common.DB
	sqlStr := "INSERT INTO exam_class_relation(CLASS_ID, EXAM_ID) VALUE(?, ?)"
	_, err := db.Exec(sqlStr, classID, examID)
	return err
}

func GetAllClass() (ids []int, err error) {
	db := common.DB
	sqlStr := "SELECT id FROM class"
	err = db.Select(&ids, sqlStr)
	return
}

func GetClassScore(classID, examID string) (scores []ScoreMsg, err error) {
	db := common.DB
	sqlStr := "SELECT e.id AS `exam_id`, u.uid AS `student_id`, p.mark AS `score` FROM users u, paper p, exam e, exam_class_relation c WHERE p.owner_id = u.uid AND p.exam_id = e.id AND c.class_id = ? AND c.exam_id = e.id AND e.id = ?"
	err = db.Select(&scores, sqlStr, classID, examID)
	return
}

func GetClassScoreHistory(classID, n string) (scores []ScoreMsg, err error) {
	db := common.DB
	sqlStr := "SELECT e.id AS `exam_id`, u.uid AS `student_id`, p.mark AS `score` FROM users u, paper p, exam e, exam_class_relation c WHERE p.owner_id = u.uid AND p.exam_id = e.id AND c.class_id = ? AND c.exam_id = e.id ORDER BY e.create_time DESC LIMIT ?"
	err = db.Select(&scores, sqlStr, classID, n)
	return
}
