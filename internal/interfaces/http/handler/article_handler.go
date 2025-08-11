package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/usecase"
	"github.com/formal-you/clean-architecture-blog/internal/errorx"
	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/dto"
)

// ArticleHandler handles HTTP requests for articles
type ArticleHandler struct {
	usecase *usecase.ArticleUsecase
	logger  *zap.Logger
}

// NewArticleHandler creates a new ArticleHandler
func NewArticleHandler(usecase *usecase.ArticleUsecase, logger *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		usecase: usecase,
		logger:  logger.Named("ArticleHandler"),
	}
}

// Create handles the creation of a new article
func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(errorx.CodeInvalidParams, err))
		return
	}

	article := &domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Tags:    make([]domain.Tag, len(req.Tags)),
	}
	for i, tagName := range req.Tags {
		article.Tags[i] = domain.Tag{Name: tagName}
	}

	if err := h.usecase.CreateArticle(c.Request.Context(), article); err != nil {
		c.Error(err)
		return
	}

	h.logger.Info("article created successfully", zap.Int64("article_id", article.ID))
	c.JSON(http.StatusCreated, article)
}

// GetByID handles retrieving a single article by its ID
func (h *ArticleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.CodeInvalidParams, err))
		return
	}

	article, err := h.usecase.GetArticleByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	h.logger.Info("article retrieved successfully", zap.Int64("article_id", article.ID))
	c.JSON(http.StatusOK, article)
}

// GetAll handles retrieving all articles
func (h *ArticleHandler) GetAll(c *gin.Context) {
	articles, err := h.usecase.GetAllArticles(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	h.logger.Info("retrieved all articles", zap.Int("count", len(articles)))
	c.JSON(http.StatusOK, articles)
}

// Update handles updating an existing article
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.CodeInvalidParams, err))
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(errorx.CodeInvalidParams, err))
		return
	}

	article := &domain.Article{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.usecase.UpdateArticle(c.Request.Context(), article); err != nil {
		c.Error(err)
		return
	}

	h.logger.Info("article updated successfully", zap.Int64("article_id", id))
	c.JSON(http.StatusOK, gin.H{"status": "article updated"})
}

// Delete handles deleting an article
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.CodeInvalidParams, err))
		return
	}

	if err := h.usecase.DeleteArticle(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	h.logger.Info("article deleted successfully", zap.Int64("article_id", id))
	c.JSON(http.StatusOK, gin.H{"status": "article deleted"})
}