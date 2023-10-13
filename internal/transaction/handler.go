package transaction

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"transaction-service/internal/errorhandling"
)

// Storer is the interface of the transaction business service expected by the handler that deals with storing
// transactions.
type Storer interface {
	Store(transaction StoreRequest) (StoreResponse, error)
}

// ConfigureStoreHandler configures the supplied router with a store handler that uses the supplied service to store
// transactions.
func ConfigureStoreHandler(router *gin.Engine, service Storer) {
	router.POST("/transaction", NewStoreHandler(service))
}

// NewStoreHandler is responsible for mapping the incoming 'fetch transaction' http request into the call to the
// business service and mapping the result back to a http response.
func NewStoreHandler(service Storer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var request StoreRequest
		if err := ctx.Bind(&request); err != nil {
			ctx.Error(errors.New(errorhandling.BadRequest))
			return
		}
		response, err := service.Store(request)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}

// Fetcher is the interface of the transaction business service expected by the handler that deals with fetching
// transactions.
type Fetcher interface {
	Fetch(ctx context.Context, transactionID, country string) (FetchResponse, error)
}

// ConfigureFetchHandler configures the supplied router with a fetch handler that uses the supplied service to fetch
// transactions.
func ConfigureFetchHandler(router *gin.Engine, service Fetcher) {
	router.GET("/transaction/:id", NewFetchHandler(service))
}

// NewFetchHandler is responsbile for mapping the incoming 'store transaction' http request into the call to the
// business service and mapping the result back to a http response.
func NewFetchHandler(service Fetcher) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		country := ctx.Query("country")
		response, err := service.Fetch(ctx, id, country)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}
