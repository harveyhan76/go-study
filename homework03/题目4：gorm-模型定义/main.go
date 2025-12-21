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

}
