package router

import (
	"net/http"
	"servers/api/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ApiRouter(conn *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

  imageHandler := handlers.CreateImageHandler(conn)
	mux.HandleFunc("/images.json", imageHandler.HandleImages)
	
	return mux
}
