# GCP-AlloyDB-Connection-Golang

这是一个用于在Go语言中连接Google Cloud Platform (GCP) AlloyDB的库，它封装了使用GORM进行数据库操作的功能。

[English Documentation](README.md)

## 功能特点

- 简单易用的API接口
- 支持环境变量配置
- 集成GORM ORM功能
- 完善的错误处理机制
- 详细的中文注释和文档
- 灵活的连接池配置

## 前置条件

1. 已创建GCP AlloyDB实例
2. 已获取服务账号密钥文件（JSON格式）
3. Go 1.16或更高版本

## 安装

```bash
go get github.com/PolyhedraZK/GCP-AlloyDB-Connection-Golang@v1.0.0
```

## 环境变量配置

### 必需的环境变量

| 环境变量 | 说明 |
|---------|------|
| DB_HOST | AlloyDB实例URI |
| DB_USER | 数据库用户名 |
| DB_PASS | 数据库密码 |
| DB_NAME | 数据库名称 |
| DB_CERT_PATH | 服务账号密钥文件路径 |

### 可选的环境变量（连接池配置）

| 环境变量 | 说明 | 默认值 | 说明 |
|---------|------|--------|------|
| DB_MAX_OPEN_CONNS | 最大开放连接数 | 0 | 默认0表示无限制 |
| DB_MAX_IDLE_CONNS | 最大空闲连接数 | 2 | Go标准库默认值 |
| DB_CONN_MAX_LIFETIME | 连接最大生命周期（分钟） | 0 | 默认0表示无限制 |
| DB_CONN_MAX_IDLE_TIME | 空闲连接最大生命周期（分钟） | 0 | 默认0表示无限制 |

## 快速开始

### 1. 基本使用

```go
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

    // 现在可以使用db进行数据库操作
}
```

### 2. 使用GORM进行数据库操作

```go
// 定义模型
type User struct {
    ID   uint   `gorm:"primarykey"`
    Name string `gorm:"type:varchar(100)"`
}

// 自动迁移
if err := db.AutoMigrate(&User{}); err != nil {
    log.Fatal(err)
}

// 创建记录
user := User{Name: "测试用户"}
db.Create(&user)

// 查询记录
var users []User
db.Find(&users)
```

## 连接池配置说明

### 1. 最大开放连接数 (DB_MAX_OPEN_CONNS)
- 默认值：0（无限制）
- 建议：
  * 小型应用：5-20
  * 中型应用：20-50
  * 大型应用：50-100
  * 根据服务器资源和实际负载调整

### 2. 最大空闲连接数 (DB_MAX_IDLE_CONNS)
- 默认值：2（Go标准库默认值）
- 建议：
  * 设置为最大开放连接数的1/2到1/3
  * 避免设置过大造成资源浪费
  * 避免设置过小导致频繁创建连接

### 3. 连接最大生命周期 (DB_CONN_MAX_LIFETIME)
- 默认值：0（无限制）
- 建议：
  * 生产环境建议设置为30-120分钟
  * 考虑数据库和网络环境的稳定性
  * 较短的生命周期有助于资源回收

### 4. 空闲连接最大生命周期 (DB_CONN_MAX_IDLE_TIME)
- 默认值：0（无限制）
- 建议：
  * 生产环境建议设置为5-30分钟
  * 较短的时间有助于释放不活跃连接
  * 根据应用访问模式调整

## 最佳实践

1. 环境变量管理
   - 使用.env文件或环境变量管理工具
   - 不要在代码中硬编码敏感信息
   - 在生产环境中妥善保管密钥文件

2. 连接池配置
   - 生产环境建议配置所有连接池参数
   - 根据实际负载调整参数
   - 定期监控连接池状态

3. 错误处理
   - 始终检查InitDB()的返回错误
   - 对数据库操作进行适当的错误处理

4. 安全性
   - 确保服务账号具有适当的权限
   - 使用安全的方式管理数据库凭证
   - 在生产环境中使用SSL连接

## 常见问题

1. 连接失败
   - 检查环境变量是否正确设置
   - 确认服务账号密钥文件路径正确
   - 验证网络连接和防火墙设置

2. 性能问题
   - 检查连接池配置是否合理
   - 监控连接使用情况
   - 考虑调整最大连接数

3. 权限问题
   - 检查服务账号权限配置
   - 确认数据库用户权限设置

## 依赖

- cloud.google.com/go/alloydbconn
- gorm.io/gorm
- gorm.io/driver/postgres

## 版本管理

当前版本：v1.0.0

版本号遵循[语义化版本 2.0.0](https://semver.org/lang/zh-CN/)规范：
- 主版本号：不兼容的API修改
- 次版本号：向下兼容的功能性新增
- 修订号：向下兼容的问题修正

### 版本检查
```go
import "github.com/PolyhedraZK/GCP-AlloyDB-Connection-Golang/connector"

version := connector.GetVersion()
fmt.Printf("当前版本: %s\n", version)
```

## 许可证

MIT License
