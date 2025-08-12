package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/formal-you/clean-architecture-blog/internal/application/usecase"
	"github.com/formal-you/clean-architecture-blog/internal/infrastructure/auth"
	"github.com/formal-you/clean-architecture-blog/internal/infrastructure/cache"
	"github.com/formal-you/clean-architecture-blog/internal/infrastructure/config"
	"github.com/formal-you/clean-architecture-blog/internal/infrastructure/log"
	gorm_db "github.com/formal-you/clean-architecture-blog/internal/infrastructure/persistence/gorm"
	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/handler"
	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/handler/middleware"
)

var (
	TestRouter  *gin.Engine
	testDB      *gorm.DB
	redisClient *redis.Client
	testLogger  *zap.Logger
	testConfig  config.Config
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	var err error
	testConfig, err = config.LoadConfig("../../../../configs")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	log.InitLogger()
	testLogger = log.GetLogger()

	dbConfig := gorm_db.DSNConfig{
		User:     testConfig.Database.User,
		Password: testConfig.Database.Password,
		Host:     testConfig.Database.Host,
		Port:     testConfig.Database.Port,
		DBName:   testConfig.Database.DBName,
	}
	testDB, err = gorm_db.NewDB(dbConfig)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	redisClient, err = cache.NewRedisClient()
	if err != nil {
		panic("failed to connect redis: " + err.Error())
	}

	TestRouter = SetupRouter(&testConfig, testDB, redisClient, testLogger)

	code := m.Run()

	CleanUp()

	os.Exit(code)
}

func SetupRouter(cfg *config.Config, db *gorm.DB, rds *redis.Client, logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger))

	userRepo := gorm_db.NewGormUserRepository(db)
	articleRepo := gorm_db.NewGormArticleRepository(db)
	articleCacheRepo := cache.NewArticleCacheRepository(rds)

	authService := auth.NewJWTAuthService(cfg.JWT.Secret)

	loggerAdapter := log.NewZapAdapter(logger)
	userUsecase := usecase.NewUserUsecase(userRepo, authService, time.Duration(cfg.JWT.ExpiresInMinutes)*time.Minute, loggerAdapter)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, articleCacheRepo, authService, loggerAdapter)

	userHandler := handler.NewUserHandler(userUsecase, logger)
	articleHandler := handler.NewArticleHandler(articleUsecase, logger)

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", userHandler.Register)
		userRoutes.POST("/login", userHandler.Login)
	}

	articleRoutes := router.Group("/articles")
	authMiddleware := middleware.AuthMiddleware(authService, logger)
	articleRoutes.Use(authMiddleware)
	{
		articleRoutes.POST("", articleHandler.Create)
		articleRoutes.GET("/:id", articleHandler.GetByID)
		articleRoutes.GET("", articleHandler.GetAll)
		articleRoutes.PUT("/:id", articleHandler.Update)
		articleRoutes.DELETE("/:id", articleHandler.Delete)
	}

	return router
}

func CleanUp() {
	testDB.Exec("DELETE FROM comments")
	testDB.Exec("DELETE FROM article_tags")
	testDB.Exec("DELETE FROM tag_models")
	testDB.Exec("DELETE FROM article_models")
	testDB.Exec("DELETE FROM user_models")

	redisClient.FlushAll(context.Background())
}

func PerformRequest(r http.Handler, method, path string, body io.Reader, token ...string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	if len(token) > 0 {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ToJSON(v interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		panic(err)
	}
	return &buf
}

func ClearAllData() {
	CleanUp()
}
