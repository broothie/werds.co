package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "[werds] ", log.LstdFlags)

func main() {
	loggerMiddleware := newLoggerMiddleware(logger)

	handler := loggerMiddleware(http.HandlerFunc(handler))
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logger.Printf("serving @ %s", addr)
	logger.Panic(http.ListenAndServe(addr, handler))
}
