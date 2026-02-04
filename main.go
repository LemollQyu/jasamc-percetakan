package main

import (
	"fmt"
	"jasamc/cmd/app/handler"
	"jasamc/cmd/app/repository"
	"jasamc/cmd/app/resource"
	"jasamc/cmd/app/service"
	"jasamc/cmd/app/storage"
	"jasamc/cmd/app/usecase"
	"jasamc/config"
	grpcJasa "jasamc/grpc"
	"jasamc/infrastructure/log"
	"jasamc/middleware"
	"jasamc/proto/jasapb"
	"jasamc/routes"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Service Jasa")

	cfg := config.LoadConfig()
	fmt.Println("Config semua disni")
	fmt.Println("APP CONFIG:", cfg.App)
	fmt.Println("DATABASE CONFIG:", cfg.Database)
	fmt.Println("REDIS CONFIG:", cfg.Redis)
	fmt.Println("PATH UPLOADS:", cfg.Storage)

	db := resource.InitDB(&cfg)
	redis := resource.InitRedis(&cfg)
	log.SetupLogger()

	jasaRepostory := repository.NewJasaRepository(db, redis)
	jasaService := service.NewJasaService(*jasaRepostory)
	jasaStorage := storage.NewLocalStorage(cfg.Storage.UploadBaseDir, cfg.App.Url)
	jasaUsecase := usecase.NewJasaUsecase(
		*jasaService,
		jasaStorage,
	)

	jasaHandler := handler.NewJasaHandler(*jasaUsecase)

	port := cfg.App.Port

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS([]string{"http://localhost:3000", "http://localhost:3001"}))
	router.Use(gin.Logger())
	router.Static("/static", "./uploads")
	routes.SetupRoutes(router, *jasaHandler, cfg.Secret.JWTSecret)

	// ---- HTTP SERVER ----
	go func() {
		log.Logger.Printf("HTTP server running on port : %s", port)
		if err := router.Run(":" + port); err != nil {
			log.Logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	// ---- gRPC SERVER ----
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Logger.Fatalf("Failed to listen gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	jasapb.RegisterJasaServiceServer(
		grpcServer,
		&grpcJasa.GRPCServer{JasaUsecase: *jasaUsecase},
	)

	for service, info := range grpcServer.GetServiceInfo() {
		log.Logger.Println("gRPC Service:", service)
		for _, method := range info.Methods {
			log.Logger.Println("  └─ Method:", method.Name)
		}
	}

	log.Logger.Println("gRPC server running on port :50052")

	if err := grpcServer.Serve(lis); err != nil {
		log.Logger.Fatalf("Failed to serve gRPC: %v", err)
	}

}
