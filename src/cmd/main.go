package main

import (
	"fmt"
	"os"

	"github.com/hellomyzn/nf-analysis/internal/controller"
	"github.com/hellomyzn/nf-analysis/internal/repository"
	"github.com/hellomyzn/nf-analysis/internal/service"
)

func main() {
	repo := repository.NewNetflixRepository()
	srv := service.NewNetflixService(repo)
	ctrl := controller.NewNetflixController(srv)

	if err := ctrl.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Println("Netflix CSV transformation completed!")
}
