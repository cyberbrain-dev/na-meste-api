package main

import (
	"fmt"

	"github.com/cyberbrain-dev/na-meste-api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: setup logger

	// TODO: init database
}
