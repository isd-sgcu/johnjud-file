package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/isd-sgcu/johnjud-file/cfgldr"
	"github.com/isd-sgcu/johnjud-file/database"
	"github.com/isd-sgcu/johnjud-file/pkg/client/bucket"
	imageRepo "github.com/isd-sgcu/johnjud-file/pkg/repository/image"
	imageSvc "github.com/isd-sgcu/johnjud-file/pkg/service/image"
	imagePb "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		sig := <-s

		log.Info().
			Str("service", "graceful shutdown").
			Msgf("got signal \"%v\" shutting down service", sig)

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Error().
				Str("service", "graceful shutdown").
				Msgf("timeout %v ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("cleaning up: %v", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Error().
						Str("service", "graceful shutdown").
						Err(err).
						Msgf("%v: clean up failed: %v", innerKey, err.Error())
					return
				}

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("%v was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}

func main() {
	conf, err := cfgldr.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "file").
			Msg("Failed to load config")
	}

	db, err := database.InitPostgresDatabase(&conf.Database, conf.App.Debug)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "file").
			Msg("Failed to init postgres connection")
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "file").
			Msg("Failed to load AWS SDK config")
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", conf.App.Port))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "file").
			Msg("Failed to start service")
	}

	grpcServer := grpc.NewServer()

	awsClient := s3.NewFromConfig(sdkConfig)
	bucketClient := bucket.NewClient(conf.S3, awsClient)

	imageRepository := imageRepo.NewRepository(db)

	imageService := imageSvc.NewService(bucketClient, imageRepository)

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	imagePb.RegisterImageServiceServer(grpcServer, imageService)

	reflection.Register(grpcServer)
	go func() {
		log.Info().
			Str("service", "file").
			Msgf("JohnJud file starting at port %v", conf.App.Port)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().
				Err(err).
				Str("service", "file").
				Msg("Failed to start service")
		}
	}()

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"server": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
		"database": func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				return nil
			}
			return sqlDB.Close()
		},
	})

	<-wait

	grpcServer.GracefulStop()
	log.Info().
		Str("service", "file").
		Msg("Closing the listener")
	lis.Close()
	log.Info().
		Str("service", "file").
		Msg("End the program")
}
