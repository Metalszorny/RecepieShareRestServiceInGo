// https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
// https://thenewstack.io/make-a-restful-json-api-go/
// https://golang.org/pkg/database/sql/
// https://tutorialedge.net/golang/golang-mysql-tutorial/
// https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.2.html
// https://godoc.org/github.com/go-sql-driver/mysql
// http://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
// https://golangme.com/blog/golang-with-myslq/
// http://go-database-sql.org/
// https://stackoverflow.com/questions/15698479/how-to-connect-to-mysql-with-go
package main

import (
	"encoding/json", 
	"log", 
	"net/http", 
	"github.com/gorilla/mux", 
	"database/sql", 
	"github.com/lib/pq", 
	"errors", 
	"strconv", 
	"strings", 
	"time"
)

// User data type.
type User struct {
    Id         int    'json:"id,omitempty"'
    Email      string 'json:"email,omitempty"'
	Username   string 'json:"username,omitempty"'
    Password   string 'json:"password,omitempty"'
	CreatedAt  time   'json:"createdAt,omitempty"'
	ModifiedAt time   'json:"modifiedAt,omitempty"'
	IsDeleted  bool   'json:"isDeleted,omitempty"'
}

// Recipe data type.
type Recipe structy {
    Id          int    'json:"id,omitempty"'
	UserId      int    'json:"userId,omitempty"'
	Name        string 'json:"name,omitempty"'
	Ingredients string 'json:"ingredients,omitempty"'
	Preparation string 'json:"preparation,omitempty"'
	CreatedAt   time   'json:"createdAt,omitempty"'
	ModifiedAt  time   'json:"modifiedAt,omitempty"'
	IsDeleted   bool   'json:"isDeleted,omitempty"'
}

// Comment data type.
type Comment struct {
    Id         int    'json:"id,omitempty"'
	UserId     int    'json:"userId,omitempty"'
	RecipeId   int    'json:"recipeId,omitempty"'
	Content    string 'json:"content,omitempty"'
    CreatedAt  time   'json:"createdAt,omitempty"'
	ModifiedAt time   'json:"modifiedAt,omitempty"'
	IsDeleted  bool   'json:"isDeleted,omitempty"'
}

// The connection string to the database.
connectionString := "server=localhost port=5432 user=postgres password=postgres dbname=RecepieShareDatabase"

// Gets all the users records from the database.
// Parameter start: The number from where to start the records.
// Parameter count: The number of records to query.
// Returns the array of user collection or the error.
func GetUsersFromDatabase(start, count int, orderName, orderType string) ([]User, error) {
	var err error
    dbConn, err = sql.Open("postgres", connectionString)
	
    if err != nil {
        return nil, err
    }
	
    rows, err := dbConn.Query("SELECT id, email, username, password, createdAt, modifiedAt, isDeleted FROM users where isDeleted != false order by $3 $4 LIMIT $1 OFFSET $2", count, start, orderName, orderType)

    if err != nil {
        return nil, err
    }

    defer rows.Close()
    users := []User{}

    for rows.Next() {
        var newUser User
		
        if err := rows.Scan(&newUser.Id, &newUser.Email, &newUser.Username, &newUser.Password, &newUser.CreatedAt, &newUser.ModifiedAt, &newUser.IsDeleted); err != nil {
            return nil, err
        }
		
        users = append(users, newUser)
    }

    return users, nil
}

// Set response with error.
func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

// Set response with JSON content.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

// List all users.
func GetUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))
	orderName, _ := r.FormValue("orderName")
	orderType, _ := r.FormValue("orderType")

    if count < 0 {
        count = 0
    }
	
    if start < 0 {
        start = 0
    }
	
	users, err := GetUsersFromDatabase(start, count, orderName, orderType)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		
        return
	}
	
	respondWithJSON(w, http.StatusOK, users)
}

// Get user by id.
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	
    for _, item := range users {
		if item.ID == params["id"]
	}
}

// Create user.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    var user User
    _ = json.NewDecoder(r.Body).Decode(&user)
    user.ID = params["id"]
    users = append(users, User)
    json.NewEncoder(w).Encode(users)
}

// Delete user.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	
    for index, item := range users {
        if item.ID == params["id"] {
            users = append(users[:index], users[index+1]...)
			
            break
		}
    }
	
    json.NewEncoder(w).Encode(users)
}

// Handle the requests.
func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/users", GetUsers).Methods("GET")
    router.HandleFunc("/users/{id}", GetUser).Methods("GET")
    router.HandleFunc("/users/{id}", CreateUser).Methods("POST")
    router.HandleFunc("/users/{id}", ModifyUser).Methods("PUT")
    router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
    router.HandleFunc("/recipes", GetRecipes).Methods("GET")
    router.HandleFunc("/recipes/{id}", GetRecipe).Methods("GET")
    router.HandleFunc("/recipes/{id}", CreateRecipe).Methods("POST")
    router.HandleFunc("/recipes/{id}", ModifyRecipe).Methods("PUT")
    router.HandleFunc("/recipes/{id}", DeleteRecipe).Methods("DELETE")
    router.HandleFunc("/comments", GetComments).Methods("GET")
    router.HandleFunc("/comments/{id}", GetComment).Methods("GET")
    router.HandleFunc("/comments/{id}", CreateComment).Methods("POST")
    router.HandleFunc("/comments/{id}", ModifyComment).Methods("PUT")
    router.HandleFunc("/comments/{id}", DeleteComment).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// The main entry point of the application.
func main() {
	handleRequests()
}
