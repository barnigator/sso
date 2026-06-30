package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/barnigator/sso/internal/auth/delivery/grpc"
	"github.com/barnigator/sso/internal/auth/deps"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService deps.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{log, gRPCServer, port}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const fn = "grpcapp.Run"

	log := a.log.With(
		slog.String("fn", fn),
		slog.Int("port", a.port),
	)

	log.Info("starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) {
	const fn = "grpcapp.Stop"

	a.log.With(slog.String("fn", fn)).
		Info("stopping gRPC server")

	stopped := make(chan struct{})

	go func() {
		a.gRPCServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		a.log.Warn("graceful shutdown timed out, forced stop")
		a.gRPCServer.Stop()
	case <-stopped:
		a.log.Info("gRPC server stopped gracefully")
	}
}
