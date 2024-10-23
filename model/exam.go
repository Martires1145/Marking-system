package model

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"marking/common"
	"strconv"
)

type Exam struct {
	ID         uint   `db:"id" json:"id"`
	CreateTime uint   `db:"create_time" json:"createTime"`
	State      int    `db:"state" json:"state"`
	Name       string `db:"name" json:"name"`
	Desc       string `db:"desc" json:"desc"`
	Subjects   string `db:"subjects" json:"subjects"`
}

type ExamDataReturn struct {
	ExamID     uint   `db:"id" json:"examID"`
	PaperID    []uint `db:"paper_id" json:"paperID"`
	State      int    `db:"state" json:"state"`
	Name       string `db:"name" json:"name"`
	Desc       string `db:"desc" json:"desc"`
	Subjects   string `db:"subjects" json:"subjects"`
	CreateTime uint   `db:"create_time" json:"createTime"`
	ClassIDs   []uint `json:"classIDs"`
}

type Mission struct {
	ExamID string `db:"exam_id"`
	Start  int    `db:"start"`
	End    int    `db:"end"`
}

type ScoreMsg struct {
	StudentID uint `json:"studentID" db:"student_id"`
	ExamID    uint `json:"examID" db:"exam_id"`
	Score     int  `json:"score" db:"score"`
}

var (
	Preparation = 1
	Sharing     = 2
	Marking     = 3
	Ended       = 4
)

type CutMg struct {
	ExamID string   `json:"examID"`
	Cut    []string `json:"cut"`
}

func (e Exam) Save() error {
	db := common.DB
	sqlStr := "INSERT INTO exam(name, `desc`, create_time, subjects) VALUE(:name, :desc, :create_time, :subjects)"
	_, err := db.NamedExec(sqlStr, &e)
	return err
}

func (e Exam) Update() error {
	db := common.DB
	sqlStr := "INSERT INTO exam(name, `desc`, subjects, state) VALUE(:name, :desc, :subjects, :state)"
	_, err := db.NamedExec(sqlStr, &e)
	return err
}

func DeleteExam(id string) error {
	db := common.DB
	sqlStr := "DELETE FROM exam_class_relation WHERE exam_id = ?"
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		return err
	}

	sqlStr = "DELETE FROM exam WHERE id = ?"
	_, err = db.Exec(sqlStr, id)
	return err
}

func GetExams(ids []uint) (eInfos []ExamDataReturn, err error) {
	db := common.DB
	sqlStr := "SELECT * FROM exam WHERE id in (?)"
	q, args, err := sqlx.In(sqlStr, ids)
	err = db.Select(&eInfos, q, args...)
	if err != nil {
		return nil, err
	}

	for idx := range eInfos {
		sqlStr = "SELECT id FROM paper WHERE exam_id = ?"
		err = db.Select(&eInfos[idx].PaperID, sqlStr, eInfos[idx].ExamID)

		sqlStr = "SELECT class_id FROM exam_class_relation WHERE exam_id = ?"
		err = db.Select(&eInfos[idx].ClassIDs, sqlStr, eInfos[idx].ExamID)
		if err != nil {
			return nil, err
		}
	}
	return
}

func GetExamNum(examID string, t string) (num int, err error) {
	switch t {
	case "teacher":
		return getExamNumOfTeacher(examID)
	case "student":
		return getExamNumOfStudent(examID)
	case "task":
		return getExamNumOfTask(examID)
	}
	return
}

func getExamNumOfTeacher(examID string) (num int, err error) {
	db := common.DB
	sqlStr := "SELECT COUNT(1) FROM (SELECT DISTINCT teacher_id FROM class_teacher_relation WHERE class_id in (SELECT class_id FROM exam_class_relation WHERE exam_id = ?)) t"
	err = db.Get(&num, sqlStr, examID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return
}

func getExamNumOfStudent(examID string) (num int, err error) {
	db := common.DB
	sqlStr := "SELECT COUNT(1) FROM (SELECT DISTINCT student_id FROM class_student_relation WHERE class_id in (SELECT class_id FROM exam_class_relation WHERE exam_id = ?)) t"
	err = db.Get(&num, sqlStr, examID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return
}

func getExamNumOfTask(examID string) (num int, err error) {
	db := common.DB
	sqlStr := "SELECT  IF(COUNT(offset) = 0, -1, MAX(offset)) cnt FROM paper_part WHERE paper_id = (SELECT paper_id FROM exam WHERE exam.id = ?)"
	err = db.Get(&num, sqlStr, examID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return num + 1, err
}

func SavePaperMission(examID string, cut []string) error {
	db := common.DB
	sqlStr := "INSERT INTO paper_mission(exam_id, start, end) VALUES (:exam_id, :start, :end)"

	var cutPoint []Mission
	last := 0
	for _, point := range cut {
		offset, _ := strconv.Atoi(point)
		cutPoint = append(cutPoint, Mission{
			ExamID: examID,
			Start:  last,
			End:    last + offset - 1,
		})

		last = last + offset
	}

	_, err := db.NamedExec(sqlStr, cutPoint)

	sqlStr = "UPDATE exam SET state = 3 WHERE id = ?"
	_, err = db.Exec(sqlStr, examID)
	if err != nil {
		return err
	}
	return err
}

func GetMission(teacher uint, examID string) (m Mission, err error) {
	db := common.DB

	sqlStr := `SELECT exam_id, start, end
     		   FROM (
    				(SELECT
    				 exam_id, start, end, 1 AS r
    				 FROM paper_mission
					 WHERE teacher_id = ?
					 AND exam_id= ?)
					 UNION ALL
					 (SELECT
					 exam_id, start, end, 0 AS r
					 FROM paper_mission
					 WHERE teacher_id <> ?
					 AND exam_id = ?
					 AND teacher_id = -1
					 ORDER BY RAND()
					 LIMIT 1)
					 ) s
			   ORDER BY s.r DESC
			   LIMIT 1;`
	err = db.Get(&m, sqlStr, teacher, examID, teacher, examID)
	if err != nil {
		return
	}

	sqlStr = "UPDATE paper_mission SET state = 1, teacher_id = ? WHERE exam_id = ? AND start = ?"

	_, err = db.Exec(sqlStr, teacher, m.ExamID, m.Start)
	return
}

func GetAllExam() (ids []int, err error) {
	db := common.DB
	sqlStr := "SELECT id FROM exam"
	err = db.Select(&ids, sqlStr)
	return
}

func GetAllScore(id string) (scores []ScoreMsg, err error) {
	db := common.DB
	sqlStr := "SELECT e.id AS exam_id, u.uid AS student_id, p.mark AS score FROM users u, paper p, exam e WHERE p.owner_id = u.uid AND p.exam_id = e.id AND e.id = ?"
	err = db.Select(&scores, sqlStr, id)
	return
}

func GetExamIn(uid uint) (exams []Exam, err error) {
	db := common.DB
	sqlStr := `SELECT 
    		exam.id, exam.create_time, exam.state, exam.name, exam.desc, exam.subjects 
			FROM exam, class_teacher_relation c, exam_class_relation e 
			WHERE c.teacher_id = ? 
			AND  e.class_id = c.class_id
			AND e.exam_id = exam.id`

	err = db.Select(&exams, sqlStr, uid)
	return
}
