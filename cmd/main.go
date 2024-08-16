package main

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/cache"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/config"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/database"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/logger"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/metrics"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/msgqueue"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/proto"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/taskservice"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	defaultEnv := "local"
	env, exists := os.LookupEnv("APP_ENV")
	if exists {
		defaultEnv = env
	}
	// Load configuration
	conf, err := config.LoadConfig(fmt.Sprintf("./%s-env-config.yaml", defaultEnv))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	taskLogger := logger.NewLogger(conf.Logging)
	metricConf := conf.MetricsConfig
	metricConf.Logger = taskLogger
	metrics.InitializePrometheus(&metricConf)
	//conf.Database
	// Setup database, cache, message queue, and logger
	db := database.NewDatabase(conf.Database)

	cacheInst := cache.NewCache(conf.Cache)

	msgQueue := msgqueue.NewMesageQueue(&conf.MessageQueue)

	tskService := taskservice.NewTaskService(db, cacheInst, msgQueue, taskLogger)

	server := proto.NewTaskServiceServer(tskService)

	//// Create and start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go metrics.StartServer()

	s := grpc.NewServer()
	proto.RegisterTaskServiceServer(s, server)

	log.Println("Starting gRPC server on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	//taskRepo := repository.NewPostgresRepository(db)
	//taskCache := cache.newRedisCache()
	//taskQueue := queue.NewRabbitMQ()
	//taskLogger := logger.NewConsoleLogger()
	//
	//// Create TaskService
	//taskService := service.NewTaskService(taskRepo, taskCache, taskQueue, taskLogger)
	//
	//// Create and start gRPC server
	//lis, err := net.Listen("tcp", ":50051")
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//
	//grpcServer := grpc.NewTaskServiceServer(taskService)
	//
	//s := grpc.NewServer()
	//pb.RegisterTaskServiceServer(s, grpcServer)
	//
	//log.Println("Starting gRPC server on port 50051...")
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
}
