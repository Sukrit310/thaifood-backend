// ไฟล์ mw/http.go: middleware ส่วนกลางสำหรับ CORS บันทึก request และรับมือ panic ด้วย JSON error
// Package mw provides HTTP authentication, CORS and request logging.
package mw

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// จำกัดเวลาที่ middleware ใช้โหลดสมาชิกจากฐานข้อมูลต่อ request
const databaseTimeout = 5 * time.Second

// writeError ส่งรหัสสถานะและ JSON error รูปแบบเดียวกันจาก middleware
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": code, "message": message}})
}

// CORS อนุญาต origin ที่ตรงกับรายการตั้งค่า และตรวจ method/header ของ preflight ก่อนส่งต่อ
func CORS(origins []string, next http.Handler) http.Handler {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin != "" && origin != "*" {
			allowed[origin] = true
		}
	}
	methods := map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true, "OPTIONS": true}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}
		if len(r.Header.Values("Origin")) != 1 || !allowed[origin] {
			writeError(w, http.StatusForbidden, "CORS_FORBIDDEN", "Origin is not allowed")
			return
		}
		if r.Method == http.MethodOptions {
			// ตรวจคำขอ preflight ของ browser ก่อนอนุญาตให้ส่ง method และ headers จริง
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")
			method := r.Header.Get("Access-Control-Request-Method")
			valid := methods[method]
			for _, header := range strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",") {
				header = strings.TrimSpace(header)
				if header != "" && !strings.EqualFold(header, "Authorization") && !strings.EqualFold(header, "Content-Type") {
					valid = false
				}
			}
			if !valid {
				writeError(w, http.StatusForbidden, "CORS_FORBIDDEN", "Preflight request is not allowed")
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		next.ServeHTTP(w, r)
	})
}

// responseRecorder ห่อ ResponseWriter เพื่อเก็บ status ที่ Logging ต้องใช้หลัง handler ทำงานจบ
type responseRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader จดจำและส่งสถานะ HTTP เฉพาะครั้งแรก เพื่อให้ log ตรงกับ response จริง
func (w *responseRecorder) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
		w.ResponseWriter.WriteHeader(status)
	}
}

// Write ตั้งสถานะเริ่มต้นเป็น 200 เมื่อเริ่มเขียน body โดยยังไม่ได้ส่ง header
func (w *responseRecorder) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

// Unwrap คืน ResponseWriter เดิมเพื่อให้เครื่องมือของ net/http เข้าถึงความสามารถของตัวที่ถูกห่อได้
func (w *responseRecorder) Unwrap() http.ResponseWriter { return w.ResponseWriter }

// Logging บันทึก method, path, สถานะ และเวลาที่ใช้ โดยไม่บันทึก query, body หรือ token
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		status := recorder.status
		if status == 0 {
			status = http.StatusOK
		}
		// %q escape อักขระควบคุมใน path เพื่อไม่ให้ผู้ส่ง request ปลอมบรรทัด log ได้
		log.Printf("method=%q path=%q status=%d duration=%s", r.Method, r.URL.Path, status, time.Since(start))
	})
}

// Recover ดัก panic ระหว่างประมวลผล request และตอบ 500 โดยไม่ส่งรายละเอียดภายในให้ผู้ใช้
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				log.Print("request panic recovered")
				writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
