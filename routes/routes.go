package routes

import (
	"database/sql"

	"github.com/cuddest/dz-skills/controllers"
	"github.com/cuddest/dz-skills/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine, db *sql.DB) {
	/*


		// Student Routes
		StudentGroup := router.Group("/students")
		StudentGroup.POST("/login", controllers.GenerateToken)
		StudentGroup.POST("/CreateUser", controllers.CreateUser)
		StudentGroup.Use(middlewares.AuthMiddleware())
		{
			StudentGroup.GET("/all", controllers.GetAllUsers)
			StudentGroup.POST("/GetUser", controllers.GetUser)
			StudentGroup.PUT("/UpdateUser", controllers.UpdateUser)
			StudentGroup.DELETE("/DeleteUser", controllers.DeleteUser)
		}

		// Teacher Routes
		TeacherGroup := router.Group("/teachers")
		TeacherGroup.POST("/login", controllers.GenerateToken)
		TeacherGroup.POST("/CreateTeacher", controllers.CreateTeacher)
		TeacherGroup.Use(middlewares.AuthMiddleware())
		{
			TeacherGroup.GET("/all", controllers.GetAllTeachers)
			TeacherGroup.POST("/GetTeacher", controllers.GetTeacher)
			TeacherGroup.PUT("/UpdateTeacher", controllers.UpdateTeacher)
			TeacherGroup.DELETE("/DeleteTeacher", controllers.DeleteTeacher)
		}



	*/

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
		answerGroup.POST("/AnswerById", answerController.GetAnswersByQuestion)
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
		CoursesGroup.POST("/get", CourseController.GetCourse)
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
	// SubCat Routes
	SubcattCourseController := controllers.NewSubCatController(db)
	SubCatGroup := router.Group("/subcats")
	SubCatGroup.Use(middlewares.AuthMiddleware())
	{
		SubCatGroup.GET("/all", SubcattCourseController.GetAllSubCats)
		SubCatGroup.POST("/get", SubcattCourseController.GetSubCat)
		SubCatGroup.POST("/createSubCat", SubcattCourseController.CreateSubCat)
		SubCatGroup.PUT("/updateSubCat", SubcattCourseController.UpdateSubCat)
		SubCatGroup.PUT("/updateSubCat", SubcattCourseController.GetSubCatsByCategory)
		SubCatGroup.DELETE("/DeleteSubCat", SubcattCourseController.DeleteSubCat)
	}
	// Video Routes
	VideoCourseController := controllers.NewVideoController(db)
	VideoGroup := router.Group("/videos")
	VideoGroup.Use(middlewares.AuthMiddleware())
	{
		VideoGroup.GET("/all", VideoCourseController.GetAllVideos)
		VideoGroup.POST("/get", VideoCourseController.GetVideo)
		VideoGroup.POST("/get", VideoCourseController.GetVideosByCourse)
		VideoGroup.POST("/createVideo", VideoCourseController.CreateVideo)
		VideoGroup.PUT("/updateVideo", VideoCourseController.UpdateVideo)
		VideoGroup.DELETE("/DeleteVideo", VideoCourseController.DeleteVideo)
	}

}
