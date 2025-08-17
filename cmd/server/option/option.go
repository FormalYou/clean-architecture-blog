package option

import (
	"log"
	"time"

	"github.com/FormalYou/clean-architecture-blog/internal/application/usecase"
	"github.com/FormalYou/clean-architecture-blog/internal/infrastructure/auth"
	"github.com/FormalYou/clean-architecture-blog/internal/infrastructure/cache"
	"github.com/FormalYou/clean-architecture-blog/internal/infrastructure/config"
	zaplog "github.com/FormalYou/clean-architecture-blog/internal/infrastructure/log"
	gorm_infra "github.com/FormalYou/clean-architecture-blog/internal/infrastructure/persistence/gorm"
	"github.com/FormalYou/clean-architecture-blog/internal/interfaces/http/handler"
	"github.com/FormalYou/clean-architecture-blog/internal/interfaces/http/handler/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(configPath string) *gin.Engine {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize Logger
	zaplog.InitLogger()
	zapLogger := zaplog.GetLogger()
	// defer zapLogger.Sync() // Sync will be called in main
	logger := zaplog.NewZapAdapter(zapLogger)

	// 2. Initialize Database
	dsnCfg := gorm_infra.DSNConfig{
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		DBName:   cfg.Database.DBName,
	}
	db, err := gorm_infra.NewDB(dsnCfg)
	if err != nil {
		zapLogger.Fatal("could not connect to db", zap.Error(err))
	}

	// 4. Dependency Injection
	redisClient, err := cache.NewRedisClient()
	if err != nil {
		zapLogger.Fatal("could not connect to redis", zap.Error(err))
	}

	jwtAuth := auth.NewJWTAuthService(cfg.JWT.Secret)
	jwtExpires := time.Duration(cfg.JWT.ExpiresInMinutes) * time.Minute

	articleRepo := gorm_infra.NewGormArticleRepository(db)
	articleCacheRepo := cache.NewArticleCacheRepository(redisClient)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, articleCacheRepo, jwtAuth, logger)
	articleHandler := handler.NewArticleHandler(articleUsecase, zapLogger)

	userRepo := gorm_infra.NewGormUserRepository(db)
	// auditSvc := usecase.NewAuditService(logger)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtAuth, jwtExpires, logger)
	userHandler := handler.NewUserHandler(userUsecase, zapLogger)

	commentRepo := gorm_infra.NewGormCommentRepository(db)
	_ = commentRepo // Placeholder for future use
	tagRepo := gorm_infra.NewGormTagRepository(db)
	_ = tagRepo // Placeholder for future use

	authMiddleware := middleware.AuthMiddleware(jwtAuth, zapLogger)
	errorHandler := middleware.ErrorHandler(zapLogger)

	// 5. Setup Router
	router := gin.Default()
	router.Use(errorHandler)

	v1 := router.Group("/api/v1")
	{
		// User routes
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)

		articles := v1.Group("/articles")
		{
			articles.GET("", articleHandler.GetAll)
			articles.GET("/:id", articleHandler.GetByID)

			authorized := articles.Group("/")
			authorized.Use(authMiddleware)
			{
				authorized.POST("", articleHandler.Create)
				authorized.PUT("/:id", articleHandler.Update)
				authorized.DELETE("/:id", articleHandler.Delete)
			}
		}
	}
	return router
}
