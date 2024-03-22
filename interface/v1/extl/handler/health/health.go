package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"dot/core"
	"dot/infrastructure/v1/persistence/mongodb"
	"dot/infrastructure/v1/persistence/redisdb"
)

type Handler struct {
	mdbc mongodb.Client
	rdbc redisdb.Client
}

func New(mdbc mongodb.Client, rdbc redisdb.Client) *Handler {
	return &Handler{
		mdbc: mdbc,
		rdbc: rdbc,
	}
}

// Health godoc
// @Summary 	Health Check.
// @Description	Health Check.
// @Tags 		Health
// @Accept	 	json
// @Produce 	json
// @Success		200
// @Failure 	400
// @Failure 	401
// @Failure 	500
// @Router /api/v1/health [get]
func (handler *Handler) Health(c echo.Context) error {
	ic := core.NewInternalContext(uuid.New().String())
	var healthStatus interface{}
	checker := health.NewChecker(
		health.WithCacheDuration(1*time.Second),
		health.WithTimeout(10*time.Second),
		health.WithCheck(health.Check{
			Name:    "mongodb",
			Timeout: 2 * time.Second,
			Check: func(ctx context.Context) error {
				return handler.mdbc.Ping(ctx, readpref.Primary())
			},
		}),
		health.WithCheck(health.Check{
			Name:    "redis",
			Timeout: 2 * time.Second,
			Check: func(ctx context.Context) error {
				_, err := handler.rdbc.Ping(ctx)
				if err != nil {
					return err
				}
				return nil
			},
		}),
		health.WithStatusListener(func(ctx context.Context, state health.CheckerState) {
			healthStatus = state.Status
		}),
	)

	result := checker.Check(ic.ToContext())

	if fmt.Sprintf("%v", healthStatus) == "down" {
		return c.JSON(http.StatusServiceUnavailable, result)
	}

	return c.JSON(http.StatusOK, result)
}
