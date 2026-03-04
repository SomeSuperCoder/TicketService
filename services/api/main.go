package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/OnlineShop/handlers"
	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	appConfig := internal.LoadAppConfig()
	pool, repo, redisClient := internal.DatabaseConnect(ctx, appConfig)
	defer pool.Close()

	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	humaConfig := huma.DefaultConfig(
		"Ticket Huma + Gin API",
		"1.0.0",
	)
	humaConfig.Servers = []*huma.Server{
		{URL: "http://localhost:8888/api/v1", Description: "Local API version 1"},
	}
	api := humagin.NewWithGroup(r, apiGroup, humaConfig)

	MountRoutes(api, repo, pool, redisClient, appConfig)

	r.Run(fmt.Sprintf(":%s", appConfig.Port))
}

func MountRoutes(api huma.API, repo *repository.Queries, pool *pgxpool.Pool, redisClient *redis.Client, appConfig *internal.AppConfig) {
	categoryHandler := handlers.CategoryHandler{Repo: repo}
	{
		huma.Register(api, huma.Operation{
			OperationID: "get-categories",
			Method:      http.MethodGet,
			Path:        "/public/categories",
			Description: "Get a tree of categories and subcategories",
			Tags:        []string{"Categories"},
		}, categoryHandler.Get)
	}

	ticketHandler := handlers.TicketHandler{Repo: repo, Pool: pool}

	// Ticket Routes
	{
		// Create
		huma.Register(api, huma.Operation{
			OperationID: "create-ticket",
			Method:      http.MethodPost,
			Path:        "/public/tickets",
			Description: "Create a ticket with default values (status=init, is_hidden=false)",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Post)

		// Read
		huma.Register(api, huma.Operation{
			OperationID: "get-ticket",
			Method:      http.MethodGet,
			Path:        "/tickets/{id}",
			Description: "Get a ticket by ID",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Get)

		huma.Register(api, huma.Operation{
			OperationID: "list-tickets",
			Method:      http.MethodGet,
			Path:        "/tickets",
			Description: "List all tickets with pagination",
			Tags:        []string{"Tickets"},
		}, ticketHandler.List)

		// Update
		huma.Register(api, huma.Operation{
			OperationID: "update-ticket",
			Method:      http.MethodPatch,
			Path:        "/tickets/{id}",
			Description: "Update a ticket",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Update)

		huma.Register(api, huma.Operation{
			OperationID: "delete-ticket",
			Method:      http.MethodDelete,
			Path:        "/tickets/{id}",
			Description: "Permanently delete a ticket",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Delete)

		// Merge
		huma.Register(api, huma.Operation{
			OperationID: "merge-duplicates",
			Method:      http.MethodPost,
			Path:        "/tickets/merge",
			Description: "Merge duplicate tickets",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Merge)
	}

	statisticsHandler := handlers.StatisticsHandler{Repo: repo}
	{
		huma.Register(api, huma.Operation{
			OperationID: "stats-summary",
			Method:      http.MethodGet,
			Path:        "/statistics/summary",
			Description: "General ticket stats",
			Tags:        []string{"Statistics"},
		}, statisticsHandler.GetSummary)
		huma.Register(api, huma.Operation{
			OperationID: "category-stats",
			Method:      http.MethodGet,
			Path:        "/statistics/categories",
			Description: "Stats for each category",
			Tags:        []string{"Statistics"},
		}, statisticsHandler.GetCategoryStatistics)
		huma.Register(api, huma.Operation{
			OperationID: "stats-dynamics",
			Method:      http.MethodGet,
			Path:        "/statistics/dynamics",
			Description: "Ticket dynamics over time (received and resolved)",
			Tags:        []string{"Statistics"},
		}, statisticsHandler.GetDynamics)
	}

	// History Routes
	historyHandler := handlers.HistoryHandler{Repo: repo, Pool: pool}
	{
		huma.Register(api, huma.Operation{
			OperationID: "get-ticket-history",
			Method:      http.MethodGet,
			Path:        "/tickets/{id}/history",
			Description: "Get history of actions and status changes for a ticket",
			Tags:        []string{"History"},
		}, historyHandler.GetTicketHistory)

		huma.Register(api, huma.Operation{
			OperationID: "get-recent-history",
			Method:      http.MethodGet,
			Path:        "/history/recent",
			Description: "Get recent history across all tickets",
			Tags:        []string{"History"},
		}, historyHandler.GetRecentHistory)
	}

	// Comments Routes
	commentsHandler := handlers.CommentsHandler{Repo: repo, Pool: pool}
	{
		huma.Register(api, huma.Operation{
			OperationID: "create-comment",
			Method:      http.MethodPost,
			Path:        "/tickets/{id}/comments",
			Description: "Add a comment to a ticket",
			Tags:        []string{"Comments"},
		}, commentsHandler.Post)

		huma.Register(api, huma.Operation{
			OperationID: "get-comments",
			Method:      http.MethodGet,
			Path:        "/tickets/{id}/comments",
			Description: "Get all comments for a ticket",
			Tags:        []string{"Comments"},
		}, commentsHandler.Get)
	}

	// Heatmap Routes
	heatmapHandler := handlers.HeatmapHandler{Repo: repo, Pool: pool}
	{
		huma.Register(api, huma.Operation{
			OperationID: "get-heatmap-points",
			Method:      http.MethodGet,
			Path:        "/heatmap/points",
			Description: "Get points for rendering on the map with intensity",
			Tags:        []string{"Heatmap"},
		}, heatmapHandler.GetPoints)

		huma.Register(api, huma.Operation{
			OperationID: "get-heatmap-stats",
			Method:      http.MethodGet,
			Path:        "/heatmap/stats",
			Description: "Get heatmap statistics including top problem locations",
			Tags:        []string{"Heatmap"},
		}, heatmapHandler.GetStats)
	}
}
