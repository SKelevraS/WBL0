package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"

	"github.com/nats-io/nats.go"
)

type CheckBeforeSending struct {
	OrderUid string `json:"order_uid"`
}

func start(w http.ResponseWriter, r *http.Request) {

	m := template.Must(template.ParseFiles("templates/SendingForm.html"))

	executeErr := m.Execute(w, "form")
	checkErr1(executeErr)

	newMsg := r.FormValue("Json")

	switch newMsg {
	case "":
		break
	default:
		var data CheckBeforeSending
		unmarshalErr := json.Unmarshal([]byte(newMsg), &data)
		if unmarshalErr != nil {
			_, err := fmt.Fprintf(w, "<h1>Data sending error! Please send json information format.</h1>")
			checkErr1(err)
		} else {
			nc, err := nats.Connect(nats.DefaultURL, nats.Name("NatsConn")) // подключение к серверу Nats
			checkErr1(err)

			publishErr := nc.Publish("login", []byte(newMsg)) // Отправка новых данных в Nats
			checkErr1(publishErr)
		}
	}
}

func main() {
	http.HandleFunc("/", start)

	err := http.ListenAndServe(":5000", nil) // Подключение к сайту 5001
	checkErr1(err)

	runtime.Goexit()
}

func checkErr1(err error) {
	if err != nil {
		panic(err)
	}
}
