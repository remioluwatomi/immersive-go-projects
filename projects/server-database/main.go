package main

import(
	"net/http"
	"encoding/json"
	"fmt"
	"log"
  "strings"
	"os"
	"io"
	"context"
	
	"github.com/joho/godotenv"
	"server-database/db"
	"server-database/utils"
	"server-database/types"
	"github.com/jackc/pgx/v5"
)


func writeJsonError(w http.ResponseWriter, statusCode int, message string, code, convIndent int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
  
	res := types.FailureJson{Ok: false, Message: message, Code: code}
	jsonRes, _ := json.MarshalIndent(res, "", strings.Repeat(" ", convIndent))

	w.Write(jsonRes)
}


func setHeaders(w http.ResponseWriter, contentLength, statusCode int) {
 w.Header().Set("Content-Type", "application/json")
 w.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
 w.WriteHeader(statusCode)
}


func handlePostImage(w http.ResponseWriter, r *http.Request, indent int, conn *pgx.Conn) {
    defer r.Body.Close()
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to read request body:", err)
        writeJsonError(w, http.StatusInternalServerError, "Could not read request", 2, indent)
        return
    }

    var image types.Image
    if err := json.Unmarshal(bodyBytes, &image); err != nil {
        fmt.Fprintln(os.Stderr, "Invalid JSON:", err)
        writeJsonError(w, http.StatusBadRequest, "Invalid JSON format", 2, indent)
        return
    }

    if err := db.UploadImage(conn, image); err != nil {
        fmt.Fprintln(os.Stderr, "Failed to upload image:", err)
        writeJsonError(w, http.StatusInternalServerError, "Upload failed", 6, indent)
        return
    }

    response, err := json.MarshalIndent(image, "", strings.Repeat(" ", indent))
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to marshal response:", err)
        writeJsonError(w, http.StatusInternalServerError, "Could not create response", 2, indent)
        return
    }

    setHeaders(w, len(response), http.StatusCreated)
    w.Write(response)
}


func handleGetImages(w http.ResponseWriter, r *http.Request, indent int, conn *pgx.Conn) {
    images, err := db.FetchImages(conn)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to fetch images:", err)
        writeJsonError(w, http.StatusInternalServerError, "Error fetching images", 3, indent)
        return
    }

    response, err := json.MarshalIndent(images, "", strings.Repeat(" ", indent))
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to marshal response:", err)
        writeJsonError(w, http.StatusInternalServerError, "Could not create response", 2, indent)
        return
    }

    setHeaders(w, len(response), http.StatusOK)
    w.Write(response)
}




func main()  {
	err := godotenv.Load()
	if err != nil {
    log.Fatal("err on laoding env..")
	}

	conn, err := db.InitPostgresDB()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
  defer conn.Close(context.Background())

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
    convIndent := utils.GetConvIndent(r.URL)
    
    switch r.Method {
    case http.MethodGet:
        handleGetImages(w, r, convIndent, conn)
    case http.MethodPost:
        handlePostImage(w, r, convIndent, conn)
    default:
        writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed", 1, convIndent)
    }
     
	})

	http.ListenAndServe(":8080", nil)
}
