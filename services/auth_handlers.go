package services

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

func SetDB(dbObject *sql.DB){
	DB = dbObject
}

func AuthHandlers(){
	fmt.Print("HIß")
}