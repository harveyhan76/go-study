package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Students struct {
	ID    int `gorm:"primaryKey;autoIncrement"`
	Name  string
	Age   int
	Grade string
}

// 1. 插入学生姓名为"张三"，年龄为20，年级为"三年级"
func insertStudent(db *gorm.DB) {
	student := &Students{Name: "张三", Age: 20, Grade: "三年级"}
	result := db.Create(student)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println(student)
}

// 2. 查询年龄大于18岁的学生
func queryStudents(db *gorm.DB) {
	var students []Students
	result := db.Where("age > ?", 18).Find(&students)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println(students)
}

// 3. 将姓名为"张三"的学生年级更新为"四年级"
func updateStudents(db *gorm.DB) {
	result := db.Model(&Students{}).
		Where("name = ?", "张三").
		Update("grade", "四年级")

	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println(result.RowsAffected)
}

// 4. 删除年龄小于15岁的学生记录
func deleteStudents(db *gorm.DB) {
	result := db.Where("age < ?", 15).Delete(&Students{})
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println(result.RowsAffected)
}

func main() {
	// fmt.Println("Hello, World!")
	db, err := gorm.Open(mysql.Open("root:demo12!@@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	//建表
	db.AutoMigrate(&Students{})

	// 执行各个操作
	//insertStudent(db) // 插入张三
	//queryStudents(db) // 查询年龄大于18岁的学生
	//updateStudents(db) // 更新张三的年级
	//deleteStudents(db) // 删除年龄小于15岁的学生记录
}
