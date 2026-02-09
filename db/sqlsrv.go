package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

func GetsSQLsrvDB() (*sql.DB, error) {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		os.Getenv("server"), os.Getenv("user"), os.Getenv("password"), os.Getenv("port"), os.Getenv("database"))
	conn, errCon := sql.Open("mssql", connString)
	if errCon != nil {
		log.Fatal("Open connection failed:", errCon.Error())
	}
	fmt.Println("connected")
	return conn, nil

}

func GetsSQLsrvDB2() (*sql.DB, error) {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		os.Getenv("server2"), os.Getenv("user2"), os.Getenv("password2"), os.Getenv("port2"), os.Getenv("database2"))
	conn, errCon := sql.Open("mssql", connString)
	if errCon != nil {
		log.Fatal("Open connection failed:", errCon.Error())
	}
	fmt.Println("connected")
	return conn, nil

}

func GetsSQLsrvDB3() (*sql.DB, error) {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		os.Getenv("server3"), os.Getenv("user3"), os.Getenv("password3"), os.Getenv("port3"), os.Getenv("database3"))
	conn, errCon := sql.Open("mssql", connString)
	if errCon != nil {
		log.Fatal("Open connection failed:", errCon.Error())
	}
	fmt.Println("connected")
	return conn, nil

}
