package routes

import (
	"jasamc/cmd/app/handler"
	"jasamc/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, jasaHandler handler.JasaHandler, jwtSecret string) {
	// cek request logging
	router.Use(middleware.RequestLogger())

	// router.GET("/coba-get", userHandler.CobaGet)
	// router.POST("/coba-post", userHandler.CobaPost)
	groupV1 := router.Group("/api/v1")

	groupV1.GET("/categories", jasaHandler.GetAllCategoryJasa)
	groupV1.GET("/category/:id", jasaHandler.GetCategoryJasaByID)

	groupV1.GET("/services", jasaHandler.GetAllJasa)
	groupV1.GET("/service/:serviceID", jasaHandler.GetService)

	// route delete services, berarti dengan service_media
	// dan juga storagenya yang ada di service itu

	// ini jangan di implemetasikan terlebih dulu, karena butuh yang lainnya
	// karena masih ada service_spesificationnya, sama service_spesificationnya_value
	// dan itu belum ter implementasikan

	// add photo and url unutk service

	// route hapus hanya di service media, dan juga di storage yang media itu
	// hapusnya di media tertentu yang ingin dihapus, hapus medianya, dan di storagenya

	// route service spesification

	// route value service spesification

	// routes admin only
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AdminMiddleware(jwtSecret, "admin"))

	// category
	admin.POST("/category", jasaHandler.CreateCategoryJasa)
	admin.DELETE("/category/:id", jasaHandler.DeleteCategoryJasa)
	admin.PATCH("/category/:id/icon", jasaHandler.UpdateCategoryIcon)
	admin.PATCH("/category/:id/status", jasaHandler.SetStatusJasaCategory)

	// jasa
	admin.POST("/service", jasaHandler.CreateService)
	admin.PATCH("/service/:serviceID/status", jasaHandler.SetStatusService)
	admin.DELETE("/service/:serviceID", jasaHandler.DeleteService)
	admin.POST("/service/:serviceID/media", jasaHandler.AddServiceMedia)
	admin.DELETE("/service/:serviceID/media/:mediaID", jasaHandler.DeleteMediaInService)
	admin.POST("/service-spesification", jasaHandler.CreateServiceSpesification)
	admin.DELETE("/service/:serviceID/spesification/:specID", jasaHandler.DeleteServiceSpesification)
	admin.PATCH("/service/:serviceID/spesification/:specID/status", jasaHandler.ToggleServiceSpesificationStatus)
	admin.PATCH("/service/:serviceID/spesification/:specID/required", jasaHandler.ToggleServiceSpesificationRequired)
	admin.POST("/service-spesification-value", jasaHandler.CreateServiceSpesificationValue)
	admin.PATCH(
		"/service/:serviceID/specification/:specID/value/:valueID",
		jasaHandler.UpdateServiceSpesificationValue,
	)

}
