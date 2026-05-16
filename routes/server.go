package routes

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/internal/repositories"
	"github.com/sagar-rathod-tech53/provenloop/internal/services"
	"github.com/sagar-rathod-tech53/provenloop/migrations"
)

func SetupServer() *gin.Engine {

	// ==========================================
	// CONFIG
	// ==========================================
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// ==========================================
	// DB
	// ==========================================
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// ==========================================
	// MIGRATIONS
	// ==========================================
	if err := migrations.Migrate(db); err != nil {
		log.Fatal(err)
	}

	// ==========================================
	// REDIS
	// ==========================================
	config.ConnectRedis()

	// ==========================================
	// REPOSITORIES
	// ==========================================
	userRepo := &repositories.UserRepository{DB: db}

	profileRepo := &repositories.UserProfileRepository{DB: db}

	educationRepo := &repositories.EducationRepository{DB: db}

	experienceRepo := &repositories.ExperienceRepository{DB: db}

	projectRepo := &repositories.ProjectRepository{DB: db}

	skillRepo := &repositories.SkillRepository{DB: db}

	certificationRepo := &repositories.CertificationRepository{DB: db}

	// ==========================================
	// SERVICES
	// ==========================================
	authService := &services.AuthService{
		UserRepo: userRepo,
	}

	profileService := &services.UserProfileService{
		Repo: profileRepo,
	}

	educationService := &services.EducationService{
		Repo: educationRepo,
	}

	experienceService := &services.ExperienceService{
		Repo: experienceRepo,
	}

	projectService := &services.ProjectService{
		Repo: projectRepo,
	}

	skillService := &services.SkillService{
		Repo: skillRepo,
	}

	certificationService := &services.CertificationService{
		Repo: certificationRepo,
	}

	// ==========================================
	// CONTROLLERS
	// ==========================================
	authController := &controllers.AuthController{
		AuthService: authService,
	}

	profileController := &controllers.UserProfileController{
		Service: profileService,
	}

	educationController := &controllers.EducationController{
		Service: educationService,
	}

	experienceController := &controllers.ExperienceController{
		Service: experienceService,
	}

	projectController := &controllers.ProjectController{
		Service: projectService,
	}

	skillController := &controllers.SkillController{
		Service: skillService,
	}

	certificationController := &controllers.CertificationController{
		Service: certificationService,
	}

	// ==========================================
	// ROUTER
	// ==========================================
	router := gin.Default()

	// ==========================================
	// REGISTER ROUTES
	// ==========================================

	RegisterAuthRoutes(router, authController, db)
	SkillRoutes(router, skillController, db)
	ExperienceRoutes(router, experienceController, db)
	RegisterUserProfileRoutes(router, profileController, db)
	RegisterEducationRoutes(router, educationController, db)
	ProjectRoutes(router, projectController, db)
	RegisterCertificationRoutes(router, certificationController, db)
	return router
}
