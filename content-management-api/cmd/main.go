package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/service"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/aws"
	appConfig "github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/config"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/database"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/kafka"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/web/controller"
	"github.com/guilhermefaleiros/codestream/content-management-system/pkg"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	cfg := appConfig.LoadConfig()

	connection, err := database.NewConnection(ctx, cfg)

	defer connection.Close()

	if err != nil {
		panic(err)
	}
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.Aws.AccessKeyID, cfg.Aws.SecretAccessKey, ""),
		),
		config.WithRegion(cfg.Aws.Region),
	)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.BootstrapServers,
		"group.id":          cfg.Kafka.GroupId,
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	kafkaVideoUploadedHandler := kafka.NewVideoUploadedHandler(kafkaProducer, cfg.Kafka.VideoUploadedTopic)

	eventMediator := pkg.NewEventMediator()

	err = eventMediator.Register("VideoUploaded", kafkaVideoUploadedHandler)
	if err != nil {
		log.Fatalf("Failed to register handler: %v", err)
	}

	videoRepository := database.NewPGVideoRepository(connection)
	movieRepository := database.NewPGMovieRepository(connection)
	imageRepository := database.NewPGImageRepository(connection)

	storageGateway := aws.NewS3StorageGateway(awsCfg, cfg.Aws.S3.Bucket)
	videoService := service.NewVideoService(storageGateway, videoRepository, eventMediator, cfg.Aws.S3.BaseFolder)
	movieService := service.NewMovieService(movieRepository)
	imageService := service.NewImageService(imageRepository, storageGateway)

	videoController := controller.NewVideoController(r, videoService, cfg.App.MaxFileSize)
	movieController := controller.NewMovieController(r, movieService)
	imageController := controller.NewImageController(r, imageService, cfg.App.MaxFileSize)

	videoController.SetupRoutes()
	movieController.SetupRoutes()
	imageController.SetupRoutes()

	log.Println(fmt.Sprintf("Server started on port %s", cfg.App.Port))
	if err := http.ListenAndServe(cfg.App.Port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
