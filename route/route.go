package route

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"marking/api"
	"marking/docs"
)

func GetRoute() *gin.Engine {
	r := gin.Default()

	r.Use(api.Cors)

	r.Static("/static", "./static")

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.POST("/api/upload", api.VerifyToken, api.Upload)

	user := r.Group("/api/user")
	{
		user.POST("/new", api.Register)
		user.POST("/check", api.CheckUserName)
		user.POST("/verify", api.Verify)
		user.POST("/login", api.Login)
		user.POST("/info", api.UserInfo)
		user.POST("/ru", api.UserRevise)
		user.POST("/c", api.UserClass)
		user.POST("/all", api.UserAll)
		user.POST("/list", api.UserInfoList)
	}

	class := r.Group("/api/c")
	{
		class.POST("/new", api.VerifyToken, api.ClassCreate)
		class.POST("/rc", api.VerifyToken, api.ClassRevise)
		class.POST("/info", api.VerifyToken, api.ClassInfo)
		class.POST("/all", api.VerifyToken, api.AllClass)
	}

	classMg := r.Group("/api/c/mg")
	{
		classMg.POST("/quit", api.VerifyToken, api.QuitClass)
		classMg.POST("/clear", api.VerifyToken, api.ClearStudent)
		classMg.POST("/add", api.VerifyToken, api.AddClass)
		classMg.POST("/score", api.VerifyToken, api.ClassScore)
		classMg.POST("/sh", api.VerifyToken, api.ClassHistoryScore)
	}

	exam := r.Group("/api/exam")
	{
		exam.POST("/new", api.VerifyToken, api.CreateExam)
		exam.POST("/delete", api.VerifyToken, api.DeleteExam)
		exam.POST("/re", api.VerifyToken, api.ReviseExam)
		exam.POST("/info", api.VerifyToken, api.ExamInfo)
		exam.POST("/all", api.VerifyToken, api.AllExam)
	}

	examMg := r.Group("/api/exam/mg")
	{
		examMg.POST("/join", api.VerifyToken, api.ClassJoinExam)
		examMg.POST("/quit", api.VerifyToken, api.ClassQuitExam)
		examMg.POST("/plan", api.VerifyToken, api.DistributePaper)
		examMg.POST("/num", api.VerifyToken, api.GetExamNum)
		examMg.POST("/score", api.VerifyToken, api.GetScore)
	}

	paper := r.Group("/api/paper")
	{
		paper.POST("/new", api.VerifyToken, api.NewPaper)
		paper.POST("/na", api.VerifyToken, api.NewAnswer)
		paper.POST("/sum", api.VerifyToken, api.TotalScore)
		paper.POST("/delete", api.VerifyToken, api.DeletePaper)
		paper.POST("/info", api.VerifyToken, api.PaperInfos)
		paper.POST("/all", api.VerifyToken, api.AllPaper)
	}

	student := r.Group("/api/stu")
	{
		student.POST("/scores", api.VerifyToken, api.GetStudentScores)
		student.POST("/anal", api.VerifyToken, api.PersonalAnalysis)
	}

	teacher := r.Group("/api/seisei")
	{
		teacher.POST("/mission", api.VerifyToken, api.GetTaskMission)
		teacher.POST("/mark", api.VerifyToken, api.MarkPaperPart)
		teacher.POST("/c/anal", api.VerifyToken, api.GetClassAnalysis)
		teacher.POST("/as", api.VerifyToken, api.GetAnswerOfPaper)
		teacher.POST("/ms/part", api.VerifyToken, api.GetMissionPart)
		teacher.POST("/ms/fm", api.VerifyToken, api.FinishMission)
		teacher.POST("/ms/exam", api.VerifyToken, api.GetExamIn)
	}

	return r
}
