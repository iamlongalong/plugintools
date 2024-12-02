package server

import (
	"log"
	"net/http"
	"time"

	"gay/plugintools/internal/config"
)

// Logger 日志中间件
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 包装ResponseWriter以捕获状态码
		wrapper := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next(wrapper, r)

		// 记录请求信息
		log.Printf(
			"%s %s %s %d %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrapper.status,
			time.Since(start),
		)
	}
}

// Auth 认证中间件
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Get()
		if !cfg.Security.EnableAuth {
			next(w, r)
			return
		}

		// 从请求头获取API密钥
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "API key is required", http.StatusUnauthorized)
			return
		}

		// 验证API密钥
		valid := false
		for _, key := range cfg.Security.APIKeys {
			if apiKey == key {
				valid = true
				break
			}
		}

		if !valid {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// responseWriter 包装http.ResponseWriter以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Chain 链接多个中间件
func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}
