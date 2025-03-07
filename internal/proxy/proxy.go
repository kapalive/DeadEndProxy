package proxy

import (
	"io"
	"log"
	"net/http"
)

// DynamicRouter перенаправляет запросы на нужный бэкенд
func DynamicRouter(backend string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📡 Проксируем запрос: %s -> %s", r.URL.Path, backend)

		req, err := http.NewRequest(r.Method, backend+r.URL.Path, r.Body)
		if err != nil {
			http.Error(w, "Ошибка запроса", http.StatusInternalServerError)
			return
		}

		// Копируем заголовки
		req.Header = r.Header

		// Отправляем запрос
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Ошибка соединения с бэкендом", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Копируем статус-код и заголовки
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
