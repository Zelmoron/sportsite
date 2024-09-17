package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"site/settings"
	"site/structs"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"
)

func pull(w http.ResponseWriter, r *http.Request) {

	//подключение к бд и парсинг оттуда имен и результаты, цвет статичный , взависимости от упражнений
	names := []string{}
	count := []int{}

	mapUser := make(map[string]int)
	db, err := pgxpool.Connect(context.Background(), "postgres://igor:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
	}
	defer db.Close()
	result, _ := db.Query(context.Background(), "SELECT Name,Surname,Pull FROM users")
	for result.Next() {
		var t structs.Users
		result.Scan(&t.Name, &t.Surname, &t.Pull)
		if t.Pull == 0 {
			continue
		} else {
			mapUser[t.Name+" "+t.Surname] = t.Pull

		}

	}
	type key_value struct {
		Key   string
		Value int
	}

	var sorted_struct []key_value

	for key, value := range mapUser {
		sorted_struct = append(sorted_struct, key_value{key, value})
	}

	sort.Slice(sorted_struct, func(i, j int) bool {
		return sorted_struct[i].Value < sorted_struct[j].Value
	})

	for _, key_value := range sorted_struct {
		names = append(names, key_value.Key)
		count = append(count, key_value.Value)
	}

	colors := settings.Color(names)

	s := Ex{
		XValues:   names,
		YValues:   count,
		BarColors: colors,
	}
	// Сериализация структуры в строку
	jsonBytes, err := json.Marshal(&s)
	if err != nil {
		f, err := os.OpenFile("text.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString("Ошибка выполнения Marshal 463\n"); err != nil {
			log.Println(err)
		}

	}

	w.Write(jsonBytes)
}
