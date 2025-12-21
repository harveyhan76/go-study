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
	CreatedAt time.Time
	UpdatedAt time.Time

	// 一对多关系：一个用户有多篇文章
	Posts []Post `gorm:"foreignKey:UserID"`
}

// Post 模型 - 文章
type Post struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time

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

func main() {
	fmt.Println("Hello, World!")
	db, err := gorm.Open(mysql.Open("root:demo12!@@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	//建表
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	fmt.Println("Database tables created successfully!")
	//insertTable(db)

	// 查询用户发布的文章及评论
	var userID uint = 1
	queryUserPostsAndComments(db, userID)

	// 查询评论数量最多的文章
	queryPostWithMostComments(db)
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
