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
		"Ticke API Huma + Gin API",
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
	ticketHandler := handlers.TicketHandler{Repo: repo}

	// Ticket Routes
	{
		// Create
		huma.Register(api, huma.Operation{
			OperationID: "create-ticket",
			Method:      http.MethodPost,
			Path:        "/tickets",
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

		huma.Register(api, huma.Operation{
			OperationID: "get-tickets-by-status",
			Method:      http.MethodGet,
			Path:        "/tickets/status/{status}",
			Description: "Get tickets by status",
			Tags:        []string{"Tickets"},
		}, ticketHandler.GetByStatus)

		huma.Register(api, huma.Operation{
			OperationID: "get-tickets-by-subcategory",
			Method:      http.MethodGet,
			Path:        "/tickets/subcategory/{subcategoryId}",
			Description: "Get tickets by subcategory",
			Tags:        []string{"Tickets"},
		}, ticketHandler.GetBySubcategory)

		huma.Register(api, huma.Operation{
			OperationID: "get-tickets-by-department",
			Method:      http.MethodGet,
			Path:        "/tickets/department/{departmentId}",
			Description: "Get tickets by department",
			Tags:        []string{"Tickets"},
		}, ticketHandler.GetByDepartment)

		// Search
		huma.Register(api, huma.Operation{
			OperationID: "search-tickets-by-meaning",
			Method:      http.MethodGet,
			Path:        "/tickets/search/meaning",
			Description: "Search for tickets based upon semantic meaning",
			Tags:        []string{"Tickets"},
		}, ticketHandler.SearchByMeaning)

		// Update
		huma.Register(api, huma.Operation{
			OperationID: "update-ticket",
			Method:      http.MethodPatch,
			Path:        "/tickets/{id}",
			Description: "Update a ticket",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Update)

		huma.Register(api, huma.Operation{
			OperationID: "update-ticket-status",
			Method:      http.MethodPatch,
			Path:        "/tickets/{id}/status",
			Description: "Update ticket status",
			Tags:        []string{"Tickets"},
		}, ticketHandler.UpdateStatus)

		// Delete/Hide
		huma.Register(api, huma.Operation{
			OperationID: "hide-ticket",
			Method:      http.MethodDelete,
			Path:        "/tickets/{id}/hide",
			Description: "Soft delete a ticket (hide it)",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Hide)

		huma.Register(api, huma.Operation{
			OperationID: "delete-ticket",
			Method:      http.MethodDelete,
			Path:        "/tickets/{id}",
			Description: "Permanently delete a ticket",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Delete)

		// Count
		huma.Register(api, huma.Operation{
			OperationID: "count-tickets",
			Method:      http.MethodGet,
			Path:        "/tickets/count",
			Description: "Get total count of non-hidden tickets",
			Tags:        []string{"Tickets"},
		}, ticketHandler.Count)
	}
}
