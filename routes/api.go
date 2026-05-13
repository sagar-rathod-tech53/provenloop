package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/internal/repositories"
	"github.com/sagar-rathod-tech53/provenloop/internal/services"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
	"github.com/sagar-rathod-tech53/provenloop/migrations"
)

func SetupServer() *gin.Engine {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize DB connection
	db, err := config.ConnectDB(&cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Run DB migrations
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.UserRepository{DB: db}
	otpRepo := repositories.OTPRepository{DB: db}
	// postRepo := &repositories.PostRepository{DB: db}                     // pointer matches PostService.Repo
	// jobRepo := &repositories.JobRepository{DB: db}                       // pointer matches JobService.Repo
	// userProfileRepo := &repositories.UserProfileRepository{DB: db}       // pointer matches UserProfileService.Repo
	// videoRepo := &repositories.VideoProfileRepository{DB: db}            // pointer matches VideoService.Repo
	// educationRepo := &repositories.UserEducationRepository{DB: db}       // pointer matches EducationService.Repo
	// userExperienceRepo := &repositories.UserExperienceRepository{DB: db} // pointer matches ExperienceService.Repo
	// postLikeRepo := &repositories.PostLikeRepository{DB: db}             // pointer matches PostLikeService.Repo
	// postCommentRepo := &repositories.PostCommentRepository{DB: db}       // pointer matches PostCommentService.Repo
	// followRepo := &repositories.FollowRepository{DB: db}                 // pointer matches FollowService.Repo
	// notificationRepo := &repositories.NotificationRepository{DB: db}     // pointer matches NotificationService.Repo

	// Initialize services
	authService := services.AuthService{
		DB:              db,
		UserRepository:  userRepo,
		OTPRepository:   otpRepo,
		TokenExpiration: 3600,
		OTPLifespan:     300,
	}
	// postService := services.PostService{Repo: postRepo}
	// jobService := services.JobService{Repo: jobRepo}                                                      // pointer matches JobService.Repo
	// userProfileService := services.UserProfileService{Repo: userProfileRepo}                              // pointer matches UserProfileService.Repo
	// videoService := services.VideoProfileService{Repo: videoRepo}                                         // pointer matches VideoService.Repo
	// educationService := services.UserEducationService{Repo: educationRepo}                                // pointer matches EducationService.Repo
	// userExperienceService := services.UserExperienceService{UserExperienceRepository: userExperienceRepo} // pointer matches ExperienceService.Repo
	// postLikeService := services.PostLikeService{PostLikeRepository: postLikeRepo}                         // pointer matches PostLikeService.Repo
	// postCommentService := services.PostCommentService{PostCommentRepository: postCommentRepo}             // pointer matches PostCommentService.Repo
	// followService := services.FollowService{FollowRepository: followRepo}                                 // pointer matches FollowService.Repo
	// notificationService := services.NotificationService{NotificationRepository: notificationRepo}         // pointer matches NotificationService.Repo

	// Initialize controllers
	authController := controllers.AuthController{AuthService: &authService}
	// postController := controllers.PostController{PostService: &postService}
	// jobController := controllers.JobController{JobService: &jobService}
	// userProfileController := controllers.UserProfileController{UserProfileService: &userProfileService}             // pointer matches UserProfileController.Service
	// videoProfileController := controllers.VideoProfileController{VideoProfileService: &videoService}                // pointer matches VideoController.Service
	// educationController := controllers.UserEducationController{Service: &educationService}                          // pointer matches EducationController.Service
	// userExperienceController := controllers.UserExperienceController{UserExperienceService: &userExperienceService} // pointer matches ExperienceController.Service
	// postLikeController := controllers.PostLikeController{PostLikeService: &postLikeService}                         // pointer matches PostLikeController.Service
	// postCommentController := controllers.PostCommentController{PostCommentService: &postCommentService}             // pointer matches PostCommentController.Service
	// followController := controllers.FollowController{FollowService: &followService}                                 // pointer matches FollowController.Service
	// notificationController := controllers.NotificationController{NotificationService: &notificationService}         // pointer matches NotificationController.Service
	// uploadController, err := controllers.NewUploadController(cfg)
	// if err != nil {
	// 	log.Fatalf("Failed to create upload controller: %v", err)
	// }

	// Set up Gin router
	router := gin.Default()

	// Public API routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/verify-otp", authController.VerifyOTP)
		authGroup.POST("/forgot-password", authController.ForgotPassword)
		authGroup.POST("/reset-password", authController.ResetPassword)
	}

	// Protected auth routes
	authProtected := router.Group("/auth")
	authProtected.Use(middlewares.DeserializeUser(db))
	{
		authProtected.GET("/logout", authController.LogoutUser)
	}

	// // Protected post routes
	// postGroup := router.Group("/posts")
	// postGroup.Use(middlewares.DeserializeUser(db))
	// {
	// 	postGroup.POST("/content", postController.CreatePost)
	// 	postGroup.GET("/user/:user_id", postController.GetPostsByUserID)
	// 	postGroup.GET("/all-content", postController.GetAllContentPosts)
	// 	postGroup.POST("/job", jobController.CreateJobPost)
	// 	postGroup.GET("/all-job", jobController.GetAllJobPosts)
	// }

	// userGroup := router.Group("/user")
	// userGroup.Use(middlewares.DeserializeUser(db))
	// {
	// 	userGroup.POST("/profile", userProfileController.Create)
	// 	userGroup.GET("/profile/:user_id", userProfileController.GetByUserID)
	// 	userGroup.GET("/profile", userProfileController.GetAll)
	// 	userGroup.PUT("/profile/update/:user_id", userProfileController.Update)
	// 	userGroup.DELETE("/profile/delete/:user_id", userProfileController.Delete)
	// 	userGroup.POST("/video", videoProfileController.UploadVideo)
	// 	userGroup.GET("/video/:user_id", videoProfileController.GetVideoProfilesByUser)
	// 	userGroup.PUT("/video/:id", videoProfileController.UpdateVideo)    // New
	// 	userGroup.DELETE("/video/:id", videoProfileController.DeleteVideo) // New
	// 	userGroup.GET("/stream", videoProfileController.StreamVideo)
	// 	userGroup.POST("/upload", uploadController.UploadFile)
	// 	userGroup.POST("/education", educationController.Create)
	// 	userGroup.GET("/education/:user_id", educationController.GetByUser)
	// 	userGroup.PUT("/education/:id", educationController.Update)
	// 	userGroup.DELETE("/education/:id", educationController.Delete)
	// 	userGroup.POST("/experience", userExperienceController.Create)
	// 	userGroup.GET("/experience/:user_id", userExperienceController.GetByUserID)
	// 	userGroup.PUT("/experience/:id", userExperienceController.Update)
	// 	userGroup.DELETE("/experience/:id", userExperienceController.Delete)

	// }

	// likeGroup := router.Group("/post")
	// likeGroup.Use(middlewares.DeserializeUser(db))
	// {
	// 	// Routes for post likes
	// 	likeGroup.POST("/:post_id/like", postLikeController.LikePost)
	// 	likeGroup.POST("/:post_id/unlike", postLikeController.UnlikePost)
	// 	likeGroup.GET("/:post_id/likes", postLikeController.GetPostLikes)

	// }

	// commentGroup := router.Group("/post")
	// commentGroup.Use(middlewares.DeserializeUser(db))
	// {
	// 	// Routes for post comments
	// 	commentGroup.POST("/:post_id/comment", postCommentController.CommentOnPost)
	// 	commentGroup.GET("/:post_id/comments", postCommentController.GetPostComments)
	// }

	// follorshipGroup := router.Group("/user")
	// follorshipGroup.Use(middlewares.DeserializeUser(db))
	// {
	// 	// Routes for following and unfollowing
	// 	follorshipGroup.POST("/:followed_id/follow", followController.FollowUser)
	// 	follorshipGroup.POST("/:followed_id/unfollow", followController.UnfollowUser)
	// 	follorshipGroup.GET("/:user_id/followers", followController.GetFollowers)
	// 	follorshipGroup.GET("/:user_id/followings", followController.GetFollowings)
	// }

	// notificationGroup := router.Group("/notifications")
	// notificationGroup.Use(middlewares.DeserializeUser(db))
	// {
	// 	notificationGroup.POST("/create", notificationController.CreateNotification)
	// 	notificationGroup.GET("/:user_id", notificationController.GetNotifications)
	// }

	return router
}

func RunServer() {
	router := SetupServer()
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
