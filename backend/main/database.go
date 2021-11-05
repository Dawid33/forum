package main

import (
	"database/sql"
	_ "embed"
	"fmt"
)

const createSchemaSQL = `
CREATE SCHEMA IF NOT EXISTS forum
CREATE TABLE forum.posts {
	userid	text
	postid 	text PRIMARY
	post	text
}
`

const (
	host     = "database"
	port     = 5432
	user     = "postgres"
	password = "test"
	dbname   = "postgres"
)

func connectToDB() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return sql.Open("postgres", psqlconn)
}

func DropAllSchemas(db* sql.DB, schemas []string) error {
	for _, x := range schemas {
		// TODO: Make this work without Sprintf
		_, err := db.Exec(fmt.Sprintf("drop schema if exists %s cascade;", x))
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateMissingSchemas(db *sql.DB, schemas []string) error {
	exists, err := CheckIfSchemasExists(db, schemas)
	CheckError(err)
	for i, x := range exists {
		if x {
			fmt.Printf("%s : YES\n", schemas[i])
		} else {
			fmt.Printf("Does %s exist? : NO\n", schemas[i])
			fmt.Printf("Creating %s Schema...\n", schemas[i])
			query := GetSQLFile(fmt.Sprintf("%sCreateSchema", schemas[i]))
			_, err := db.Exec(query)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
func checkIfTablesExistInSchema() {

}

func checkDbConnection(db *sql.DB) {
	err := db.Ping()
	CheckError(err)
}