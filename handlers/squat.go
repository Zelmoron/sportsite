package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"
)

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
	mapUser := make(map[string]int)
	db, err := pgxpool.Connect(context.Background(), "postgres://igor:132313Igor@localhost:5432/sportsite")

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
			mapUser[t.Name+" "+t.Surname] = t.Squat

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

	for i := 0; i < len(names); i++ {
		j := i
		if j == 7 {
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
		f, err := os.OpenFile("text.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString("Ошибка выполнения Marshal 272\n"); err != nil {
			log.Println(err)
		}

	}

	w.Write(jsonBytes)
}
