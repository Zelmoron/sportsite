package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	cache map[string]Admin
}
type Admin struct {
	Id       int
	Name     string
	Password string
}
type Ex struct {
	XValues   []string `json:"x"`
	YValues   []int    `json:"y"`
	BarColors []string `json:"color"`
}

// функция для вывода главной странички
func index(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		log.Println("Ошибка обработки html")
		return
	}
	tmpl.ExecuteTemplate(w, "index", nil)

}
func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
		return value, errors.New("you are trying to read empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return value, err
	}
	str := cookie.Value
	value, _ = url.QueryUnescape(str)
	return value, err
}

// обработчик для перехода на стрницу авторизации админа
func (c *App) admin(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie("token", r)
	if err != nil {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println("Ошибка обработки html")
			return
		}
		tmpl.ExecuteTemplate(w, "login", nil)
		fmt.Println("Кука устарела")
		return
	}
	fmt.Println(token)

	if _, ok := c.cache[token]; !ok {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println("Ошибка обработки html")
			return
		}
		tmpl.ExecuteTemplate(w, "login", nil)
		fmt.Println("Токена нет в кеше")
		return
	}
	fmt.Println(c.cache[token].Id)
	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		log.Println("Ошибка обработки html в get")
		return
	}
	tmpl.ExecuteTemplate(w, "admin", nil)

}

func (c *App) get(w http.ResponseWriter, r *http.Request) {
	var a Admin
	name := r.FormValue("username")
	password := r.FormValue("userpassword")

	db, err := pgxpool.Connect(context.Background(), "postgres://igor:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
		return
	}
	defer db.Close()

	err = db.QueryRow(context.Background(), "SELECT Id,Name,password FROM admin WHERE Name = $1 AND password = $2", name, password).Scan(&a.Id, &a.Name, &a.Password)
	if err != nil {

		type Check struct {
			Text string
		}
		e := Check{
			Text: "Неправильно введен логин или пароль",
		}
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println("Ошибка обработки html в get")
		}
		tmpl.ExecuteTemplate(w, "login", e)

	} else {
		//логин и пароль совпадают, поэтому генерируем токен, пишем его в кеш и в куки

		time64 := time.Now().Unix()
		timeInt := string(time64)
		token := name + password + timeInt

		hashToken := md5.Sum([]byte(token))
		hashedToken := hex.EncodeToString(hashToken[:])
		c.cache = make(map[string]Admin)
		c.cache[hashedToken] = a

		livingTime := 1 * time.Minute
		expiration := time.Now().Add(livingTime)
		//кука будет жить 1 час
		cookie := http.Cookie{Name: "token", Value: url.QueryEscape(hashedToken), Expires: expiration}
		http.SetCookie(w, &cookie)
		tmpl, err := template.ParseFiles("templates/admin.html")
		if err != nil {
			log.Println("Ошибка обработки html в get")
			return
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
			return
		} else {
			ton += value
		}

	}

	db, err := pgxpool.Connect(context.Background(), "postgres://igor:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
		return
	}
	defer db.Close()

	_, err = db.Exec(context.Background(), "INSERT INTO users (Name,Surname,Bench,Squat,Dead,Pull,Ton) VALUES ($1,$2,$3,$4,$5,$6,$7)", name, surname, bench, squat, dead, pull, ton)
	if err != nil {
		f, err := os.OpenFile("text.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString("Ошибка выполнения insert строка 120\n"); err != nil {
			log.Println(err)
		}

	}

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

	}
	http.Redirect(w, r, "http://127.0.0.1:8080", http.StatusSeeOther)

}

func update(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	surname := r.FormValue("usersurname")

	ex := r.FormValue("update")

	new_ex := r.FormValue("ex")

	newIntEx, _ := strconv.Atoi(new_ex)

	db, err := pgxpool.Connect(context.Background(), "postgres://igor:132313Igor@localhost:5432/sportsite")

	if err != nil {
		log.Println("Error with connection")
		return
	}
	defer db.Close()

	type Up struct {
		exercises int
		ton       int
	}
	var u Up

	err = db.QueryRow(context.Background(), fmt.Sprintf("SELECT %s,Ton Bench FROM users WHERE Name = $1 AND Surname = $2", ex), name, surname).Scan(&u.exercises, &u.ton)
	if err != nil {
		f, err := os.OpenFile("text.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString("Ошибка выполнения sql select строка 177\n"); err != nil {
			log.Println(err)
		}

	}

	different := newIntEx - u.exercises
	newTon := different + u.ton

	_, err = db.Exec(context.Background(), fmt.Sprintf("UPDATE users SET %s = %d WHERE Name = $1 AND Surname = $2", ex, newIntEx), name, surname)
	if err != nil {
		f, err := os.OpenFile("text.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString("Ошибка выполнения sql select строка 194\n"); err != nil {
			log.Println(err)
		}

	}
	_, err = db.Exec(context.Background(), fmt.Sprintf("UPDATE users SET Ton = %d WHERE Name = $1 AND Surname = $2", newTon), name, surname)
	if err != nil {
		f, err := os.OpenFile("text.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString("Ошибка выполнения sql select строка 207\n"); err != nil {
			log.Println(err)
		}

	}

	http.Redirect(w, r, "http://127.0.0.1:8080", http.StatusSeeOther)

}
