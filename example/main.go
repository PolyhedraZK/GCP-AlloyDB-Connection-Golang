package main

import (
	"log"

	"github.com/PolyhedraZK/GCP-AlloyDB-Connection-Golang/connector"
)

func main() {
	// 初始化数据库连接
	if err := connector.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 获取GORM实例
	db := connector.GetDB()

	// 示例：定义一个用户模型
	type User struct {
		ID   uint   `gorm:"primarykey"`
		Name string `gorm:"type:varchar(100)"`
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("迁移表结构失败: %v", err)
	}

	// 创建用户
	user := User{Name: "测试用户"}
	if err := db.Create(&user).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
	}

	// 查询用户
	var users []User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("查询用户失败: %v", err)
	}

	log.Printf("查询到 %d 个用户", len(users))
}
