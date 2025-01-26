package routes

import (
	"database/sql"

	"github.com/cuddest/dz-skills/controllers"
	"github.com/cuddest/dz-skills/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(router *gin.Engine, db *sql.DB) {
	// swagger docs route
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//base routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Dz Skills API, go to https://dzskiils-production.up.railway.app/docs/index.html#/ for documentation, good to see you :D",
		})
	})
	// Answer Routes
	answerController := controllers.NewAnswerController(db)
	answerGroup := router.Group("/answers")
	answerGroup.Use(middlewares.AuthMiddleware())
	{

		answerGroup.POST("/CreateAnswer", answerController.CreateAnswer)
		answerGroup.POST("/GetAnswer", answerController.GetAnswer)
		answerGroup.GET("/GetAllAnswer", answerController.GetAllAnswers)
		answerGroup.PUT("/UpdateAnswer", answerController.UpdateAnswer)
		answerGroup.DELETE("/DeleteAnswer", answerController.DeleteAnswer)
		answerGroup.POST("/GetAnswersByQuestion", answerController.GetAnswersByQuestion)
	}
	// Article Routes
	articleController := controllers.NewArticleController(db)
	ArticleGroup := router.Group("/articles")
	ArticleGroup.Use(middlewares.AuthMiddleware())
	{
		ArticleGroup.GET("/all", articleController.GetAllArticles)
		ArticleGroup.POST("/get", articleController.GetArticle)
		ArticleGroup.POST("/GetArticlesByCourse", articleController.GetArticlesByCourse)
		ArticleGroup.POST("/createArticle", articleController.CreateArticle)
		ArticleGroup.PUT("/updateArticle", articleController.UpdateArticle)
		ArticleGroup.DELETE("/DeleteArticle", articleController.DeleteArticle)
	}
	// Category Routes
	CategoryController := controllers.NewCategoryController(db)
	CategoryGroup := router.Group("/categories")
	CategoryGroup.Use(middlewares.AuthMiddleware())
	{
		CategoryGroup.GET("/all", CategoryController.GetAllCategories)
		CategoryGroup.POST("/get", CategoryController.GetCategory)
		CategoryGroup.POST("/createCategory", CategoryController.CreateCategory)
		CategoryGroup.PUT("/updateCategory", CategoryController.UpdateCategory)
		CategoryGroup.DELETE("/DeleteCategory", CategoryController.DeleteCategory)
	}
	// Course Routes
	CourseController := controllers.NewCourseController(db)
	CoursesGroup := router.Group("/Courses")

	CoursesGroup.Use(middlewares.AuthMiddleware())
	{
		CoursesGroup.GET("/all", CourseController.GetAllCourses)
		CoursesGroup.POST("/createCourse", CourseController.CreateCourse)
		CoursesGroup.PUT("/updateCourse", CourseController.UpdateCourse)
		CoursesGroup.DELETE("/DeleteCourse", CourseController.DeleteCourse)

	}
	// coursequizz Routes
	CourseQuizzController := controllers.NewCourseQuizzController(db)
	CourseQuizzGroup := router.Group("/coursequizzs")
	CourseQuizzGroup.Use(middlewares.AuthMiddleware())
	{
		CourseQuizzGroup.GET("/all", CourseQuizzController.GetAllQuizzes)
		CourseQuizzGroup.POST("/get", CourseQuizzController.GetQuizz)
		CourseQuizzGroup.POST("/createCourseQuizz", CourseQuizzController.CreateQuizz)
		CourseQuizzGroup.PUT("/updateCourseQuizz", CourseQuizzController.UpdateQuizz)
		CourseQuizzGroup.DELETE("/DeleteCourseQuizz", CourseQuizzController.DeleteQuizz)
		ArticleGroup.POST("/GetQuizzesByCourse", CourseQuizzController.GetQuizzesByCourse)
	}
	// crating Routes
	CratingController := controllers.NewCratingController(db)
	CratingGroup := router.Group("/cratings")
	CratingGroup.Use(middlewares.AuthMiddleware())
	{
		CratingGroup.POST("/GetCratingsByCourse", CratingController.GetCratingsByCourse)
		CratingGroup.POST("/GetCratingsByStudent", CratingController.GetCratingsByStudent)
		CratingGroup.POST("/GetCratingByCourseAndStudent", CratingController.GetCratingsByStudent)
		CratingGroup.POST("/createCrating", CratingController.CreateCrating)
		CratingGroup.PUT("/updateCrating", CratingController.UpdateCrating)
		CratingGroup.DELETE("/DeleteCrating", CratingController.DeleteCrating)
		CratingGroup.GET("/GetAllCratings", CratingController.GetAllCratings)
		CratingGroup.POST("/GetCourseAverageRating", CratingController.GetCourseAverageRating)
	}
	// Exam Routes
	ExamController := controllers.NewExamController(db)
	ExamGroup := router.Group("/exams")
	ExamGroup.Use(middlewares.AuthMiddleware())
	{
		ExamGroup.GET("/all", ExamController.GetAllExams)
		ExamGroup.POST("/get", ExamController.GetExam)
		ExamGroup.POST("/createExam", ExamController.CreateExam)
		ExamGroup.PUT("/updateExam", ExamController.UpdateExam)
		ExamGroup.DELETE("/DeleteExam", ExamController.DeleteExam)
		ExamGroup.POST("/GetExamsByCourse", ExamController.GetExamsByCourse)
	}

	// ExamQuiz Routes
	ExamQuizController := controllers.NewExamQuizzController(db)
	ExamQuizGroup := router.Group("/examquizzes")
	ExamQuizGroup.Use(middlewares.AuthMiddleware())
	{
		ExamQuizGroup.GET("/all", ExamQuizController.GetAllExamQuizzes)
		ExamQuizGroup.POST("/get", ExamQuizController.GetExamQuizz)
		ExamQuizGroup.POST("/GetExamQuizzesByExam", ExamQuizController.GetExamQuizzesByExam)
		ExamQuizGroup.POST("/createExamQuiz", ExamQuizController.CreateExamQuizz)
		ExamQuizGroup.PUT("/updateExamQuiz", ExamQuizController.UpdateExamQuizz)
		ExamQuizGroup.DELETE("/DeleteExamQuiz", ExamQuizController.DeleteExamQuizz)
	}
	// feedback Routes
	FeedbackQuizController := controllers.NewFeedbackController(db)
	FeedbackGroup := router.Group("/feedbacks")
	FeedbackGroup.Use(middlewares.AuthMiddleware())
	{
		FeedbackGroup.GET("/all", FeedbackQuizController.GetAllFeedbacks)
		FeedbackGroup.POST("/get", FeedbackQuizController.GetFeedback)
		FeedbackGroup.POST("/getFeedbacksByStudent", FeedbackQuizController.GetFeedbacksByStudent)
		FeedbackGroup.POST("/createFeedback", FeedbackQuizController.CreateFeedback)
		FeedbackGroup.PUT("/updateFeedback", FeedbackQuizController.UpdateFeedback)
		FeedbackGroup.DELETE("/DeleteFeedback", FeedbackQuizController.DeleteFeedback)
	}
	// Question Routes
	QuestionkQuizController := controllers.NewQuestionController(db)
	QuestionGroup := router.Group("/questions")
	QuestionGroup.Use(middlewares.AuthMiddleware())
	{
		QuestionGroup.GET("/all", QuestionkQuizController.GetAllQuestions)
		QuestionGroup.POST("/get", QuestionkQuizController.GetQuestion)
		QuestionGroup.POST("/createQuestion", QuestionkQuizController.CreateQuestion)
		QuestionGroup.PUT("/updateQuestion", QuestionkQuizController.UpdateQuestion)
		QuestionGroup.DELETE("/DeleteQuestion", QuestionkQuizController.DeleteQuestion)
	}

	// student_course Routes
	studentCourseController := controllers.NewStudentCourseController(db)
	StudentCourseGroup := router.Group("/student_courses")
	StudentCourseGroup.Use(middlewares.AuthMiddleware())
	{
		StudentCourseGroup.GET("/all", studentCourseController.GetAllStudentCourses)
		StudentCourseGroup.POST("/get", studentCourseController.GetStudentCourse)
		StudentCourseGroup.POST("/SubmitExamAnswers", studentCourseController.SubmitExamAnswers)
		StudentCourseGroup.POST("/createStudentCourse", studentCourseController.CreateStudentCourse)
		StudentCourseGroup.PUT("/updateStudentCourse", studentCourseController.UpdateStudentCourse)
		StudentCourseGroup.DELETE("/DeleteStudentCourse", studentCourseController.DeleteStudentCourse)
	}
	// Student Routes
	StudentCourseController := controllers.NewStudentController(db)
	StudentGroup := router.Group("/students")
	StudentGroup.POST("/CreateStudent", StudentCourseController.CreateStudent)
	StudentGroup.Use(middlewares.AuthMiddleware())
	{
		StudentGroup.GET("/all", StudentCourseController.GetAllStudents)
		StudentGroup.POST("/GetStudent/:id", StudentCourseController.GetStudent)
		StudentGroup.PUT("/UpdateUser", StudentCourseController.UpdateStudent)
		StudentGroup.DELETE("/DeleteUser", StudentCourseController.DeleteStudent)
	}
	// Teacher Routes
	TeacherCourseController := controllers.NewTeacherController(db)
	TeacherGroup := router.Group("/teachers")
	TeacherGroup.POST("/login", controllers.GenerateToken)
	TeacherGroup.POST("/CreateTeacher", TeacherCourseController.CreateTeacher)
	TeacherGroup.Use(middlewares.AuthMiddleware())
	{
		TeacherGroup.GET("/all", TeacherCourseController.GetAllTeachers)
		TeacherGroup.POST("/GetTeacher", TeacherCourseController.GetTeacher)
		TeacherGroup.PUT("/UpdateTeacher", TeacherCourseController.UpdateTeacher)
		TeacherGroup.DELETE("/DeleteTeacher", TeacherCourseController.DeleteTeacher)
	}

}
