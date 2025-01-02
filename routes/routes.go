package routes

import (
	"database/sql"

	"github.com/cuddest/dz-skills/controllers"
	"github.com/cuddest/dz-skills/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine, db *sql.DB) {
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



		// Video Routes
		VideoGroup := router.Group("/videos")
		VideoGroup.Use(middlewares.AuthMiddleware())
		{
			VideoGroup.GET("/all", controllers.GetAllVideos)
			VideoGroup.GET("/get", controllers.GetVideo)
			VideoGroup.POST("/createVideo", controllers.CreateVideo)
			VideoGroup.PUT("/updateVideo", controllers.UpdateVideo)
			VideoGroup.DELETE("/DeleteVideo", controllers.DeleteVideo)
		}*/

	// Answer Routes
	answerController := controllers.NewAnswerController(db)
	answerGroup := router.Group("/answers")
	answerGroup.Use(middlewares.AuthMiddleware())
	{

		answerGroup.POST("/", answerController.CreateAnswer)                  // Create new answer
		answerGroup.POST("/GetAnswer", answerController.GetAnswer)            // Get specific answer
		answerGroup.GET("/GetAllAnswer", answerController.GetAllAnswers)      // Get all answers
		answerGroup.PUT("/UpdateAnswer", answerController.UpdateAnswer)       // Update specific answer
		answerGroup.DELETE("/DeleteAnswer", answerController.DeleteAnswer)    // Delete specific answer
		answerGroup.GET("/AnswerById", answerController.GetAnswersByQuestion) // Get answers by question
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
		CategoryGroup.GET("/get", CategoryController.GetCategory)
		CategoryGroup.POST("/createCategory", CategoryController.CreateCategory)
		CategoryGroup.PUT("/updateCategory", CategoryController.UpdateCategory)
		CategoryGroup.DELETE("/DeleteCategory", CategoryController.DeleteCategory)
	}

	/*

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

	   }
	*/
}
