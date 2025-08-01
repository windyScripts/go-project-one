package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb(dbname string) (*sql.DB, error){

	connectionString := "db_user:your_password@tcp(127.0.0.1:3306)/" + dbname
	db, err := sql.Open("mysql", connectionString)
	if err != nil{
		//panic(err)
		return nil, err
	}
	fmt.Println("Connected to MariaDB")
	return db, nil
}

//create a limited access user that can only access the single db.

