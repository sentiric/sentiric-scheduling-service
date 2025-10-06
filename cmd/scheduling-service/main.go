// sentiric-scheduling-service/cmd/scheduling-service/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/sentiric/sentiric-scheduling-service/internal/config"
	"github.com/sentiric/sentiric-scheduling-service/internal/logger"
	"github.com/sentiric/sentiric-scheduling-service/internal/server"

	externalv1 "github.com/sentiric/sentiric-contracts/gen/go/sentiric/external/v1"
)

var (
	ServiceVersion string
	GitCommit      string
	BuildDate      string
)

const serviceName = "scheduling-service"

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Kritik Hata: Konfigürasyon yüklenemedi: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(serviceName, cfg.Env, cfg.LogLevel)

	log.Info().
		Str("version", ServiceVersion).
		Str("commit", GitCommit).
		Str("build_date", BuildDate).
		Str("profile", cfg.Env).
		Msg("🚀 Sentiric Scheduling Service başlatılıyor...")

	// HTTP ve gRPC sunucularını oluştur
	grpcServer := server.NewGrpcServer(cfg.CertPath, cfg.KeyPath, cfg.CaPath, log)
	httpServer := startHttpServer(cfg.HttpPort, log)

	// gRPC Handler'ı kaydet
	externalv1.RegisterSchedulingServiceServer(grpcServer, &schedulingHandler{})

	// gRPC sunucusunu bir goroutine'de başlat
	go func() {
		log.Info().Str("port", cfg.GRPCPort).Msg("gRPC sunucusu dinleniyor...")
		if err := server.Start(grpcServer, cfg.GRPCPort); err != nil && err.Error() != "http: Server closed" {
			log.Error().Err(err).Msg("gRPC sunucusu başlatılamadı")
		}
	}()

	// Graceful shutdown için sinyal dinleyicisi
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Warn().Msg("Kapatma sinyali alındı, servisler durduruluyor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Stop(grpcServer)
	log.Info().Msg("gRPC sunucusu durduruldu.")

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP sunucusu düzgün kapatılamadı.")
	} else {
		log.Info().Msg("HTTP sunucusu durduruldu.")
	}

	log.Info().Msg("Servis başarıyla durduruldu.")
}

func startHttpServer(port string, log zerolog.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok"}`)
	})

	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{Addr: addr, Handler: mux}

	go func() {
		log.Info().Str("port", port).Msg("HTTP sunucusu (health) dinleniyor")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP sunucusu başlatılamadı")
		}
	}()
	return srv
}

// =================================================================
// GRPC HANDLER IMPLEMENTASYONU (Placeholder)
// =================================================================

type schedulingHandler struct {
	externalv1.UnimplementedSchedulingServiceServer
}

func (*schedulingHandler) CreateAppointment(ctx context.Context, req *externalv1.CreateAppointmentRequest) (*externalv1.CreateAppointmentResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("rpc", "CreateAppointment").Str("user_id", req.GetUserId()).Logger()
	log.Info().Str("title", req.GetTitle()).Msg("CreateAppointment isteği alındı (Placeholder)")

	// Simüle edilmiş başarılı randevu
	return &externalv1.CreateAppointmentResponse{
		AppointmentId: "appt-fake-98765",
		Success:       true,
	}, nil
}
