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
	router.PUT("/api/v1/admin/:username/username", adminController.UpdateUsernameAdmin)
	router.PUT("/api/v1/admin/:username/password", adminController.UpdatePasswordAdmin)
	router.GET("/api/v1/admin/:id_admin", adminController.GetAdminById)

	kontenRepo := repositories.NewKontenRepository()
	kontenService := services.NewKontenService(kontenRepo, db)
	kontenController := controllers.NewKontenController(kontenService)

	router.GET("/api/v1/kontak/:id_kontak", kontenController.GetKontak)
	router.PUT("/api/v1/kontak/:id_kontak", kontenController.UpdateKontak)

	aparatRepo := repositories.NewAparatRepository()
	aparatService := services.NewAparatService(aparatRepo, db)
	aparatController := controllers.NewAparatController(aparatService)

	router.POST("/api/v1/aparat", aparatController.CreateAparat)
	router.GET("/api/v1/aparat", aparatController.GetAllAparat)
	router.GET("/api/v1/aparat/:id_aparat", aparatController.GetAparatById)
	router.PUT("/api/v1/aparat/:id_aparat", aparatController.UpdateAparat)
	router.DELETE("/api/v1/aparat/:id_aparat", aparatController.DeleteAparat)
	router.DELETE("/api/v1/bulk/aparat", aparatController.BulkDeleteAparat)

	beritaRepo := repositories.NewBeritaRepository()
	beritaService := services.NewBeritaService(beritaRepo, db)
	beritaController := controllers.NewBeritaController(beritaService)

	router.POST("/api/v1/berita", beritaController.CreateBerita)
	router.GET("/api/v1/berita", beritaController.GetAllBerita)
	router.GET("/api/v1/berita/:id_berita", beritaController.GetBeritaById)
	router.PUT("/api/v1/berita/:id_berita", beritaController.UpdateBerita)
	router.DELETE("/api/v1/berita/:id_berita", beritaController.DeleteBerita)
	router.POST("/api/v1/photo/berita", beritaController.CreatePhoto)
	router.DELETE("/api/v1/photo/berita/:filename", beritaController.DeletePhotoByFilename)
	router.DELETE("/api/v1/bulk/photo/berita", beritaController.BulkDeletePhoto)

	pendudukRepo := repositories.NewPendudukRepository()
	pendudukService := services.NewPendudukService(pendudukRepo, db)
	pendudukController := controllers.NewPendudukController(pendudukService)

	router.GET("/api/v1/penduduk", pendudukController.GetPenduduk)
	router.PUT("/api/v1/penduduk", pendudukController.UpdatePenduduk)

	// Serve static files from uploads directory
	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

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
