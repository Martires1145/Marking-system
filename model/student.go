package model

import "marking/common"

func StudentQuitClass(uid uint, classID string) error {
	db := common.DB
	sqlStr := "DELETE FROM class_student_relation WHERE student_id = ? AND class_id = ?"
	_, err := db.Exec(sqlStr, uid, classID)
	return err
}

func StudentAddClass(uid uint, classCode string) error {
	db := common.DB
	sqlStr := "INSERT INTO class_student_relation(STUDENT_ID, CLASS_ID)  VALUE(?, (SELECT id FROM class WHERE code = ?))"
	_, err := db.Exec(sqlStr, uid, classCode)
	return err
}

func GetScore(uid uint) (scores []ScoreMsg, err error) {
	db := common.DB
	sqlStr := "SELECT mark as score, exam_id, owner_id AS student_id FROM paper WHERE owner_id = ? AND state = 2"
	err = db.Select(&scores, sqlStr, uid)
	return
}
