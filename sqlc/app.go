package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/tutorial"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

func run() error {
	ctx := context.Background()

	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/mydb?parseTime=true")
	if err != nil {
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	result, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}

	insertedAuthorID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Println(insertedAuthorID)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthorID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthorID, fetchedAuthor.ID))
	return nil
}

func _select(q string) {
	ctx := context.Background()

	//DB接続
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/mydb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	queries := tutorial.New(db)

	//クエリ実行
	authors, err := queries.GetName(ctx, q)
	fmt.Println(authors, err)
}

func main() {
	_select("Brian%")
	// if err := run(); err != nil {
	// log.Fatal(err)
	// }
}
