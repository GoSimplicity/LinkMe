package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:v6SxhWHyZC7S@tcp(localhost:33306)/linkme?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	fmt.Println("✅ Successfully connected to MySQL!")

	// 插入数据
	insertSQL := `INSERT INTO users (created_at, updated_at, deleted_at, username, password_hash, deleted, roles) 
				  VALUES (UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), NULL, ?, ?, ?, ?);`

	_, err = db.Exec(insertSQL, "admin2", "$2a$10$wH8sHZTflD5vKj5iHxD5reZ6eYPs1E4/RyTc7HbYJxU6sphGzHl7i", 0, `[{"role":"admin"}]`)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}

	fmt.Println("✅ User 'admin' inserted successfully!")
}
