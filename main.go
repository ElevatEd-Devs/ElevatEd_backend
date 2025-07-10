package main

import (
	"context"
	"elevated_backend/database"
	"elevated_backend/router"
)

func main() {
	conn, _ := database.ConnectDB()
	defer conn.Close(context.Background())
	router.SetRouter(conn)
}
