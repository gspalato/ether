package main

import (
	"context"

	"github.com/joho/godotenv"

	"unreal.sh/ether/internal/server"
)

func main() {
	godotenv.Load(".env");

	ctx := context.Background();
	server.Start(ctx);
}
