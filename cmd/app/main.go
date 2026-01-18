package main

import (
	"fmt"

	"github.com/Weit145/REST_API_golang/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// This is the entry point of the application.
	// You can initialize your application here.
}
