package app

import (
	"github.com/gin-gonic/gin"

	"transaction-service/internal/errorhandling"
	"transaction-service/internal/transaction"
)

// newRouter returns a new default gin http router with the application's routes wired to their handlers and middleware
// configured.
func newRouter(deps Dependencies) *gin.Engine {
	router := gin.Default()
	router.Use(errorhandling.NewMiddleware)
	transaction.ConfigureStoreHandler(router, deps.TxnService)
	transaction.ConfigureFetchHandler(router, deps.TxnService)
	return router
}
