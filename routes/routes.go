package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	/*
		CoursesGroup := router.Group("/Courses")

		CoursesGroup.Use(middlewares.AuthMiddleware())
		{
			CoursesGroup.GET("/all", controllers.GetAllPlayers)
			CoursesGroup.GET("/get", controllers.GetPlayer)
			CoursesGroup.POST("/createCourse", controllers.CreatePlayer)
			CoursesGroup.PUT("/updateCourse", controllers.UpdatePlayer)
			CoursesGroup.DELETE("/DeleteCourse", controllers.DeletePlayer)

		}

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

		// Article Routes
		ArticleGroup := router.Group("/articles")
		ArticleGroup.Use(middlewares.AuthMiddleware())
		{
			ArticleGroup.GET("/all", controllers.GetAllArticles)
			ArticleGroup.GET("/get", controllers.GetArticle)
			ArticleGroup.POST("/createArticle", controllers.CreateArticle)
			ArticleGroup.PUT("/updateArticle", controllers.UpdateArticle)
			ArticleGroup.DELETE("/DeleteArticle", controllers.DeleteArticle)
		}

		// Video Routes
		VideoGroup := router.Group("/videos")
		VideoGroup.Use(middlewares.AuthMiddleware())
		{
			VideoGroup.GET("/all", controllers.GetAllVideos)
			VideoGroup.GET("/get", controllers.GetVideo)
			VideoGroup.POST("/createVideo", controllers.CreateVideo)
			VideoGroup.PUT("/updateVideo", controllers.UpdateVideo)
			VideoGroup.DELETE("/DeleteVideo", controllers.DeleteVideo)
		}

		// Answer Routes
		AnswerGroup := router.Group("/answers")
		AnswerGroup.Use(middlewares.AuthMiddleware())
		{
			AnswerGroup.GET("/all", controllers.GetAllAnswers)
			AnswerGroup.GET("/get", controllers.GetAnswer)
			AnswerGroup.POST("/createAnswer", controllers.CreateAnswer)
			AnswerGroup.PUT("/updateAnswer", controllers.UpdateAnswer)
			AnswerGroup.DELETE("/DeleteAnswer", controllers.DeleteAnswer)
		}

		// Category Routes
		CategoryGroup := router.Group("/categories")
		CategoryGroup.Use(middlewares.AuthMiddleware())
		{
			CategoryGroup.GET("/all", controllers.GetAllCategories)
			CategoryGroup.GET("/get", controllers.GetCategory)
			CategoryGroup.POST("/createCategory", controllers.CreateCategory)
			CategoryGroup.PUT("/updateCategory", controllers.UpdateCategory)
			CategoryGroup.DELETE("/DeleteCategory", controllers.DeleteCategory)
		}

		// SubCat Routes
		SubCatGroup := router.Group("/subcats")
		SubCatGroup.Use(middlewares.AuthMiddleware())
		{
			SubCatGroup.GET("/all", controllers.GetAllSubCats)
			SubCatGroup.GET("/get", controllers.GetSubCat)
			SubCatGroup.POST("/createSubCat", controllers.CreateSubCat)
			SubCatGroup.PUT("/updateSubCat", controllers.UpdateSubCat)
			SubCatGroup.DELETE("/DeleteSubCat", controllers.DeleteSubCat)
		}

		// Question Routes
		QuestionGroup := router.Group("/questions")
		QuestionGroup.Use(middlewares.AuthMiddleware())
		{
			QuestionGroup.GET("/all", controllers.GetAllQuestions)
			QuestionGroup.GET("/get", controllers.GetQuestion)
			QuestionGroup.POST("/createQuestion", controllers.CreateQuestion)
			QuestionGroup.PUT("/updateQuestion", controllers.UpdateQuestion)
			QuestionGroup.DELETE("/DeleteQuestion", controllers.DeleteQuestion)
		}

		// Exam Routes
		ExamGroup := router.Group("/exams")
		ExamGroup.Use(middlewares.AuthMiddleware())
		{
			ExamGroup.GET("/all", controllers.GetAllExams)
			ExamGroup.GET("/get", controllers.GetExam)
			ExamGroup.POST("/createExam", controllers.CreateExam)
			ExamGroup.PUT("/updateExam", controllers.UpdateExam)
			ExamGroup.DELETE("/DeleteExam", controllers.DeleteExam)
		}

		// coursequizz Routes
		CourseQuizzGroup := router.Group("/coursequizzs")
		CourseQuizzGroup.Use(middlewares.AuthMiddleware())
		{
			CourseQuizzGroup.GET("/all", controllers.GetAllCourseQuizzs)
			CourseQuizzGroup.GET("/get", controllers.GetCourseQuizz)
			CourseQuizzGroup.POST("/createCourseQuizz", controllers.CreateCourseQuizz)
			CourseQuizzGroup.PUT("/updateCourseQuizz", controllers.UpdateCourseQuizz)
			CourseQuizzGroup.DELETE("/DeleteCourseQuizz", controllers.DeleteCourseQuizz)
		}

		// crating Routes
		CratingGroup := router.Group("/cratings")
		CratingGroup.Use(middlewares.AuthMiddleware())
		{
			CratingGroup.GET("/all", controllers.GetAllCratings)
			CratingGroup.GET("/get", controllers.GetCrating)
			CratingGroup.POST("/createCrating", controllers.CreateCrating)
			CratingGroup.PUT("/updateCrating", controllers.UpdateCrating)
			CratingGroup.DELETE("/DeleteCrating", controllers.DeleteCrating)
		}

		// feedback Routes
		FeedbackGroup := router.Group("/feedbacks")
		FeedbackGroup.Use(middlewares.AuthMiddleware())
		{
			FeedbackGroup.GET("/all", controllers.GetAllFeedbacks)
			FeedbackGroup.GET("/get", controllers.GetFeedback)
			FeedbackGroup.POST("/createFeedback", controllers.CreateFeedback)
			FeedbackGroup.PUT("/updateFeedback", controllers.UpdateFeedback)
			FeedbackGroup.DELETE("/DeleteFeedback", controllers.DeleteFeedback)
		}

		// ExamQuiz Routes
		ExamQuizGroup := router.Group("/examquizzes")
		ExamQuizGroup.Use(middlewares.AuthMiddleware())
		{
			ExamQuizGroup.GET("/all", controllers.GetAllExamQuizzes)
			ExamQuizGroup.GET("/get", controllers.GetExamQuiz)
			ExamQuizGroup.POST("/createExamQuiz", controllers.CreateExamQuiz)
			ExamQuizGroup.PUT("/updateExamQuiz", controllers.UpdateExamQuiz)
			ExamQuizGroup.DELETE("/DeleteExamQuiz", controllers.DeleteExamQuiz)
		}

		// student_course Routes
		StudentCourseGroup := router.Group("/student_courses")
		StudentCourseGroup.Use(middlewares.AuthMiddleware())
		{
			StudentCourseGroup.GET("/all", controllers.GetAllStudentCourses)
			StudentCourseGroup.GET("/get", controllers.GetStudentCourse)
			StudentCourseGroup.POST("/createStudentCourse", controllers.CreateStudentCourse)
			StudentCourseGroup.PUT("/updateStudentCourse", controllers.UpdateStudentCourse)
			StudentCourseGroup.DELETE("/DeleteStudentCourse", controllers.DeleteStudentCourse)
		}
	*/
}
