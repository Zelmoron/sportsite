package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Admin struct {
	Name     string
	Password string
}
type Squat struct {
	XValues   []string `json:"x"`
	YValues   []int    `json:"y"`
	BarColors []string `json:"color"`
}

func main() {

	handlers()

}

func handlers() {
	//создал роутер
	rtr := mux.NewRouter()

	//подключил свои стили
	http.Handle("/templates/",
		http.StripPrefix("/templates", http.FileServer(http.Dir("./templates/"))))

	rtr.HandleFunc("/", index)
	rtr.HandleFunc("/logo", logo)
	rtr.HandleFunc("/get", get).Methods("POST")

	//обработчики силовых
	rtr.HandleFunc("/squat", squat).Methods("GET")
	rtr.HandleFunc("/bench", bench).Methods("GET")

	//Все адреса будут обрабатываться через rtr
	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)

}

// функция для вывода главной странички
func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		log.Println("Ошибка обработки html")
	}
	tmpl.ExecuteTemplate(w, "index", nil)

}

// обработчик для перехода на стрницу авторизации админа
func logo(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Println("Ошибка обработки html")
	}
	tmpl.ExecuteTemplate(w, "login", nil)
}

func get(w http.ResponseWriter, r *http.Request) {
	var a Admin
	name := r.FormValue("username")
	password := r.FormValue("userpassword")
	fmt.Println(name, password)
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()

	err = db.QueryRow(context.Background(), "SELECT Name,password FROM admin WHERE Name = $1 AND password = $2", name, password).Scan(&a.Name, &a.Password)
	if err != nil {

		type Check struct {
			Text string
		}
		c := Check{
			Text: "Неправильно введен логин или пароль",
		}
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println("Ошибка обработки html в get")
		}
		tmpl.ExecuteTemplate(w, "login", c)

	} else {
		tmpl, err := template.ParseFiles("templates/admin.html")
		if err != nil {
			log.Println("Ошибка обработки html в get")
		}
		tmpl.ExecuteTemplate(w, "admin", nil)

	}

}

// функция для отправки json на присед
func squat(w http.ResponseWriter, r *http.Request) {
	//подключение к бд и парсинг оттуда имен и результаты, цвет статичный , взависимости от упражнений
	names := []string{"Igor", "Nikita", "Gosha", "Rusia", "FFF"}
	colors := []string{}
	for i := 0; i < len(names); i++ {
		colors = append(colors, "red")

	}

	s := Squat{
		XValues:   names,
		YValues:   []int{55, 49, 44, 24, 15},
		BarColors: colors,
	}
	// Сериализация структуры в строку
	jsonBytes, err := json.Marshal(&s)
	if err != nil {
		// Используем Fatal только для примера,
		// нельзя использовать в реальных приложениях
		log.Fatalln("marshal ", err.Error())
	}

	w.Write(jsonBytes)
}

// функция для отправки json на жим
func bench(w http.ResponseWriter, r *http.Request) {
	names := []string{"Igor", "Nikita", "Gosha", "Rusia", "FFF"}
	colors := []string{}
	for i := 0; i < len(names); i++ {
		colors = append(colors, "white")

	}

	s := Squat{
		XValues:   names,
		YValues:   []int{55, 49, 30, 24, 40},
		BarColors: colors,
	}
	// Сериализация структуры в строку
	jsonBytes, err := json.Marshal(&s)
	if err != nil {
		// Используем Fatal только для примера,
		// нельзя использовать в реальных приложениях
		log.Fatalln("marshal ", err.Error())
	}

	w.Write(jsonBytes)
}
