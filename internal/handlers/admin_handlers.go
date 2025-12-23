package handlers

import (
	"net/http"
	"strconv"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"github.com/Yw332/campus-moments-go/pkg/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var adminService = service.NewAdminService()

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	admin, err := adminService.Authenticate(req.Account, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "账号或密码错误",
			"data":    nil,
		})
		return
	}

	// 生成 token（admin 不依赖 userId，使用 "0"）
	token, _ := jwt.GenerateToken("0", admin.Username)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
		},
	})
}

// AdminMenu 返回管理员导航菜单
func AdminMenu(c *gin.Context) {
	menu := []gin.H{
		{
			"name": "用户管理",
			"path": "/admin/users",
			"children": []gin.H{
				{"name": "用户列表", "path": "/admin/users"},
				{"name": "添加用户", "path": "/admin/users/create"},
			},
		},
		{
			"name": "内容管理",
			"path": "/admin/contents",
			"children": []gin.H{
				{"name": "内容列表", "path": "/admin/contents"},
				{"name": "违规内容", "path": "/admin/contents/violations"},
			},
		},
		{
			"name": "系统设置",
			"path": "/admin/settings",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    menu,
	})
}

// AdminRegisterRequest 管理员注册请求
type AdminRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// AdminRegister 管理员注册
func AdminRegister(c *gin.Context) {
	var req AdminRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "用户名已存在",
			"data":    nil,
		})
		return
	}

	// 检查手机号是否已存在
	if err := db.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "手机号已注册",
			"data":    nil,
		})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码加密失败",
			"data":    nil,
		})
		return
	}

	// 生成用户ID（简单实现，实际应该更复杂）
	userID := "admin_" + strconv.FormatInt(int64(len(req.Username)), 10) + strconv.FormatInt(int64(len(req.Phone)), 10)

	// 创建管理员用户
	admin := models.User{
		ID:       userID,
		Username: req.Username,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Status:   1, // 正常状态
	}

	if err := db.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "注册失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"userId":   admin.ID,
			"username": admin.Username,
			"phone":    admin.Phone,
		},
	})
}

// UserListRequest 用户列表查询参数
type UserListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []UserInfo `json:"users"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
	TotalPages int        `json:"totalPages"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID             string `json:"userId"`
	Username       string `json:"username"`
	Phone          string `json:"phone"`
	Avatar         string `json:"avatar"`
	AvatarType     int    `json:"avatarType"`
	AvatarUpdatedAt string `json:"avatarUpdatedAt"`
	Signature      string `json:"signature"`
	Status         int    `json:"status"`
	CreatedAt      string `json:"createdAt"`
}

// ListUsers 获取用户列表
func ListUsers(c *gin.Context) {
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var users []models.User
	var total int64

	// 构建查询条件
	query := db.Model(&models.User{})
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ? OR id LIKE ?", 
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 获取分页数据
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询用户列表失败",
			"data":    nil,
		})
		return
	}

	// 转换数据格式
	userInfos := make([]UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = UserInfo{
			ID:             user.ID,
			Username:       user.Username,
			Phone:          user.Phone,
			Avatar:         user.Avatar,
			AvatarType:     user.AvatarType,
			Signature:      user.Signature,
			Status:         user.Status,
			CreatedAt:      user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		
		// 格式化 AvatarUpdatedAt
		if user.AvatarUpdatedAt != nil {
			userInfos[i].AvatarUpdatedAt = user.AvatarUpdatedAt.Format("2006-01-02 15:04:05")
		}
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	response := UserListResponse{
		Users:      userInfos,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    response,
	})
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=6,max=50"`
	Confirm     bool   `json:"confirm" binding:"required"`
}

// ResetUserPassword 重置用户密码
func ResetUserPassword(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 验证二次确认
	if !req.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请确认操作",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 检查用户是否存在
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码加密失败",
			"data":    nil,
		})
		return
	}

	// 更新密码
	if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码重置失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码重置成功",
		"data":    nil,
	})
}

// SetAdminRequest 设置管理员请求
type SetAdminRequest struct {
	IsAdmin bool `json:"isAdmin"`
}

// SetUserAsAdmin 设置用户为管理员
func SetUserAsAdmin(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	var req SetAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 检查用户是否存在
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 这里需要添加role字段到User模型，暂时使用status字段区分
	// 或者创建一个新的Admin模型来管理管理员权限
	// 现在简单实现：status=3表示管理员
	newStatus := 1 // 普通用户
	if req.IsAdmin {
		newStatus = 3 // 管理员
	}

	if err := db.Model(&user).Update("status", newStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "设置管理员失败",
			"data":    nil,
		})
		return
	}

	message := "取消管理员成功"
	if req.IsAdmin {
		message = "设置管理员成功"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": message,
		"data":    nil,
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	// 验证二次确认参数
	confirm := c.Query("confirm")
	if confirm != "true" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请确认删除操作（?confirm=true）",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 检查用户是否存在
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 开始事务
	tx := db.Begin()

	// 软删除用户（更新状态为已删除）
	if err := tx.Model(&user).Update("status", 0).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除用户失败",
			"data":    nil,
		})
		return
	}

	// 软删除用户的所有动态（更新状态为已删除）
	if err := tx.Model(&models.Moment{}).Where("author_id = ?", userID).Update("status", 2).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除用户动态失败",
			"data":    nil,
		})
		return
	}

	// 提交事务
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除用户成功",
		"data":    nil,
	})
}

// GetUserFriends 获取用户好友列表
func GetUserFriends(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	// 这里需要实现好友关系查询
	// 暂时返回空列表，因为当前模型中没有好友关系表
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"friends": []interface{}{},
			"total":   0,
		},
	})
}

// ========== 内容管理相关结构体 ==========

// ContentListRequest 内容列表查询参数
type ContentListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`    // 按内容关键词搜索
	Author   string `form:"author"`    // 按发布者搜索
}

// ContentListResponse 内容列表响应
type ContentListResponse struct {
	Contents   []ContentInfo `json:"contents"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
	TotalPages int           `json:"totalPages"`
}

// ContentInfo 内容信息
type ContentInfo struct {
	ID           int64         `json:"id"`
	Content      string        `json:"content"`
	AuthorID     string        `json:"authorId"`
	AuthorName   string        `json:"authorName"`
	AuthorAvatar string        `json:"authorAvatar"`
	Tags         []string      `json:"tags"`
	Media        []MediaItem   `json:"media"`
	LikeCount    int           `json:"likeCount"`
	CommentCount int           `json:"commentCount"`
	Status       int           `json:"status"`
	Visibility   int           `json:"visibility"`
	CreatedAt    string        `json:"createdAt"`
}

// ContentDetailResponse 内容详情响应
type ContentDetailResponse struct {
	ContentInfo
	Comments []CommentInfo `json:"comments"`
}

// CommentInfo 评论信息
type CommentInfo struct {
	ID            int64      `json:"id"`
	Content       string     `json:"content"`
	UserID        string     `json:"userId"`
	UserName      string     `json:"userName"`
	UserAvatar    string     `json:"userAvatar"`
	ParentID      *int64     `json:"parentId"`
	ReplyToUserID *string    `json:"replyToUserId"`
	ReplyToName   *string    `json:"replyToName"`
	LikeCount     int        `json:"likeCount"`
	Status        int        `json:"status"`
	CreatedAt     string     `json:"createdAt"`
}

// MediaItem 媒体项（复制自moment.go，避免循环导入）
type MediaItem struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration int    `json:"duration"`
}

// ListContents 获取内容列表
func ListContents(c *gin.Context) {
	var req ContentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 确保默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := database.GetDB()
	var moments []models.Moment
	var total int64

	// 构建查询条件
	query := db.Model(&models.Moment{}).Preload("Author")
	
	// 关键词搜索（内容）
	if req.Keyword != "" {
		query = query.Where("content LIKE ?", "%"+req.Keyword+"%")
	}
	
	// 发布者搜索
	if req.Author != "" {
		query = query.Where("author_id = ?", req.Author)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询总数失败",
			"data":    nil,
		})
		return
	}

	// 获取分页数据
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&moments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询内容列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 转换数据格式
	contents := make([]ContentInfo, len(moments))
	for i, moment := range moments {
		// 转换Tags
		tags := []string{}
		if moment.Tags != nil {
			tags = moment.Tags
		}
		
		// 转换Media
		media := []MediaItem{}
		if moment.Media != nil {
			for _, m := range moment.Media {
				media = append(media, MediaItem{
					URL:      m.URL,
					Type:     m.Type,
					Size:     m.Size,
					Width:    m.Width,
					Height:   m.Height,
					Duration: m.Duration,
				})
			}
		}

		authorName := ""
		authorAvatar := ""
		if moment.Author != nil {
			authorName = moment.Author.Username
			authorAvatar = moment.Author.Avatar
		}

		contents[i] = ContentInfo{
			ID:           moment.ID,
			Content:      moment.Content,
			AuthorID:     moment.AuthorID,
			AuthorName:   authorName,
			AuthorAvatar: authorAvatar,
			Tags:         tags,
			Media:        media,
			LikeCount:    moment.LikeCount,
			CommentCount: moment.CommentCount,
			Status:       moment.Status,
			Visibility:   moment.Visibility,
			CreatedAt:    moment.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	response := ContentListResponse{
		Contents:   contents,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    response,
	})
}

// GetContentDetail 获取内容详情
func GetContentDetail(c *gin.Context) {
	contentID := c.Param("id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "内容ID不能为空",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 获取动态详情
	var moment models.Moment
	if err := db.Preload("Author").Where("id = ?", contentID).First(&moment).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "内容不存在",
			"data":    nil,
		})
		return
	}

	// 获取评论列表
	var comments []models.Comment
	if err := db.Preload("User").Preload("ReplyToUser").Where("moment_id = ? AND status = 1", contentID).Order("created_at ASC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询评论失败",
			"data":    nil,
		})
		return
	}

	// 转换动态信息
	tags := []string{}
	if moment.Tags != nil {
		tags = moment.Tags
	}
	
	media := []MediaItem{}
	if moment.Media != nil {
		for _, m := range moment.Media {
			media = append(media, MediaItem{
				URL:      m.URL,
				Type:     m.Type,
				Size:     m.Size,
				Width:    m.Width,
				Height:   m.Height,
				Duration: m.Duration,
			})
		}
	}

	authorName := ""
	authorAvatar := ""
	if moment.Author != nil {
		authorName = moment.Author.Username
		authorAvatar = moment.Author.Avatar
	}

	contentInfo := ContentInfo{
		ID:           moment.ID,
		Content:      moment.Content,
		AuthorID:     moment.AuthorID,
		AuthorName:   authorName,
		AuthorAvatar: authorAvatar,
		Tags:         tags,
		Media:        media,
		LikeCount:    moment.LikeCount,
		CommentCount: moment.CommentCount,
		Status:       moment.Status,
		Visibility:   moment.Visibility,
		CreatedAt:    moment.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 转换评论信息
	commentInfos := make([]CommentInfo, len(comments))
	for i, comment := range comments {
		userName := ""
		userAvatar := ""
		if comment.User != nil {
			userName = comment.User.Username
			userAvatar = comment.User.Avatar
		}

		replyToName := (*string)(nil)
		if comment.ReplyToUser != nil {
			replyToName = &comment.ReplyToUser.Username
		}

		commentInfos[i] = CommentInfo{
			ID:            comment.ID,
			Content:       comment.Content,
			UserID:        comment.UserID,
			UserName:      userName,
			UserAvatar:    userAvatar,
			ParentID:      comment.ParentID,
			ReplyToUserID: comment.ReplyToUserID,
			ReplyToName:   replyToName,
			LikeCount:     comment.LikeCount,
			Status:        comment.Status,
			CreatedAt:     comment.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	response := ContentDetailResponse{
		ContentInfo: contentInfo,
		Comments:    commentInfos,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    response,
	})
}

// DeleteContent 删除内容
func DeleteContent(c *gin.Context) {
	contentID := c.Param("id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "内容ID不能为空",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 检查内容是否存在
	var moment models.Moment
	if err := db.Where("id = ?", contentID).First(&moment).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "内容不存在",
			"data":    nil,
		})
		return
	}

	// 开始事务
	tx := db.Begin()

	// 软删除动态（更新状态为已删除）
	if err := tx.Model(&moment).Update("status", 2).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除内容失败",
			"data":    nil,
		})
		return
	}

	// 软删除相关评论
	if err := tx.Model(&models.Comment{}).Where("moment_id = ?", contentID).Update("status", 2).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除相关评论失败",
			"data":    nil,
		})
		return
	}

	// 删除相关点赞
	if err := tx.Where("target_id = ? AND target_type = ?", contentID, models.LikeTargetTypeMoment).Delete(&models.Like{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除相关点赞失败",
			"data":    nil,
		})
		return
	}

	// 提交事务
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除内容成功",
		"data":    nil,
	})
}

// GetUserDetail 获取用户详情（管理员用）
func GetUserDetail(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 获取用户详情
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 获取用户动态数量
	var momentCount int64
	db.Model(&models.Moment{}).Where("author_id = ? AND status = 1", userID).Count(&momentCount)

	// 获取用户评论数量
	var commentCount int64
	db.Model(&models.Comment{}).Where("user_id = ? AND status = 1", userID).Count(&commentCount)

	userInfo := UserInfo{
		ID:             user.ID,
		Username:       user.Username,
		Phone:          user.Phone,
		Avatar:         user.Avatar,
		AvatarType:     user.AvatarType,
		Signature:      user.Signature,
		Status:         user.Status,
		CreatedAt:      user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	
	// 格式化 AvatarUpdatedAt
	if user.AvatarUpdatedAt != nil {
		userInfo.AvatarUpdatedAt = user.AvatarUpdatedAt.Format("2006-01-02 15:04:05")
	}

	response := gin.H{
		"user":         userInfo,
		"momentCount":  momentCount,
		"commentCount": commentCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    response,
	})
}