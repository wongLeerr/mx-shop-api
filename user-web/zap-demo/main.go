package main

import (
	"time"

	"go.uber.org/zap"
)

func main() {
	// 适用于生产环境
	logger, _ := zap.NewProduction()
	// logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	url := "www.baidu.com"
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}
