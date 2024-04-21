package main

import (
	"context"
	"git.snappfood.ir/backend/go/services/bushwack/handler/app"
)

func main() {
	application := app.NewApplication(context.Background())
	application.Setup()
}
