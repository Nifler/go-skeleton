package main

import (
	"log"
	"net/http"

	"github.com/Nifler/go-skeleton/app"
)

func main() {
	a := app.New(":8080")

	// "/{$}" означає точно корінь. Без {$} шаблон "/" ловив би всі шляхи.
	a.Router.Get("/{$}", hello)

	// Run повертає помилку, а що з нею робити, вирішує main: тут завершуємо
	// програму з повідомленням.
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
