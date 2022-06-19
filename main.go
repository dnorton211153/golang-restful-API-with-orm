/**
 * @author Norton 2022
 * This is a simple example of a RESTful API,
 * using the gorm ORM library and the gorilla/mux router.
 *
 * Using Repository Interface per Model (best DB migration practice)
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"

	"restAPI/controllers"
	"restAPI/repositories"
)

var db *gorm.DB

// main function
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	fmt.Println("Starting the application...")

	var dbHost = os.Getenv("DB_HOST")
	var dbPort = os.Getenv("DB_PORT")
	var dbUser = os.Getenv("DB_USER")
	var dbPass = os.Getenv("DB_PASS")
	var dbName = os.Getenv("DB_NAME")

	// gorm.Open takes 3 parameters:
	// 1. dialect, 2. connection string, 3. options
	db, err = gorm.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Create repositories with the database connection
	userRepository := repositories.NewUserRepository(db)

	// Create handlers with the repositories
	userHandler := controllers.NewHandler(userRepository)

	// create a new router
	router := mux.NewRouter()

	// user routes
	router.HandleFunc("/user", userHandler.Create).Methods("POST")
	router.HandleFunc("/user/{id}", userHandler.Get).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.Update).Methods("PUT")
	router.HandleFunc("/user/{id}", userHandler.Delete).Methods("DELETE")
	router.HandleFunc("/user", userHandler.GetAll).Methods("GET")

	// start the server
	log.Fatal(http.ListenAndServe(":8000", router))

}
