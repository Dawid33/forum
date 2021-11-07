package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	port     = 5432
	user     = "postgres"
	password = "test"
	dbname   = "postgres"
)

type Thread struct {
	threadId uint64
	userId string
	title string
	content string
}

func getThreadsWithCategory(category string) ([]Thread, error) {
	db := connectToDB()
	rows, err := db.Query("SELECT * FROM forum.posts WHERE posts.category = $1::text;", category)
	if err != nil {
		return []Thread{}, nil
	}
	var output []Thread
	for rows.Next() {
		var thread Thread
		var postid uint64
		var userid string
		var category string
		var title string
		var content string
		err = rows.Scan(&postid, &userid, &category, &title, &content)
		if err != nil {
			return []Thread{}, nil
		}

		thread.threadId = postid
		thread.userId = userid
		thread.title = title
		thread.content = content

		output = append(output, thread)
	}
	return output, nil
}

func getCommentsInThread(db *sql.DB, threadId string) (string, error) {
	rows, err := db.Query("SELECT userID, content FROM forum.comments WHERE comments.threadID = $1::bigint;", threadId)
	if err != nil {
		return "", err
	}
	var output = ""
	for rows.Next() {
		var userId string
		var content string
		err = rows.Scan(&userId, &content)
		if err != nil {
			return "", err
		}
		output += "test"
	}
	return output, nil
}

func connectToDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	for {
		conn, err := sql.Open("postgres", psqlconn)
		if err != nil {
			log.Println("Cannot connect to database, trying again...")
			time.Sleep(time.Second * 2)
		} else {
			return conn
		}
	}
}

// DropAllSchemas This function must succeed, so it can panic all it wants.
func DropAllSchemas(db* sql.DB, schemas []string) {
	for _, x := range schemas {
		// TODO: Make this work without Sprintf
		_, err := db.Exec(fmt.Sprintf("drop schema if exists %s cascade;", x))
		Panic(err)
	}
}

// CreateMissingSchemas This function must succeed, so it can panic all it wants.
func CreateMissingSchemas(db *sql.DB, schemas []string) {
	exists, err := CheckIfSchemasExists(db, schemas)
	Panic(err)
	for i, x := range exists {
		if x {
			fmt.Printf("%s : YES\n", schemas[i])
		} else {
			fmt.Printf("Does %s exist? : NO\n", schemas[i])
			fmt.Printf("Creating %s Schema...\n", schemas[i])
			query := GetSQLFile(fmt.Sprintf("%sCreateSchema", schemas[i]))
			_, err := db.Exec(query)
			Panic(err)
		}
	}
}

func CheckIfSchemasExists(db *sql.DB, schemas []string) ([]bool, error) {
	data := GetSQLFile("checkIfSchemaExists")

	var hasSchema = make([]bool, len(schemas))

	for i, x := range schemas {
		rows, err := db.Query(data, x)
		if rows != nil {
			for rows.Next() {
				var exists bool
				err = rows.Scan(&exists)
				if err != nil {
					return []bool{false}, err
				}
				hasSchema[i] = exists
			}
		}
	}

	return hasSchema, nil
}