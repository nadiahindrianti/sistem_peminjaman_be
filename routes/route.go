package routes

import (
	"log"
	"net/http"
	"sistem_peminjaman_be/controllers"
	"sistem_peminjaman_be/middlewares"
	"sistem_peminjaman_be/repositories"
	"sistem_peminjaman_be/usecases"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func init() {
	middleware.ErrJWTMissing.Code = 401
	middleware.ErrJWTMissing.Message = "Unauthorized"
}

func Init(e *echo.Echo, db *gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	templateMessageRepository := repositories.NewTemplateMessageRepository(db)
	userRepository := repositories.NewUserRepository(db)
	notificationRepository := repositories.NewNotificationRepository(db)
	historySearchRepository := repositories.NewHistorySearchRepository(db)
	labImageRepository := repositories.NewLabImageRepository(db)
	labRepository := repositories.NewLabRepository(db)
	beritaAcaraImageRepository := repositories.NewBeritaAcaraImageRepository(db)
	jadwalRepository := repositories.NewJadwalRepository(db)
	historySeenLabRepository := repositories.NewHistorySeenLabRepository(db)
	suratRekomendasiImageRepository := repositories.NewSuratRekomendasiImageRepository(db)
	peminjamanRepository := repositories.NewPeminjamanRepository(db)
	dashboardRepository := repositories.NewDashboardRepository(db)

	templateMessageUsecase := usecases.NewTemplateMessageUsecase(templateMessageRepository)
	templateMessageController := controllers.NewTemplateMessageController(templateMessageUsecase)

	userUsecase := usecases.NewUserUsecase(userRepository, notificationRepository)
	userController := controllers.NewUserController(userUsecase)

	cloudinaryUsecase := usecases.NewMediaUpload()
	cloudinaryController := controllers.NewCloudinaryController(cloudinaryUsecase)

	notificationUsecase := usecases.NewNotificationUsecase(notificationRepository, templateMessageRepository, userRepository)
	notificationController := controllers.NewNotificationController(notificationUsecase)

	historySearchUsecase := usecases.NewHistorySearchUsecase(historySearchRepository, userRepository)
	historySearchController := controllers.NewHistorySearchController(historySearchUsecase)

	historySeenLabUsecase := usecases.NewHistorySeenLabUsecase(historySeenLabRepository, labRepository, labImageRepository)
	historySeenLabController := controllers.NewHistorySeenLabController(historySeenLabUsecase)

	labUsecase := usecases.NewLabUsecase(labRepository, labImageRepository, historySearchRepository, userRepository, historySeenLabUsecase)
	labController := controllers.NewLabController(labUsecase)

	jadwalUsecase := usecases.NewJadwalUsecase(jadwalRepository, beritaAcaraImageRepository, userRepository)
	jadwalController := controllers.NewJadwalController(jadwalUsecase)

	peminjamanUsecase := usecases.NewPeminjamanUsecase(peminjamanRepository, suratRekomendasiImageRepository, labRepository, labImageRepository, userRepository)
	peminjamanController := controllers.NewPeminjamanController(peminjamanUsecase)

	dashboardUsecase := usecases.NewDashboardUsecase(dashboardRepository, userRepository, peminjamanRepository, jadwalRepository, labRepository)
	dashboardController := controllers.NewDashboardController(dashboardUsecase)

	// Middleware CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Izinkan semua domain
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	api := e.Group("/api/v1")
	public := api.Group("/public")

	// cloudinary
	public.POST("/cloudinary/file-upload", cloudinaryController.FileUpload)
	//public.POST("/cloudinary/url-upload", cloudinaryController.UrlUpload)

	// USER
	api.POST("/login", userController.UserLogin)
	api.POST("/register", userController.UserRegister)
	api.POST("/register/exam", userController.ExamUserRegister)
	api.POST("/register/admin", userController.AdminRegister)

	user := api.Group("/user")
	user.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("user"))

	// user account
	user.Any("", userController.UserCredential)
	user.PUT("/update-password", userController.UserUpdatePassword)
	user.PUT("/update-profile", userController.UserUpdateProfile)
	user.PUT("/update-photo-profile", userController.UserUpdatePhotoProfile)
	user.DELETE("/delete-photo-profile", userController.UserDeletePhotoProfile)

	user.GET("/notification", notificationController.GetNotificationByUserID)

	// user lab
	user.GET("/lab/search", labController.SearchLabAvailable)

	// user jadwal
	user.GET("/jadwal/search", jadwalController.SearchJadwalAvailable)
	user.POST("/jadwal", jadwalController.CreateJadwal)
	public.GET("/user/jadwal/:id", jadwalController.AdminGetJadwalByID)
	user.PUT("/jadwal/:id", jadwalController.UpdateJadwalbyAdmin)

	//user search
	user.GET("/history-search", historySearchController.HistorySearchGetAll)
	user.POST("/history-search", historySearchController.HistorySearchCreate)
	user.DELETE("/history-search/:id", historySearchController.HistorySearchDelete)

	//user search lab
	user.GET("/history-seen-lab", historySeenLabController.GetAllHistorySeenLabs)

	// ADMIN
	admin := api.Group("/admin")
	admin.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("admin"))

	// users @ admin
	admin.GET("/user", userController.UserGetAll)
	admin.GET("/user/detail", userController.UserGetDetail)
	admin.POST("/user/register", userController.UserAdminRegister)
	admin.PUT("/user/update/:id", userController.UserAdminUpdate)
	admin.DELETE("/user/:id", userController.DeleteUser)

	admin.GET("/dashboard", dashboardController.DashboardGetAll)
	admin.GET("/dashboard/filter", dashboardController.DashboardGetByMonth)

	public.GET("/template-message", templateMessageController.GetAllTemplateMessages)
	public.GET("/template-message/:id", templateMessageController.GetTemplateMessageByID)
	public.PUT("/template-message/:id", templateMessageController.UpdateTemplateMessage)
	public.POST("/template-message", templateMessageController.CreateTemplateMessage)

	public.GET("/lab", labController.GetAllLabs)
	public.GET("/lab/:id", labController.GetLabByID)
	admin.PUT("/lab/:id", labController.UpdateLab)
	admin.POST("/lab", labController.CreateLab)
	admin.DELETE("/lab/:id", labController.DeleteLab)

	public.GET("/jadwal", jadwalController.GetAllJadwals)
	public.GET("/jadwal/:id", jadwalController.GetJadwalByID)
	admin.PUT("/jadwal/:id", jadwalController.UpdateJadwal)
	admin.POST("/jadwal", jadwalController.CreateJadwal)
	admin.DELETE("/jadwal/:id", jadwalController.DeleteJadwal)

	user.GET("/peminjaman", peminjamanController.GetAllPeminjamans)
	user.GET("/peminjaman/:id", peminjamanController.GetPeminjamanByID)
	admin.GET("/peminjaman", peminjamanController.GetPeminjamansByAdmin)
	public.GET("/peminjamans", peminjamanController.GetPeminjamansByAdmin)
	public.GET("/peminjaman/:id", peminjamanController.GetPeminjamanByID)
	public.GET("/admin/peminjaman/:id", peminjamanController.AdminGetPeminjamanByID)
	//admin.PUT("/peminjaman/admin/:id", peminjamanController.AdminUpdatePeminjaman)
	admin.PUT("/peminjaman/:id", peminjamanController.UpdatePeminjaman)
	user.POST("/peminjaman", peminjamanController.CreatePeminjaman)

}
