package handlers

import (
	"blog/models"
	"blog/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

// CreateCommentRequest 创建评论请求结构体
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := h.db.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Post not found")
		} else {
			utils.InternalServerError(c, "Failed to fetch post")
		}
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postID),
	}

	if err := h.db.Create(&comment).Error; err != nil {
		utils.InternalServerError(c, "Failed to create comment")
		return
	}

	utils.Success(c, gin.H{
		"id":         comment.ID,
		"content":    comment.Content,
		"user_id":    comment.UserID,
		"post_id":    comment.PostID,
		"created_at": comment.CreatedAt,
	})
}

// GetComments 获取文章评论
func (h *CommentHandler) GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := h.db.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Post not found")
		} else {
			utils.InternalServerError(c, "Failed to fetch post")
		}
		return
	}

	var comments []models.Comment
	if err := h.db.Preload("User").Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch comments")
		return
	}

	var response []gin.H
	for _, comment := range comments {
		response = append(response, gin.H{
			"id":      comment.ID,
			"content": comment.Content,
			"user": gin.H{
				"id":       comment.User.ID,
				"username": comment.User.Username,
			},
			"created_at": comment.CreatedAt,
		})
	}

	utils.Success(c, response)
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		utils.BadRequest(c, "Invalid comment ID")
		return
	}

	var comment models.Comment
	if err := h.db.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Comment not found")
		} else {
			utils.InternalServerError(c, "Failed to fetch comment")
		}
		return
	}

	// 检查权限：只能删除自己的评论
	if comment.UserID != userID {
		utils.Forbidden(c, "You can only delete your own comments")
		return
	}

	if err := h.db.Delete(&comment).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete comment")
		return
	}

	utils.Success(c, gin.H{
		"message": "Comment deleted successfully",
	})
}
