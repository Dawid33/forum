package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/lib/pq"
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
	category string
	title string
	content string
}
type Comment struct {
	commentId string
	threadId int
	parentId int
	kidsId []int
	userId string
	content string
}

func getThreads(db *sql.DB, whereClause string) ([]Thread, error) {
	query := "SELECT * FROM forum.posts " + whereClause + ";"

	rows, err := db.Query(query)
	if err != nil {
		return []Thread{}, nil
	}
	var output []Thread
	for rows.Next() {
		var postid uint64
		var userid string
		var category string
		var title string
		var content string
		err = rows.Scan(&postid, &userid, &category, &title, &content)
		if err != nil {
			return []Thread{}, nil
		}
		output = append(output, Thread{
			threadId: postid,
			userId:   userid,
			category: category,
			title:    title,
			content:  content,
		})
	}
	return output, nil
}

func getComments(db *sql.DB, whereClause string) ([]Comment, error) {
	query := "SELECT * FROM forum.comments " + whereClause + ";"
	rows, err := db.Query(query)
	if err != nil {
		return []Comment{}, err
	}
	var output []Comment
	for rows.Next() {
		var commentId string
		var threadId int
		var parentId int
		var kidsId []int
		var userId string
		var content string
		err = rows.Scan(&commentId, &threadId, &parentId, pq.Array(&kidsId), &userId, &content)
		if err != nil {
			return []Comment{}, err
		}
		output = append(output, Comment {
			commentId: commentId,
			threadId: threadId,
			parentId: parentId,
			kidsId: kidsId,
			userId: userId,
			content: content,
		})
	}
	return output, nil
}

func connectToDB() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	fmt.Println(dbHost)
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, port, user, password, dbname)
	for {
		conn, err := sql.Open("postgres", psqlconn)
		if err != nil || conn.Ping() != nil{
			log.Println("Cannot connect to database. Trying again...")
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