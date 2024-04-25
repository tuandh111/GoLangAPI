package main

import (
	"awesomeProject2/cmd/api"
	"awesomeProject2/configs"
	db2 "awesomeProject2/db"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err := db2.NewMysqlStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB successfully connected")
}
