package main

import (
	"embed"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

//go:embed sql/checkIfSchemaExists.sql
//go:embed sql/forumCreateSchema.sql
//go:embed sql/otherCreateSchema.sql
var f embed.FS

func main() {
	// Connect to database
	db, err := ConnectToDB()
	if err != nil {
		fmt.Println(err)
	}

	for {
		if err := db.Ping(); err != nil {
			log.Println("Cannot connect to database, trying again.")
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}

	var schemas = []string{
		"forum",
	}
	// Reset database at restart
	CheckError(DropAllSchemas(db, schemas))
	CheckError(CreateMissingSchemas(db, schemas))
	db.Close()
	CheckError(err)


	startHttpServer()
}

func GetSQLFile(name string) string {
	data, _ := f.ReadFile(fmt.Sprintf("sql/%s.sql",name))
	return string(data)
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}