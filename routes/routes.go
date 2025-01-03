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

		answerGroup.POST("/CreateAnswer", answerController.CreateAnswer)      // Create new answer
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
		ExamGroup.GET("/get", ExamController.GetExam)
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
