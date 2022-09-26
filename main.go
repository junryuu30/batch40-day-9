package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()

	//root for public
	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/formAddProject", formAddProject).Methods("GET")
	route.HandleFunc("/projectDetail/{index}", projectDetail).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")
	// route.HandleFunc("/edit-project/{index}", editProject).Methods("GET")

	fmt.Println("server running in port 9000")
	http.ListenAndServe("localhost:9000", route)
}

// func helloWorld(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello World jihan hallo woy ayo pasti bisa"))
// }

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message:" + err.Error()))
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, description, technologies FROM table_project2")
	fmt.Println(data)

	var result []Project
	for data.Next() {
		var each = Project{}

		var err = data.Scan(&each.ID, &each.Title, &each.Description)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	resData := map[string]interface{}{
		"Project": result,
	}

	fmt.Println(result)

	tmpl.Execute(w, resData)

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("message:" + err.Error()))
	}

	tmpl.Execute(w, nil)

}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/addProject.html")

	if err != nil {
		w.Write([]byte("message:" + err.Error()))
	}

	tmpl.Execute(w, nil)

	// http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

type Project struct {
	ID           int
	Title        string
	Description  string
	Technologies string
	NodeJs       string
	Python       string
	ReactJs      string
	Golang       string
	StartDate    int
	EndDate      int
}

var dataProject = []Project{
	{
		Title:        "Hallo Title",
		Description:  "Ini deskripsinya",
		Technologies: "node-js",
		NodeJs:       "node-js",
		ReactJs:      "react",
		Golang:       "golang",
	},
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(r.PostForm.Get("inputName"))
	// fmt.Println(r.PostForm.Get("description"))

	title := r.PostForm.Get("inputName")
	description := r.PostForm.Get("description")
	// var startDate = r.PostForm.Get("startDate")
	// var endDate = r.PostForm.Get("endDate")

	nodeJs := r.PostForm.Get("nodeJs")
	python := r.PostForm.Get("python")
	reactJs := r.PostForm.Get("react")
	golang := r.PostForm.Get("golang")

	newProject := Project{
		Title:       title,
		Description: description,
		NodeJs:      nodeJs,
		Python:      python,
		ReactJs:     reactJs,
		Golang:      golang,
	}

	dataProject = append(dataProject, newProject)
	fmt.Println(dataProject)

	http.Redirect(w, r, "/formAddProject", http.StatusMovedPermanently)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	dataProject = append(dataProject[:index], dataProject[index+1:]...)
	fmt.Println(dataProject)

	http.Redirect(w, r, "/", http.StatusFound)
}

// func editProject(w http.ResponseWriter, r *http.Request) {

// 	index, _ := strconv.Atoi(mux.Vars(r)["index"])
// 	fmt.Println(index)

// }

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/projectDetail.html")

	if err != nil {
		w.Write([]byte("message:" + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if i == index {
			ProjectDetail = Project{
				Title:       data.Title,
				Description: data.Description,
			}
		}
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}
	fmt.Println(data)
	// data := map[string]interface{}{
	// 	"Title":   "Hello Title",
	// 	"Content": "Hello Content",
	// 	"Id":      index,
	// }

	tmpl.Execute(w, data)

}
