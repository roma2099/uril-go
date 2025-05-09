package database

import (
	"fmt"
	"strconv"

	"github.com/roma2099/uril-go/internal/config"
	"github.com/roma2099/uril-go/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

)
var DB *gorm.DB
func ConnectDB(){
	var err error
	p:=config.Config("DB_PORT")
	//if any problem add port insted of _
	_, err=strconv.ParseUint(p,10,32)
	if err!=nil{
		panic("failed to parse database port")
	}
	
	/*dsn :=fmt.Sprint(
		"host=db port=%d user=%s password=%s dbname=%s sslmode=disable",
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)*/
	fmt.Println(config.Config("DB_NAME"))

	DB, err= gorm.Open(sqlite.Open(config.Config("DB_NAME")), &gorm.Config{})
	if err != nil{
		panic("failed to connect to database")

	}
	fmt.Println ("Connection Open to data base")
	DB.AutoMigrate(&model.Product{},&model.User{})
	fmt.Println("Database Migrated")

}