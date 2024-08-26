package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Admin struct {
	Name     string
	Password string
}
type Ex struct {
	XValues   []string `json:"x"`
	YValues   []int    `json:"y"`
	BarColors []string `json:"color"`
}
type Users struct {
	Name    string
	Surname string
	Bench   int
	Squat   int
	Dead    int
	Pull    int
	Ton     int
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
	rtr.HandleFunc("/insert", insert).Methods("POST")
	rtr.HandleFunc("/update", update).Methods("POST")

	//обработчики силовых
	rtr.HandleFunc("/squat", squat).Methods("GET")
	rtr.HandleFunc("/bench", bench).Methods("GET")
	rtr.HandleFunc("/dead", dead).Methods("GET")
	rtr.HandleFunc("/pull", pull).Methods("GET")
	rtr.HandleFunc("/ton", ton).Methods("GET")

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

func insert(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	surname := r.FormValue("usersurname")
	bench := r.FormValue("bench")
	dead := r.FormValue("dead")
	squat := r.FormValue("squat")
	pull := r.FormValue("pull")
	ton := 0
	dataForSum := []string{bench, squat, dead}
	for _, v := range dataForSum {
		value, err := strconv.Atoi(v)
		if err != nil {
			log.Println("Оштбка перевода в int")
		} else {
			ton += value
		}

	}

	fmt.Println(name, surname, bench, dead, squat, pull, ton)
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()

	result, err := db.Exec(context.Background(), "INSERT INTO users (Name,Surname,Bench,Squat,Dead,Pull,Ton) VALUES ($1,$2,$3,$4,$5,$6,$7)", name, surname, bench, squat, dead, pull, ton)
	fmt.Println(result)

	type Tst struct {
		Name    string
		Surname string
		Bench   int
		Squat   int
		Dead    int
		Pull    int
		Ton     int
	}

	tst, _ := db.Query(context.Background(), "SELECT Name,Surname,Bench,Squat,Dead,Pull,Ton FROM users")
	for tst.Next() {
		var t Tst
		tst.Scan(&t.Name, &t.Surname, &t.Bench, &t.Squat, &t.Dead, &t.Pull, &t.Ton)
		fmt.Println(t)
	}

}

func update(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	surname := r.FormValue("usersurname")

	ex := r.FormValue("update")

	fmt.Println(name, surname, ex)
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()

	type Up struct {
		ex string
	}
	var u Up
	result := db.QueryRow(context.Background(), "SELECT $1 FROM users WHERE Name = $2 AND Surname = $3", ex, name, surname).Scan(&u.ex)
	fmt.Println(u, result)

	// w.Write([]byte("на доработке,сделать нормальный выбор"))

}

// функция для отправки json на присед
func squat(w http.ResponseWriter, r *http.Request) {
	//подключение к бд и парсинг оттуда имен и результаты, цвет статичный , взависимости от упражнений
	colorForGraf := []string{"rgba(255, 99, 132, 0.5)",
		"rgba(255, 159, 64, 0.5)",
		"rgba(255, 205, 86, 0.5)",
		"rgba(75, 192, 192, 0.5)",
		"rgba(54, 162, 235, 0.5)",
		"rgba(153, 102, 255, 0.5)",
		"rgba(201, 203, 207, 0.5)"}
	names := []string{}
	count := []int{}
	colors := []string{}
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()
	result, _ := db.Query(context.Background(), "SELECT Name,Surname,Squat FROM users")
	for result.Next() {
		var t Users
		result.Scan(&t.Name, &t.Surname, &t.Squat)
		if t.Squat == 0 {
			continue
		} else {
			names = append(names, t.Name+" "+t.Surname)
			count = append(count, t.Squat)
		}

	}

	for i := 0; i < len(names); i++ {
		j := i
		if j == 20 {
			j = 0
		}
		colors = append(colors, colorForGraf[j])

	}

	s := Ex{
		XValues:   names,
		YValues:   count,
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
	colorForGraf := []string{"rgba(255, 99, 132, 0.5)",
		"rgba(255, 159, 64, 0.5)",
		"rgba(255, 205, 86, 0.5)",
		"rgba(75, 192, 192, 0.5)",
		"rgba(54, 162, 235, 0.5)",
		"rgba(153, 102, 255, 0.5)",
		"rgba(201, 203, 207, 0.5)"}
	//подключение к бд и парсинг оттуда имен и результаты, цвет статичный , взависимости от упражнений
	names := []string{}
	count := []int{}
	colors := []string{}
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()
	result, _ := db.Query(context.Background(), "SELECT Name,Surname,Bench FROM users")
	for result.Next() {
		var t Users

		result.Scan(&t.Name, &t.Surname, &t.Bench)
		if t.Bench == 0 {
			continue
		} else {
			names = append(names, t.Name+" "+t.Surname)
			count = append(count, t.Bench)
		}

	}

	for i := 0; i < len(names); i++ {
		j := i
		if j == 20 {
			j = 0
		}
		colors = append(colors, colorForGraf[j])

	}

	s := Ex{
		XValues:   names,
		YValues:   count,
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

func dead(w http.ResponseWriter, r *http.Request) {
	colorForGraf := []string{"rgba(255, 99, 132, 0.5)",
		"rgba(255, 159, 64, 0.5)",
		"rgba(255, 205, 86, 0.5)",
		"rgba(75, 192, 192, 0.5)",
		"rgba(54, 162, 235, 0.5)",
		"rgba(153, 102, 255, 0.5)",
		"rgba(201, 203, 207, 0.5)"}
	//подключение к бд и парсинг оттуда имен и результаты, цвет статичный , взависимости от упражнений
	names := []string{}
	count := []int{}
	colors := []string{}
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()
	result, _ := db.Query(context.Background(), "SELECT Name,Surname,Dead FROM users")
	for result.Next() {
		var t Users
		result.Scan(&t.Name, &t.Surname, &t.Dead)
		if t.Dead == 0 {
			continue
		} else {
			names = append(names, t.Name+" "+t.Surname)
			count = append(count, t.Dead)
		}

	}

	for i := 0; i < len(names); i++ {
		j := i
		if j == 20 {
			j = 0
		}
		colors = append(colors, colorForGraf[j])

	}

	s := Ex{
		XValues:   names,
		YValues:   count,
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

func pull(w http.ResponseWriter, r *http.Request) {
	colorForGraf := []string{"rgba(255, 99, 132, 0.5)",
		"rgba(255, 159, 64, 0.5)",
		"rgba(255, 205, 86, 0.5)",
		"rgba(75, 192, 192, 0.5)",
		"rgba(54, 162, 235, 0.5)",
		"rgba(153, 102, 255, 0.5)",
		"rgba(201, 203, 207, 0.5)"}
	//подключение к бд и парсинг оттуда имен и результаты, цвет статичный , взависимости от упражнений
	names := []string{}
	count := []int{}
	colors := []string{}
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()
	result, _ := db.Query(context.Background(), "SELECT Name,Surname,Pull FROM users")
	for result.Next() {
		var t Users
		result.Scan(&t.Name, &t.Surname, &t.Pull)
		if t.Pull == 0 {
			continue
		} else {
			names = append(names, t.Name+" "+t.Surname)
			count = append(count, t.Pull)
		}

	}

	for i := 0; i < len(names); i++ {
		j := i
		if j == 20 {
			j = 0
		}
		colors = append(colors, colorForGraf[j])

	}

	s := Ex{
		XValues:   names,
		YValues:   count,
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

func ton(w http.ResponseWriter, r *http.Request) {
	colorForGraf := []string{"rgba(255, 99, 132, 0.5)",
		"rgba(255, 159, 64, 0.5)",
		"rgba(255, 205, 86, 0.5)",
		"rgba(75, 192, 192, 0.5)",
		"rgba(54, 162, 235, 0.5)",
		"rgba(153, 102, 255, 0.5)",
		"rgba(201, 203, 207, 0.5)"}
	//подключение к бд и парсинг оттуда имен и результаты, цвет статичный , взависимости от упражнений
	names := []string{}
	count := []int{}
	colors := []string{}
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()
	result, _ := db.Query(context.Background(), "SELECT Name,Surname,Ton FROM users")
	for result.Next() {
		var t Users
		result.Scan(&t.Name, &t.Surname, &t.Ton)
		names = append(names, t.Name+" "+t.Surname)
		count = append(count, t.Ton)
	}

	for i := 0; i < len(names); i++ {
		j := i
		if j == 20 {
			j = 0
		}
		colors = append(colors, colorForGraf[j])

	}

	s := Ex{
		XValues:   names,
		YValues:   count,
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
