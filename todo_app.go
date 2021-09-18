package main

import (
	"html/template"
	"net/http"
	// "fmt"
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func fetchTodos(writer http.ResponseWriter, request *http.Request) {
	todoFilename := "todos.txt"
	file, err := os.Open(todoFilename)
	defer file.Close()
	todos := make(map[int]string)
	if err == nil {
		// var todos map[int]string
		scanner := bufio.NewScanner(file)
		ind := 0
		for scanner.Scan() {
			todos[ind] = scanner.Text()
			// todos = append(todos, scanner.Text())
			ind += 1
		}
	}

	t, _ := template.ParseFiles("templates/todo_template.html")
	t.Execute(writer, todos)
}

func addTodo(writer http.ResponseWriter, request *http.Request) {
	// request.ParseForm()
	newTodo := request.FormValue("item")
	todoFilename := "todos.txt"
	f, err := os.OpenFile(todoFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString("\n" + newTodo); err != nil {
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

	todoFilename := "todos.txt"
	file, err := os.Open(todoFilename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	var todos []string
	scanner := bufio.NewScanner(file)
	ind := 0
	for scanner.Scan() {
		if ind != indToDelete {
			todos = append(todos, scanner.Text())
		}
		// todos = append(todos, scanner.Text())
		ind += 1
	}

	file, err = os.OpenFile(todoFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileWriter := bufio.NewWriter(file)
	for _, data := range todos {
		_, _ = fileWriter.WriteString(data + "\n")
	}

	fileWriter.Flush()
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
