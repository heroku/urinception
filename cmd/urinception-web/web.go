package main

import (
	"os"

	"github.com/heroku/urinception/web"
)

func main() {
	web.Start(os.Getenv("PORT"))
}
