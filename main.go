// ไฟล์ main.go: จุดเริ่มต้นโปรแกรม เลือกคำสั่ง เปิดการเชื่อมต่อ PostgreSQL และจัดการวงจรชีวิต HTTP server
package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"thaifood-backend/auth"
	"thaifood-backend/handler"
	"thaifood-backend/internal/config"
	"thaifood-backend/mw"
	"thaifood-backend/store"
)

// ฝัง schema และข้อมูลเริ่มต้นลงใน binary เพื่อให้ migrate/seed ใช้ได้โดยไม่ต้องมีไฟล์ SQL ข้างโปรแกรม
//
//go:embed schema.sql seed.sql
var sqlFiles embed.FS

// main เรียกโปรแกรมหลักและออกด้วยสถานะผิดพลาดเมื่อเริ่มทำงานไม่สำเร็จ
func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

// run เลือก serve/migrate/seed/admin เตรียมฐานข้อมูล และเปิด API เมื่อเลือก serve
func run() error {
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if command != "serve" && command != "migrate" && command != "seed" && command != "admin" {
		return errors.New("usage: thaifood [serve|migrate|seed|admin EMAIL]")
	}
	var c config.Config
	var err error
	// โหมด server ต้องใช้ค่าด้าน JWT/CORS ด้วย ส่วนคำสั่งฐานข้อมูลต้องการเพียง DATABASE_URL
	if command == "serve" {
		c, err = config.Load()
		if err != nil {
			return err
		}
	} else {
		c.DatabaseURL = os.Getenv("DATABASE_URL")
		if c.DatabaseURL == "" {
			return errors.New("DATABASE_URL is required")
		}
	}
	poolConfig, err := pgxpool.ParseConfig(c.DatabaseURL)
	if err != nil {
		return errors.New("invalid DATABASE_URL configuration")
	}
	// จำกัดจำนวน connection และอายุการใช้งาน เพื่อคืน connection ที่ว่างหรือใช้นานกลับสู่ฐานข้อมูล
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 0
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.ConnConfig.ConnectTimeout = 10 * time.Second
	// ไม่ใช้ named prepared statements เพื่อให้ทำงานร่วมกับ transaction pooler ได้
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec
	poolConfig.ConnConfig.RuntimeParams["statement_timeout"] = "10000"
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		cancel()
		return errors.New("database initialization failed")
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	cancel()
	if err != nil {
		return errors.New("database connection failed; check DATABASE_URL and Supabase project status")
	}
	if command != "serve" {
		return databaseCommand(pool, command)
	}
	// ประกอบ store, JWT issuer และ routes แล้วกำหนด deadline ของ context แต่ละ request ไว้ 15 วินาที
	h := &handler.Handler{Store: &store.Store{Pool: pool}, Tickets: &auth.Issuer{Secret: []byte(c.JWTSecret), TTL: time.Hour}}
	routes := h.Routes()
	bounded := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		routes.ServeHTTP(w, r.WithContext(ctx))
	})
	server := &http.Server{
		Addr:              ":" + c.Port,
		Handler:           mw.Logging(mw.Recover(mw.CORS(c.AllowedOrigins, bounded))),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20,
	}
	// เมื่อได้รับสัญญาณหยุด ให้ server รอ request ที่กำลังทำงานก่อนปิด ภายในเวลาสูงสุด 15 วินาที
	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("ThaiFood Guide API listening on :%s", c.Port)
		serverErr <- server.ListenAndServe()
	}()
	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.New("HTTP server failed")
		}
	case <-stopCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return errors.New("HTTP shutdown timed out")
		}
	}
	return nil
}

// databaseCommand รัน SQL ที่ฝังใน binary หรือเลื่อนสมาชิกที่ระบุด้วยอีเมลให้เป็น admin
func databaseCommand(pool *pgxpool.Pool, command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if command == "admin" {
		// เลื่อนสิทธิ์ได้เฉพาะบัญชีที่สมัครไว้แล้ว โดยเทียบอีเมลแบบไม่แยกตัวพิมพ์เล็กใหญ่
		if len(os.Args) != 3 {
			return errors.New("usage: thaifood admin EMAIL (register the account first)")
		}
		tag, err := pool.Exec(ctx, `UPDATE members SET role='admin',updated_at=now() WHERE lower(email)=lower($1)`, strings.TrimSpace(os.Args[2]))
		if err != nil {
			return errors.New("admin promotion failed")
		}
		if tag.RowsAffected() != 1 {
			return errors.New("member not found; register the account first")
		}
		log.Print("Registered member promoted to admin")
		return nil
	}
	file := "schema.sql"
	if command == "seed" {
		file = "seed.sql"
	}
	data, err := sqlFiles.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", file, err)
	}
	if _, err = pool.Exec(ctx, string(data)); err != nil {
		return fmt.Errorf("%s failed; check database schema and permissions", command)
	}
	log.Printf("Database %s completed", command)
	return nil
}
