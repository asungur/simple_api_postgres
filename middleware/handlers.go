package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"simple_api_postgres/models"
	"strconv"

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// response format
type response struct {
	Id      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_STR_VPS"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// create a todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	insertID := insertTodo(todo)
	res := response{
		Id:      insertID,
		Message: "Todo created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

// return a single todo by its id
func GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	todo, err := getTodo(int64(id))

	if err != nil {
		log.Fatalf("Unable to get todo. %v", err)
	}
	json.NewEncoder(w).Encode(todo)
}

// return all the todos
func GetAllTodo(w http.ResponseWriter, r *http.Request) {
	todos, err := getAllTodos()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}
	json.NewEncoder(w).Encode(todos)
}

// update todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	var todo models.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	updatedRows := updateTodo(int64(id), todo)
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)
	res := response{
		Id:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

// delete todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	deletedRows := deleteTodo(int64(id))
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", deletedRows)
	res := response{
		Id:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one todo
func insertTodo(todo models.Todo) int64 {
	db := createConnection()
	// close the db connection
	defer db.Close()
	sqlStatement := `INSERT INTO todos (title, done) VALUES ($1, $2) RETURNING id`
	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, todo.Title, todo.Done).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one todo
func getTodo(id int64) (models.Todo, error) {
	db := createConnection()
	defer db.Close()
	// create an empty entity
	var todo models.Todo
	sqlStatement := `SELECT * FROM todos WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to todo
	err := row.Scan(&todo.Id, &todo.Title, &todo.Done)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return todo, nil
	case nil:
		return todo, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty todo on error
	return todo, err
}

// get all todos
func getAllTodos() ([]models.Todo, error) {
	db := createConnection()
	defer db.Close()
	var todos []models.Todo
	sqlStatement := `SELECT * FROM todos`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()
	// iterate over the rows
	for rows.Next() {
		var todo models.Todo
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Done)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		todos = append(todos, todo)
	}
	return todos, err
}

// update todo
func updateTodo(id int64, todo models.Todo) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `UPDATE todos SET title=$2, done=$3 WHERE id=$1`
	res, err := db.Exec(sqlStatement, id, todo.Title, todo.Done)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete todo
func deleteTodo(id int64) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM todos WHERE id=$1`
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}
