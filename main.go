package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", RequestHandler)
	http.HandleFunc("/todo", TemplateHandler)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("Failed to initialize server. Error:", err)
		os.Exit(1)
	}
	fmt.Println("Successfully started server")
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request from %s\n", r.URL.Path)
}

var example_template = template.Must(template.ParseFiles("./templates/example.html"))

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	example_template.Execute(w, data)
}
