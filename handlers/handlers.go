package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Handlers() {
	//создал роутер
	rtr := mux.NewRouter()
	var c App
	//подключил свои стили
	http.Handle("/templates/",
		http.StripPrefix("/templates", http.FileServer(http.Dir("./templates/"))))

	rtr.HandleFunc("/", c.index)
	rtr.HandleFunc("/logo", logo)
	rtr.HandleFunc("/get", c.get).Methods("POST")
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

}
