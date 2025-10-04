package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/desa-sukamaju-api/config"
	"github.com/syrlramadhan/desa-sukamaju-api/controllers"
	"github.com/syrlramadhan/desa-sukamaju-api/repositories"
	"github.com/syrlramadhan/desa-sukamaju-api/services"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		fmt.Printf("error saat memuat file .env: %v\n", errEnv)
		return
	}

	port := os.Getenv("APP_PORT")
	fmt.Println("api berjalan di port:" + port)

	db, err := config.ConnectToDatabase()
	if err != nil {
		fmt.Println("error saat koneksi ke database:", err)
		return
	}

	router := httprouter.New()

	adminRepo := repositories.NewAdminRepository()
	adminService := services.NewAdminService(adminRepo, db)
	adminController := controllers.NewAdminController(adminService)

	router.POST("/api/v1/admin/login", adminController.LoginAdmin)

	handler := CorsMiddleware(router)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	errServer := server.ListenAndServe()
	if errServer != nil {
		fmt.Printf("error saat memulai server: %v\n", errServer)
		return
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
