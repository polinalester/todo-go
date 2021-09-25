package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Todo struct {
	Id   int
	Name string
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
}

func getTodos(limit int) (todos []Todo, err error) {
	rows, err := Db.Query("select id, name from todos limit $1", limit)
	if err != nil {
		return
	}
	for rows.Next() {
		todo := Todo{}
		err = rows.Scan(&todo.Id, &todo.Name)
		if err != nil {
			return
		}
		todos = append(todos, todo)
	}
	rows.Close()
	return

}

func (todo *Todo) addTodo() (err error) {
	stm := "insert into todos (name) values ($1) returning id"
	s, err := Db.Prepare(stm)
	if err != nil {
		return
	}
	defer s.Close()
	err = s.QueryRow(todo.Name).Scan(&todo.Id)
	return
}

func (todo *Todo) Delete() (err error) {
	_, err = Db.Exec("delete from todos where id = $1", todo.Id)
	return
}

func fetchTodos(writer http.ResponseWriter, request *http.Request) {

	todos := make(map[int]string)
	results, err := getTodos(10)
	for _, t := range results {
		todos[t.Id] = t.Name
	}
	if err != nil {
		log.Fatal(err)
	}

	t, _ := template.ParseFiles("templates/todo_template.html")
	t.Execute(writer, todos)
}

func addTodo(writer http.ResponseWriter, request *http.Request) {
	// request.ParseForm()
	newItem := request.FormValue("item")
	newTodo := Todo{
		Name: newItem,
	}

	err := newTodo.addTodo()
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}

func deleteTodo(writer http.ResponseWriter, request *http.Request) {
	requestURL := request.URL.Path
	urlSplit := strings.Split(requestURL, "/")
	// indToDelete := urlSplit[len(urlSplit)-1]
	indToDelete, err := strconv.Atoi(urlSplit[len(urlSplit)-1])
	if err != nil {
		log.Fatal(err)
	}

	newTodo := Todo{
		Id: indToDelete,
	}

	_ = newTodo.Delete()
	http.Redirect(writer, request, "/", http.StatusFound)

}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8081",
	}
	http.HandleFunc("/", fetchTodos)
	http.HandleFunc("/add/", addTodo)
	http.HandleFunc("/delete/", deleteTodo)
	server.ListenAndServe()
}
