package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func main() {
	// fmt.Println("Hello, World!")

	db, err := sqlx.Open("mysql", "root:demo12!@@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 创建表
	createTable(db)

	// 插入数据
	//insertTable(db)

	// 查询部门为"技术部"的所有员工
	queryDepartmentEmployees(db)

	// 查询工资最高的员工
	queryHighestPaidEmployee(db)

}

// createTable creates the employees table if it does not exist.
func createTable(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS employees (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        department VARCHAR(100) NOT NULL,
        salary INT NOT NULL
    );`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}
	fmt.Println("表创建成功或已存在")
}

// insertTable inserts sample data into the employees table.
func insertTable(db *sqlx.DB) {
	query := `INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)`

	employees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 10000},
		{Name: "李四", Department: "市场部", Salary: 8000},
		{Name: "王五", Department: "技术部", Salary: 12000},
	}

	for _, emp := range employees {
		_, err := db.Exec(query, emp.Name, emp.Department, emp.Salary)
		if err != nil {
			log.Printf("插入数据失败: %v", err)
		}
	}
	fmt.Println("数据插入完成")
}

// 1. 查询部门为"技术部"的所有员工
func queryDepartmentEmployees(db *sqlx.DB) {
	var employees []Employee
	department := "技术部"

	// 使用sqlx的Select方法，自动映射到结构体切片
	err := db.Select(&employees,
		"SELECT id, name, department, salary FROM employees WHERE department = ?",
		department)

	if err != nil {
		log.Printf("查询部门员工失败: %v", err)
		return
	}

	fmt.Println(employees)
}

// 2. 查询工资最高的员工
func queryHighestPaidEmployee(db *sqlx.DB) {
	var employee Employee

	// 方法1：使用ORDER BY和LIMIT
	err := db.Get(&employee,
		"SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")

	// 方法2：使用子查询（如果有多个员工工资相同且都是最高）
	// err := db.Get(&employee,
	//     "SELECT id, name, department, salary FROM employees WHERE salary = (SELECT MAX(salary) FROM employees) LIMIT 1")

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			fmt.Println("没有找到员工信息")
		} else {
			log.Printf("查询最高工资员工失败: %v", err)
		}
		return
	}

	fmt.Println("\n工资最高的员工信息:")
	fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %d\n",
		employee.ID, employee.Name, employee.Department, employee.Salary)
}
