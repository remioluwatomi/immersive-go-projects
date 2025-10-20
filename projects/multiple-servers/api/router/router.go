package router

import (
	"net/http"
	"servers/api/handlers"

	"github.com/jackc/pgx/v5"
)

func ApiRouter(conn *pgx.Conn) *http.ServeMux {
	mux := http.NewServeMux()

  imageHandler := handlers.CreateImageHandler(conn)
	mux.HandleFunc("/images.json", imageHandler.HandleImages)
	
	return mux
}
