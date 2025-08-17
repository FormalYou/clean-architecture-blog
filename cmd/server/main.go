package main

import (
	"go.uber.org/zap"

	"github.com/FormalYou/clean-architecture-blog/cmd/server/option"
	zaplog "github.com/FormalYou/clean-architecture-blog/internal/infrastructure/log"
)

func main() {
	// Defer logger sync
	defer zaplog.GetLogger().Sync()

	router := option.SetupRouter("./configs")

	// 6. Start Server
	zaplog.GetLogger().Info("Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		zaplog.GetLogger().Fatal("could not run server", zap.Error(err))
	}
}
