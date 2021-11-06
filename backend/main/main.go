package main

import (
	"embed"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

//go:embed sql/checkIfSchemaExists.sql
//go:embed sql/forumCreateSchema.sql
//go:embed sql/otherCreateSchema.sql
var f embed.FS

func main() {
	// Connect to database
	db := connectToDB()
	var schemas = []string{
		"forum",
	}
	fmt.Println("Deleting all schemas for clean slate...")
	DropAllSchemas(db, schemas)
	CreateMissingSchemas(db, schemas)
	err := db.Close()
	CheckError(err)
	startHttpServer()
}

func GetSQLFile(name string) string {
	data, _ := f.ReadFile(fmt.Sprintf("sql/%s.sql",name))
	return string(data)
}

// CheckError Review and replace this function wherever possible
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PrintError(err error)  {
	if err != nil {
		log.Println(err)
	}
}

func Panic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}