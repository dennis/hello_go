package main

import (
	"github.com/dennis/hello_go/app"
)

func main() {
	app := app.App{}
	app.Initialize()
	app.Run()
}
