package api 

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"flag"
	"net/http"
	"log"
	"servers/api/db"
	"servers/api/router"
	"strconv"
	"strings"
)

var serverPort int 
var allowedOrigins []string

func init() {
err := godotenv.Load()
	if err != nil {
		log.Fatal("err on loading env..")
	}

	flag.IntVar(&serverPort, "port", 8081, "include a port number for your server")

	allowedOriginsVar := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins = strings.Split(allowedOriginsVar, ",")
}


func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")        
    if r.Method == "OPTIONS" {
      w.WriteHeader(http.StatusOK)
      return
      }
      next.ServeHTTP(w, r)
	})
}


func Run() {
	flag.Parse()
	
	conn, err := db.InitializeDB()
  
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	defer conn.Close()
  
  mux := router.ApiRouter(conn)
	handler := enableCorsMiddleware(mux)
	
  log.Println("Server running on port ", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(serverPort)), handler)
} 
