package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/contracts"
	"github.com/formal-you/clean-architecture-blog/internal/application/repository"
	"github.com/formal-you/clean-architecture-blog/internal/errorx"
)

// ArticleUsecaseInterface 定义了文章相关的业务逻辑接口
type ArticleUsecaseInterface interface {
	CreateArticle(ctx context.Context, article *domain.Article) error
	GetArticleByID(ctx context.Context, id int64) (*domain.Article, error)
	GetAllArticles(ctx context.Context) ([]*domain.Article, error)
	UpdateArticle(ctx context.Context, article *domain.Article) error
	DeleteArticle(ctx context.Context, id int64) error
}

// ArticleUsecase 封装了文章相关的业务用例
type ArticleUsecase struct {
	repo        repository.ArticleRepository
	cache       repository.ArticleCacheRepository
	authService contracts.AuthService
	logger      contracts.Logger
}

// NewArticleUsecase 创建一个新的 ArticleUsecase
func NewArticleUsecase(repo repository.ArticleRepository, cache repository.ArticleCacheRepository, authService contracts.AuthService, logger contracts.Logger) ArticleUsecaseInterface {
	return &ArticleUsecase{
		repo:        repo,
		cache:       cache,
		authService: authService,
		logger:      logger,
	}
}

// CreateArticle 创建一篇新文章
func (uc *ArticleUsecase) CreateArticle(ctx context.Context, article *domain.Article) error {
	// 从 context 中获取 userID
	userID, err := uc.authService.GetUserIDFromContext(ctx)
	if err != nil {
		uc.logger.Error("failed to get user ID from context", "error", err)
		return errorx.New(errorx.CodeUnauthorized, err)
	}
	article.AuthorID = userID

	// 1. 验证领域实体
	if err := article.Validate(); err != nil {
		uc.logger.Error("article validation failed", "error", err)
		return errorx.New(errorx.CodeInvalidParams, err)
	}

	// 2. 调用仓库进行持久化
	err = uc.repo.Create(ctx, article)
	if err != nil {
		uc.logger.Error("failed to create article", "error", err)
		return errorx.New(errorx.CodeInternalServerError, err)
	}

	uc.logger.Info("article created successfully", "article_id", article.ID)
	return nil
}

// GetArticleByID 获取单篇文章
func (uc *ArticleUsecase) GetArticleByID(ctx context.Context, id int64) (*domain.Article, error) {
	// 1. 尝试从缓存获取
	cachedArticle, err := uc.cache.GetArticle(ctx, uint(id))
	if err != nil {
		uc.logger.Error("failed to get article from cache", "error", err)
		// 如果缓存出错，我们选择忽略并继续从数据库中获取，而不是直接返回错误
	}
	if cachedArticle != nil {
		uc.logger.Info("article cache hit", "article_id", id)
		return cachedArticle, nil
	}

	uc.logger.Info("article cache miss", "article_id", id)
	// 2. 缓存未命中，从数据库获取
	article, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errorx.New(errorx.CodeArticleNotFound, err)
		}
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}

	// 3. 将结果存入缓存
	err = uc.cache.SetArticle(ctx, article, 5*time.Minute)
	if err != nil {
		uc.logger.Error("failed to set article to cache", "error", err)
		// 即使缓存设置失败，我们仍然返回从数据库中获取到的数据
	}

	return article, nil
}

// GetAllArticles 获取所有文章
func (uc *ArticleUsecase) GetAllArticles(ctx context.Context) ([]*domain.Article, error) {
	const articlesCacheKey = "articles:all"
	// 1. 尝试从缓存获取
	cachedArticles, err := uc.cache.GetArticles(ctx, articlesCacheKey)
	if err != nil {
		uc.logger.Error("failed to get articles from cache", "error", err)
	}
	if cachedArticles != nil {
		uc.logger.Info("articles cache hit")
		return cachedArticles, nil
	}

	uc.logger.Info("articles cache miss")
	// 2. 缓存未命中，从数据库获取
	articles, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}

	// 3. 将结果存入缓存
	err = uc.cache.SetArticles(ctx, articlesCacheKey, articles, 5*time.Minute)
	if err != nil {
		uc.logger.Error("failed to set articles to cache", "error", err)
	}

	return articles, nil
}

// UpdateArticle 更新文章
func (uc *ArticleUsecase) UpdateArticle(ctx context.Context, article *domain.Article) error {
	// 从 context 中获取 userID
	userID, err := uc.authService.GetUserIDFromContext(ctx)
	if err != nil {
		return errorx.New(errorx.CodeUnauthorized, err)
	}

	// 检查用户是否有权限修改文章
	existingArticle, err := uc.repo.GetByID(ctx, article.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errorx.New(errorx.CodeArticleNotFound, err)
		}
		return errorx.New(errorx.CodeInternalServerError, err)
	}
	if existingArticle.AuthorID != userID {
		return errorx.New(errorx.CodeUnauthorized, errors.New("user not authorized to update this article"))
	}

	err = uc.repo.Update(ctx, article)
	if err != nil {
		return errorx.New(errorx.CodeInternalServerError, err)
	}

	// 更新成功后，删除缓存
	return uc.cache.DeleteArticle(ctx, uint(article.ID))
}

// DeleteArticle 删除文章
func (uc *ArticleUsecase) DeleteArticle(ctx context.Context, id int64) error {
	// 从 context 中获取 userID
	userID, err := uc.authService.GetUserIDFromContext(ctx)
	if err != nil {
		return errorx.New(errorx.CodeUnauthorized, err)
	}

	// 检查用户是否有权限删除文章
	existingArticle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errorx.New(errorx.CodeArticleNotFound, err)
		}
		return errorx.New(errorx.CodeInternalServerError, err)
	}
	if existingArticle.AuthorID != userID {
		return errorx.New(errorx.CodeUnauthorized, errors.New("user not authorized to delete this article"))
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return errorx.New(errorx.CodeInternalServerError, err)
	}
	// 删除成功后，删除缓存
	return uc.cache.DeleteArticle(ctx, uint(id))
}
