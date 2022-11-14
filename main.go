package main

import (
	"os"

	"github.com/heroku/urinception/cmd/web"
)

func main() {
	web.Start(os.Getenv("PORT"))
}
