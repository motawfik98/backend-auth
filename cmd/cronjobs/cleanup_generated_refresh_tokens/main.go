package main

import (
	"backend-auth/internal/logger"
	"backend-auth/internal/utils/connection"
)

func main() {
	cronjob := connection.InitializeCronJob()
	err := cronjob.CleanupGeneratedRefreshToken()
	if err != nil {
		logger.LogFailure(err, "error deleting expired (generated refresh tokens) from db")
		panic(err)
	}
}
