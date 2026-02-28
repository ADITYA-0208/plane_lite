package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"planelite-backend/cmd/server/api"
	"planelite-backend/internal/activity"
	"planelite-backend/internal/auth"
	"planelite-backend/internal/common"
	"planelite-backend/internal/config"
	"planelite-backend/internal/middleware"
	"planelite-backend/internal/notification"
	"planelite-backend/internal/notification/providers"
	"planelite-backend/internal/project"
	"planelite-backend/internal/task"
	"planelite-backend/internal/user"
	"planelite-backend/internal/workspace"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	cfg := config.LoadEnv()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := config.ConnectMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	db := client.Database(cfg.DBName)
	if err := config.EnsureIndexes(context.Background(), db); err != nil {
		log.Printf("warning: ensure indexes: %v", err)
	}

	// Repositories
	userRepo := user.NewRepository(db)
	workspaceRepo := workspace.NewRepository(db)
	membershipRepo := workspace.NewMembershipRepository(db)
	projectRepo := project.NewRepository(db)
	taskRepo := task.NewRepository(db)

	userSvc := user.NewService(userRepo)
	authSvc := auth.NewService(userSvc, cfg)
	workspaceSvc := workspace.NewService(workspaceRepo, membershipRepo)
	projectSvc := project.NewService(projectRepo)
	taskSvc := task.NewService(taskRepo)
	activitySvc := activity.NewService(db)
	_ = activitySvc

	inApp := providers.NewInAppProvider()
	whatsApp := providers.NewWhatsAppProvider()
	_ = notification.NewService(inApp, whatsApp)

	authHandler := auth.NewHandler(authSvc, cfg)
	userHandler := user.NewHandler(userSvc)
	workspaceHandler := workspace.NewHandler(workspaceSvc)
	projectHandler := project.NewHandler(projectSvc)
	taskHandler := task.NewHandler(taskSvc)

	authMW := middleware.Auth(authSvc)
	adminOnly := middleware.RequireRole(common.RoleAdmin)
	workspaceAccess := &middleware.WorkspaceAccess{
		Membership: workspaceSvc,
		GetWorkspaceID: func(r *http.Request) (primitive.ObjectID, bool) {
			idHex := r.PathValue("id")
			if idHex == "" {
				return primitive.ObjectID{}, false
			}
			id, err := primitive.ObjectIDFromHex(idHex)
			return id, err == nil
		},
	}

	mw := api.Middleware{
		Auth:            authMW,
		AdminOnly:       adminOnly,
		WorkspaceAccess: workspaceAccess.Middleware(),
	}

	mux := http.NewServeMux()
	api.RegisterHealth(mux, client)
	api.RegisterAuth(mux, authHandler)
	api.RegisterUser(mux, userHandler, mw)
	api.RegisterWorkspace(mux, workspaceHandler, mw)
	api.RegisterProject(mux, projectHandler, mw)
	api.RegisterTask(mux, taskHandler, mw)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
