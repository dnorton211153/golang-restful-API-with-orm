/**
 * @author Norton 2022
 * This is a simple example of a RESTful API,
 * using the gorm ORM library and the gorilla/mux router.
 *
 * Using Repository Interface per Model (best DB migration practice)
 */
package main

import (
	"flag"
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
var userHandler *controllers.BaseHandler

// Track user sessions in memory
var sessions = map[string]controllers.Session{}

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
	userHandler = controllers.NewHandler(userRepository, sessions)

	var dir string
	flag.StringVar(&dir, "dir", "./static", "the directory to serve files from")
	flag.Parse()

	// create a new router
	router := mux.NewRouter()

	// user routes
	router.HandleFunc("/user", userHandler.Create).Methods("POST")
	router.HandleFunc("/user/{id}", userHandler.Get).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.Update).Methods("PUT")
	router.HandleFunc("/user/{id}", userHandler.Delete).Methods("DELETE")
	router.HandleFunc("/user", validateSession(http.HandlerFunc(userHandler.GetAll))).Methods("GET")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/sso", controllers.SSO).Methods("GET")
	router.HandleFunc("/callback", controllers.Callback).Methods("GET")

	// This will serve files under http://localhost:8000/<filename> in the 'dir' directory.
    router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	// start the server
	log.Fatal(http.ListenAndServe(":8000", router))

}

func validateSession(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value

		// We then get the session from our session map
		userSession, exists := sessions[sessionToken]
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if controllers.Expired(userSession) {
			delete(sessions, sessionToken)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

		// do something after the request
	})
}
