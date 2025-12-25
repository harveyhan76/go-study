package handlers

import (
	"blog/models"
	"blog/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{db: db}
}

// CreatePostRequest 创建文章请求结构体
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// UpdatePostRequest 更新文章请求结构体
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// CreatePost 创建文章
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := h.db.Create(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to create post")
		return
	}

	utils.Success(c, gin.H{
		"id":         post.ID,
		"title":      post.Title,
		"content":    post.Content,
		"user_id":    post.UserID,
		"created_at": post.CreatedAt,
	})
}

// GetPosts 获取文章列表
func (h *PostHandler) GetPosts(c *gin.Context) {
	var posts []models.Post

	// 预加载用户信息
	if err := h.db.Preload("User").Order("created_at DESC").Find(&posts).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch posts")
		return
	}

	var response []gin.H
	for _, post := range posts {
		response = append(response, gin.H{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
			"user": gin.H{
				"id":       post.User.ID,
				"username": post.User.Username,
			},
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		})
	}

	utils.Success(c, response)
}

// GetPost 获取单篇文章
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := h.db.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Post not found")
		} else {
			utils.InternalServerError(c, "Failed to fetch post")
		}
		return
	}

	// 构建评论响应
	var comments []gin.H
	for _, comment := range post.Comments {
		comments = append(comments, gin.H{
			"id":      comment.ID,
			"content": comment.Content,
			"user": gin.H{
				"id":       comment.User.ID,
				"username": comment.User.Username,
			},
			"created_at": comment.CreatedAt,
		})
	}

	utils.Success(c, gin.H{
		"id":      post.ID,
		"title":   post.Title,
		"content": post.Content,
		"user": gin.H{
			"id":       post.User.ID,
			"username": post.User.Username,
		},
		"comments":   comments,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
	})
}

// UpdatePost 更新文章
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Post not found")
		} else {
			utils.InternalServerError(c, "Failed to fetch post")
		}
		return
	}

	// 检查权限
	if post.UserID != userID {
		utils.Forbidden(c, "You can only update your own posts")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := h.db.Save(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to update post")
		return
	}

	utils.Success(c, gin.H{
		"id":         post.ID,
		"title":      post.Title,
		"content":    post.Content,
		"updated_at": post.UpdatedAt,
	})
}

// DeletePost 删除文章
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Post not found")
		} else {
			utils.InternalServerError(c, "Failed to fetch post")
		}
		return
	}

	// 检查权限
	if post.UserID != userID {
		utils.Forbidden(c, "You can only delete your own posts")
		return
	}

	if err := h.db.Delete(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete post")
		return
	}

	utils.Success(c, gin.H{
		"message": "Post deleted successfully",
	})
}
