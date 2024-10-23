package model

import "marking/common"

func TeacherQuitClass(uid uint, classID string) error {
	db := common.DB
	sqlStr := "SELECT count(1) FROM class_teacher_relation WHERE class_id = ?"
	cnt := 0
	_ = db.Select(&cnt, sqlStr, classID)

	sqlStr = "DELETE FROM class_teacher_relation WHERE teacher_id = ? AND class_id = ?"
	_, err := db.Exec(sqlStr, uid, classID)

	if cnt == 1 {
		sqlStr1 := "DELETE FROM class_student_relation WHERE class_id = ?"
		sqlStr2 := "DELETE FROM class WHERE id = ?"

		db.Exec(sqlStr1, classID)
		db.Exec(sqlStr2, classID)
	}

	return err
}

func TeacherAddClass(uid uint, classCode string) error {
	db := common.DB
	sqlStr := "INSERT INTO class_teacher_relation(teacher_id, class_id)  VALUE(?, (SELECT id FROM class WHERE code = ?))"
	_, err := db.Exec(sqlStr, uid, classCode)
	return err
}
