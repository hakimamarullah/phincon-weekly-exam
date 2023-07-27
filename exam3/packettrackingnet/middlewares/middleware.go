package middlewares

import (
	"log"
	"net/http"
	"packettrackingnet/config"
	"packettrackingnet/dto"
	"packettrackingnet/helpers"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Checking authentication")
		_, err := r.Cookie("uid")
		if err != nil && !config.WhiteListed(r.URL.Path) {
			helpers.ResponseJSON(w, dto.ResponseBody{Message: "Unauthorized", Code: http.StatusUnauthorized})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		defer func() {
			log.Printf("%-7s %s - %d\n", r.Method, r.URL.String(), rw.statusCode)
		}()
		next.ServeHTTP(w, r)
	})
}

func HandlerAdvice(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Handle the error
				log.Println("Internal Server Error:", err)

				helpers.ResponseJSON(w, dto.ResponseBody{Message: "An internal server error has occurred.", Code: http.StatusInternalServerError})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
