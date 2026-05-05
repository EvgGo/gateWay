package main

import (
	"context"
	"gateWay/internal/config"
	"gateWay/internal/grpc"
	"gateWay/internal/routes"
	"gateWay/pkg/logger"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	testsv1 "github.com/EvgGo/proto/proto/gen/go/tests"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net/http"
)

func main() {

	//if err := godotenv.Load(".env"); err != nil {
	//	log.Println("Внимание: файл .env не найден, используются переменные окружения по умолчанию")
	//	panic(err)
	//}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad("CONFIG_PATH_GATEWAY")

	l, err := logger.SetupLogger(cfg.Env, "", cfg.LogLevel, cfg.LogFile)
	if err != nil {
		log.Fatalf("Ошибка при инициализации логгера: %v", err)
	}

	userServiceConn, err := connection.ConnectWithRetry(ctx, cfg.Auth.Host, cfg.Auth.Port, cfg.DialConfig, l)
	if err != nil {
		l.Error("Error on connecting to the AuthProf-service:", err)
		return
	}

	workspaceConn, err := connection.ConnectWithRetry(ctx, cfg.WorkSpace.Host, cfg.WorkSpace.Port, cfg.DialConfig, l)
	if err != nil {
		l.Error("Error on connecting to the TeamAndProjects-service:", err)
		return
	}

	testinConn, err := connection.ConnectWithRetry(ctx, cfg.Testing.Host, cfg.Testing.Port, cfg.DialConfig, l)
	if err != nil {
		l.Error("Error on connecting to the Test-service:", err)
		return
	}

	defer func(userServiceConn *grpc.ClientConn) {
		err = userServiceConn.Close()
		if err != nil {
			l.Warn("Unable to close connection AuthProf-service:", err)
		}
	}(userServiceConn)

	authClient := authv1.NewAuthClient(userServiceConn)
	profileClient := authv1.NewUserProfileClient(userServiceConn)
	workspaceTeamsClient := workspacev1.NewTeamsClient(workspaceConn)
	workspaceProjectsClient := workspacev1.NewProjectsClient(workspaceConn)
	skillClient := authv1.NewSkillsClient(userServiceConn)
	testsClient := testsv1.NewAdaptiveTestingClient(testinConn)

	router := routes.NewRoutes(l, authClient, profileClient, workspaceTeamsClient, workspaceProjectsClient,
		skillClient, testsClient)
	router.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"), //The url pointing to API definition
	))

	l.Info("Starting server:", slog.String("address", cfg.HTTPServer.Address))
	server := &http.Server{
		Addr:              cfg.HTTPServer.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.ReadTimeout,
		WriteTimeout:      cfg.HTTPServer.WriteTimeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	if err = server.ListenAndServe(); err != nil {
		l.Error("failed to start server", "error", err)
	}

	l.Error("server stopped")
}
