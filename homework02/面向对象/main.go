package main

import (
	"fmt"
	"math"
)

// Shape 接口定义
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 结构体
type Rectangle struct {
	Width  float64
	Height float64
}

// Rectangle 实现 Shape 接口
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle 结构体
type Circle struct {
	Radius float64
}

// Circle 实现 Shape 接口
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 打印形状信息
func printShapeInfo(s Shape) {
	fmt.Printf("面积: %.2f\n", s.Area())
	fmt.Printf("周长: %.2f\n", s.Perimeter())
	fmt.Println("---")
}

// 题目：组合结构体

// Person 结构体
type Person struct {
	Name string
	Age  int
}

// Employee 结构体，组合 Person
type Employee struct {
	Person     // 匿名嵌入，实现组合
	EmployeeID string
	Department string // 可以添加更多字段
}

// 为 Employee 结构体实现 PrintInfo() 方法
func (e *Employee) PrintInfo() {
	fmt.Printf("员工信息：\n")
	fmt.Printf("  姓名：%s\n", e.Name) // 可以直接访问 Person 的字段
	fmt.Printf("  年龄：%d\n", e.Age)  // 可以直接访问 Person 的字段
	fmt.Printf("  工号：%s\n", e.EmployeeID)
	fmt.Printf("  部门：%s\n", e.Department)
}

func main() {
	// 创建 Rectangle 实例
	rect := Rectangle{Width: 10, Height: 5}
	fmt.Println("矩形信息:")
	fmt.Printf("宽: %.2f, 高: %.2f\n", rect.Width, rect.Height)
	printShapeInfo(rect)

	// 创建 Circle 实例
	circle := Circle{Radius: 7}
	fmt.Println("圆形信息:")
	fmt.Printf("半径: %.2f\n", circle.Radius)
	printShapeInfo(circle)

	// 使用接口类型的切片
	fmt.Println("使用接口切片:")
	shapes := []Shape{
		Rectangle{Width: 3, Height: 4},
		Circle{Radius: 5},
		Rectangle{Width: 6, Height: 8},
		Circle{Radius: 10},
	}

	for i, shape := range shapes {
		switch shape.(type) {
		case Rectangle:
			fmt.Printf("第 %d 个形状是矩形\n", i+1)
		case Circle:
			fmt.Printf("第 %d 个形状是圆形\n", i+1)
		default:
			fmt.Printf("第 %d 个形状是未知类型\n", i+1)
		}
		fmt.Printf("面积: %.2f, 周长: %.2f\n\n", shape.Area(), shape.Perimeter())
	}

	// 类型断言示例
	fmt.Println("类型断言示例:")
	var shape Shape = Circle{Radius: 3}

	if circle, ok := shape.(Circle); ok {
		fmt.Printf("这是一个圆形，半径: %.2f\n", circle.Radius)
		fmt.Printf("圆面积: %.2f\n", circle.Area())
	}

	// 创建 Employee 实例
	emp := Employee{
		Person: Person{
			Name: "张三",
			Age:  30,
		},
		EmployeeID: "E1001",
		Department: "技术部",
	}

	// 调用 PrintInfo 方法
	fmt.Println("=== 使用 PrintInfo() ===")
	emp.PrintInfo()

	// 测试字段访问
	fmt.Println("\n=== 字段访问测试 ===")
	fmt.Printf("直接访问 Name: %s\n", emp.Name)          // 匿名嵌入的便利
	fmt.Printf("完整路径访问 Name: %s\n", emp.Person.Name) // 也可以完整路径访问
}
