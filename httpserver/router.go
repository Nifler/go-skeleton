// Package httpserver містить HTTP-сервіс каркаса: роутер, а далі й сервер.
//
// Коментар прямо перед "package" називається doc-коментарем пакета.
// Саме його показують підказки IDE та команда `go doc`.
package httpserver

import "net/http"

// Router — наш власний роутер. Це тонка обгортка над стандартним
// http.ServeMux: ми не винаходимо маршрутизацію, а лише додаємо зручний
// інтерфейс (Get, Post, ...), на який далі нарощуватимемо групи та middleware.
//
// Назва типу з великої літери — отже він ЕКСПОРТОВАНИЙ: його видно з інших
// пакетів (httpserver.Router). Це правило Go: велика літера = публічне,
// мала = приватне для пакета. Ключових слів public/private в Go немає.
type Router struct {
	// Поле з малої літери — приватне. Зовні пакета до mux не дістатися,
	// тож ніхто не зможе обійти наш Router і зареєструвати маршрут напряму.
	mux *http.ServeMux
}

// NewRouter створює порожній роутер.
//
// У Go немає конструкторів як у PHP/Java. Прийнято писати звичайну функцію
// з іменем New<Тип>, яка повертає готовий до роботи об'єкт.
// Повертаємо вказівник (*Router), щоб усі, хто отримав роутер, працювали
// з ОДНИМ і тим самим екземпляром, а не з копіями.
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Handle реєструє handler для комбінації "метод + шлях".
//
// Це базовий метод: Get і Post нижче лише викликають його.
// (r *Router) перед назвою — "приймач" (receiver). Так у Go оголошують метод
// типу. Це аналог $this у PHP, але ім'я ми обираємо самі (за традицією —
// коротка літера від назви типу).
func (r *Router) Handle(method, path string, h http.HandlerFunc) {
	// Починаючи з Go 1.22 ServeMux розуміє шаблони виду "GET /orders/{id}":
	// метод, пробіл, шлях. Тому достатньо склеїти рядок.
	//
	// Якщо маршрут дублюється або конфліктує з іншим, HandleFunc викличе
	// panic. Для помилок конфігурації це навмисно: програма впаде при
	// старті, а не буде мовчки працювати неправильно.
	r.mux.HandleFunc(method+" "+path, h)
}

// Get реєструє маршрут для HTTP-методу GET.
func (r *Router) Get(path string, h http.HandlerFunc) {
	r.Handle(http.MethodGet, path, h)
}

// Post реєструє маршрут для HTTP-методу POST.
//
// http.MethodGet та http.MethodPost — це константи зі стандартної
// бібліотеки ("GET", "POST"). Краще користуватися ними, ніж писати рядки
// вручну: опечатку в константі зловить компілятор, а в рядку ні.
func (r *Router) Post(path string, h http.HandlerFunc) {
	r.Handle(http.MethodPost, path, h)
}

// ServeHTTP робить Router обробником HTTP-запитів.
//
// У Go інтерфейси задовольняються НЕЯВНО. Інтерфейс http.Handler вимагає
// єдиного методу ServeHTTP(w, r). Ми його оголосили, отже Router тепер
// автоматично є http.Handler, і його можна передати в http.ListenAndServe
// чи http.Server. Слова "implements" писати не потрібно.
//
// Тут ми просто передаємо запит далі в mux. Згодом саме в цьому місці з'явиться
// додаткова логіка (middleware).
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
