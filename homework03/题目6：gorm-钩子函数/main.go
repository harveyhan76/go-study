package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 模型 - 用户
type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Email     string
	Password  string
	PostCount int `gorm:"default:0"` // 新增：文章数量统计字段
	CreatedAt time.Time
	UpdatedAt time.Time

	// 一对多关系：一个用户有多篇文章
	Posts []Post `gorm:"foreignKey:UserID"`
}

// Post 模型 - 文章
type Post struct {
	ID            uint `gorm:"primaryKey"`
	Title         string
	Content       string
	CommentStatus string `gorm:"default:'有评论'"` // 新增：评论状态字段
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// 外键：关联用户
	UserID uint `gorm:"not null;index"`
	// Belongs To 关系：文章属于用户
	User User `gorm:"foreignKey:UserID" `

	// 一对多关系：一篇文章有多个评论
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// Comment 模型 - 评论
type Comment struct {
	ID        uint `gorm:"primaryKey"`
	Content   string
	CreatedAt time.Time

	// 外键：关联文章
	PostID uint `gorm:"not null;index" `
	// Belongs To 关系：评论属于文章
	Post Post `gorm:"foreignKey:PostID" `

	// 外键：关联用户（评论作者）
	UserID uint `gorm:"not null;index" `
	// Belongs To 关系：评论属于用户
	User User `gorm:"foreignKey:UserID"`
}

// Post模型的钩子函数，在文章创建后自动更新用户的文章数量
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 更新用户的文章数量
	err = tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + ?", 1)).Error
	if err != nil {
		return err
	}
	fmt.Printf("文章创建成功，已更新用户ID=%d的文章数量\n", p.UserID)
	return nil
}

// Comment模型的钩子函数，在评论删除后检查文章的评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 查询当前文章的评论数量
	var commentCount int64
	err = tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount).Error
	if err != nil {
		return err
	}

	// 如果评论数量为0，更新文章的评论状态为"无评论"
	if commentCount == 0 {
		err = tx.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_status", "无评论").Error
		if err != nil {
			return err
		}
		fmt.Printf("评论删除成功，文章ID=%d已无评论，评论状态已更新为'无评论'\n", c.PostID)
	} else {
		fmt.Printf("评论删除成功，文章ID=%d还有%d条评论\n", c.PostID, commentCount)
	}
	return nil
}

func main() {
	fmt.Println("Hello, World!")
	db, err := gorm.Open(mysql.Open("root:demo12!@@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	//建表
	// err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	// if err != nil {
	// 	panic("failed to migrate database: " + err.Error())
	// }

	// fmt.Println("Database tables created successfully!")
	// insertTable(db)

	// // 查询用户发布的文章及评论
	// var userID uint = 1
	// queryUserPostsAndComments(db, userID)

	// // 查询评论数量最多的文章
	// queryPostWithMostComments(db)

	// 测试钩子函数
	testHooks(db)
}

// 测试钩子函数
func testHooks(db *gorm.DB) {
	fmt.Println("\n=== 测试钩子函数 ===")

	// 创建测试用户
	testUser := User{
		Username: "TestUser",
		Email:    "test@example.com",
		Password: "testpass",
	}
	db.Create(&testUser)
	fmt.Printf("创建测试用户: %s (ID: %d)\n", testUser.Username, testUser.ID)

	// 测试1：创建文章，验证用户文章数量更新
	fmt.Println("\n--- 测试1：创建文章，验证用户文章数量更新 ---")
	post1 := Post{
		Title:   "测试文章1",
		Content: "这是测试文章1的内容",
		UserID:  testUser.ID,
	}
	db.Create(&post1)
	fmt.Printf("创建文章: %s (ID: %d)\n", post1.Title, post1.ID)

	// 查询用户文章数量
	var userAfterPost User
	db.First(&userAfterPost, testUser.ID)
	fmt.Printf("用户 %s 的文章数量: %d\n", userAfterPost.Username, userAfterPost.PostCount)

	// 测试2：创建评论
	fmt.Println("\n--- 测试2：创建评论 ---")
	comment1 := Comment{
		Content: "这是第一条评论",
		PostID:  post1.ID,
		UserID:  testUser.ID,
	}
	db.Create(&comment1)
	fmt.Printf("创建评论: %s (ID: %d)\n", comment1.Content, comment1.ID)

	// 测试3：删除评论，验证评论状态更新
	fmt.Println("\n--- 测试3：删除评论，验证评论状态更新 ---")
	db.Delete(&comment1)
	fmt.Printf("删除评论ID: %d\n", comment1.ID)

	// 查询文章评论状态
	var postAfterDelete Post
	db.First(&postAfterDelete, post1.ID)
	fmt.Printf("文章 %s 的评论状态: %s\n", postAfterDelete.Title, postAfterDelete.CommentStatus)

	// 测试4：创建第二篇文章和多个评论
	fmt.Println("\n--- 测试4：创建第二篇文章和多个评论 ---")
	post2 := Post{
		Title:   "测试文章2",
		Content: "这是测试文章2的内容",
		UserID:  testUser.ID,
	}
	db.Create(&post2)
	fmt.Printf("创建文章: %s (ID: %d)\n", post2.Title, post2.ID)

	// 创建多个评论
	comments := []Comment{
		{Content: "评论1", PostID: post2.ID, UserID: testUser.ID},
		{Content: "评论2", PostID: post2.ID, UserID: testUser.ID},
		{Content: "评论3", PostID: post2.ID, UserID: testUser.ID},
	}

	db.Create(&comments)

	// 测试5：删除部分评论，验证评论状态
	fmt.Println("\n--- 测试5：删除部分评论，验证评论状态 ---")
	// 删除第一条评论
	db.Delete(&comments[0])
	fmt.Printf("删除评论ID: %d\n", comments[0].ID)

	// 查询文章评论状态（应该还是"有评论"）
	var postAfterPartialDelete Post
	db.First(&postAfterPartialDelete, post2.ID)
	fmt.Printf("文章 %s 的评论状态: %s\n", postAfterPartialDelete.Title, postAfterPartialDelete.CommentStatus)

	// 删除剩余评论
	for i := 1; i < len(comments); i++ {
		db.Delete(&comments[i])
		fmt.Printf("删除评论ID: %d\n", comments[i].ID)
	}

	// 查询最终文章评论状态（应该变为"无评论"）
	var postFinal Post
	db.First(&postFinal, post2.ID)
	fmt.Printf("文章 %s 的最终评论状态: %s\n", postFinal.Title, postFinal.CommentStatus)

	// 最终查询用户文章数量
	var finalUser User
	db.First(&finalUser, testUser.ID)
	fmt.Printf("最终用户 %s 的文章数量: %d\n", finalUser.Username, finalUser.PostCount)

	fmt.Printf("=== 钩子函数测试完成 ===\n")
}

func insertTable(db *gorm.DB) {
	// 插入用户数据
	users := []User{
		{Username: "Alice", Email: "alice@example.com", Password: "password1"},
		{Username: "Bob", Email: "bob@example.com", Password: "password2"},
		{Username: "Charlie", Email: "charlie@example.com", Password: "password3"},
	}
	db.Create(&users)

	// 插入文章数据
	posts := []Post{
		{Title: "Post 1", Content: "Content of Post 1", UserID: users[0].ID},
		{Title: "Post 2", Content: "Content of Post 2", UserID: users[1].ID},
		{Title: "Post 3", Content: "Content of Post 3", UserID: users[2].ID},
	}
	db.Create(&posts)

	// 插入评论数据
	comments := []Comment{
		{Content: "Comment 1 on Post 1", PostID: posts[0].ID, UserID: users[1].ID},
		{Content: "Comment 2 on Post 1", PostID: posts[0].ID, UserID: users[2].ID},
		{Content: "Comment 1 on Post 2", PostID: posts[1].ID, UserID: users[0].ID},
		{Content: "Comment 2 on Post 2", PostID: posts[1].ID, UserID: users[2].ID},
		{Content: "Comment 1 on Post 3", PostID: posts[2].ID, UserID: users[0].ID},
		{Content: "Comment 2 on Post 3", PostID: posts[2].ID, UserID: users[1].ID},
	}
	db.Create(&comments)

	fmt.Println("Data inserted successfully!")
}

func queryUserPostsAndComments(db *gorm.DB, userID uint) {
	var user User

	err := db.Preload("Posts.Comments").First(&user, userID).Error
	if err != nil {
		fmt.Printf("Error retrieving user and posts: %v\n", err)
		return
	}

	fmt.Printf("User: %s\n", user.Username)
	for _, post := range user.Posts {
		fmt.Printf("Post: %s\n", post.Title)
		for _, comment := range post.Comments {
			fmt.Printf("\tComment: %s\n", comment.Content)
		}
	}
}

func queryPostWithMostComments(db *gorm.DB) {
	var result struct {
		PostID       uint
		CommentCount int64
	}

	err := db.Debug().Model(&Comment{}).
		Select("post_id, count(*) as comment_count").
		Group("post_id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		fmt.Printf("Error retrieving post with most comments: %v\n", err)
		return
	}

	var post Post
	db.Preload("Comments").First(&post, result.PostID)
	fmt.Printf("Post with most comments: %s\n", post.Title)
	fmt.Printf("Comments: %d\n", result.CommentCount)
}
