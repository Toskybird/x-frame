package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/x-frame/x-frame/migrations"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func main() {
	// 读取配置文件
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 构建连接字符串
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.SSLMode,
	)

	// 连接到 postgres 数据库以创建 x-frame 数据库
	log.Println("连接到 PostgreSQL 服务器...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 检查连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	// 创建数据库（如果不存在）
	log.Printf("检查数据库 %s 是否存在...", config.Database.DBName)
	// PostgreSQL 不支持 CREATE DATABASE IF NOT EXISTS，需要使用不同的方法
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", config.Database.DBName).Scan(&exists)
	if err != nil {
		log.Fatalf("检查数据库是否存在失败: %v", err)
	}

	if !exists {
		log.Printf("数据库 %s 不存在，正在创建...", config.Database.DBName)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Database.DBName))
		if err != nil {
			log.Fatalf("创建数据库失败: %v", err)
		}
		log.Printf("数据库 %s 创建成功", config.Database.DBName)
	} else {
		log.Printf("数据库 %s 已存在", config.Database.DBName)
	}

	// 连接到 x-frame 数据库
	connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	log.Printf("连接到数据库 %s...", config.Database.DBName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("连接到 %s 数据库失败: %v", config.Database.DBName, err)
	}
	defer db.Close()

	// 检查连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	// 获取所有迁移文件
	log.Println("获取迁移文件...")
	sqlFiles, err := getMigrationFiles()
	if err != nil {
		log.Fatalf("获取迁移文件失败: %v", err)
	}

	// 执行迁移
	log.Println("开始执行数据库迁移...")
	for _, file := range sqlFiles {
		log.Printf("执行迁移文件: %s", file)
		sqlContent, err := migrations.FS.ReadFile(file)
		if err != nil {
			log.Printf("读取迁移文件 %s 失败: %v", file, err)
			continue
		}

		// 执行 SQL
		_, err = db.ExecContext(context.Background(), string(sqlContent))
		if err != nil {
			log.Printf("执行迁移文件 %s 失败: %v", file, err)
			continue
		}

		log.Printf("迁移文件 %s 执行成功", file)
	}

	log.Println("数据库初始化完成！")
}

func loadConfig(configPath string) (*Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

func getMigrationFiles() ([]string, error) {
	var files []string

	// 读取嵌入的文件系统
	err := fs.WalkDir(migrations.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .sql 文件
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 按文件名排序（确保按顺序执行）
	sort.Strings(files)

	return files, nil
}
