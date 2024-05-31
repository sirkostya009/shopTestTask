package middleware

import (
	"fmt"
	"net/http"
	"shopTestTask/cfg"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if err := r.PathValue("__error"); err != "" {
			fmt.Println("error", err, r.Method, r.URL.Path)
		} else {
			fmt.Println(r.Method, r.URL.Path)
		}
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, pass, ok := r.BasicAuth(); !ok || user != cfg.Username || pass != cfg.Password {
			r.SetPathValue("__error", "Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type middleware func(http.Handler) http.Handler

func Combine(middlewares ...middleware) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		for _, m := range middlewares {
			h = m(h)
		}
		return h
	}
}
