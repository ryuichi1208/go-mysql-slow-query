package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// データベースに接続
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// タイムアウトを設定
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// クエリ実行
	rows, err := db.QueryContext(ctx, "SELECT sleep(10) FROM test_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 結果の取得
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
