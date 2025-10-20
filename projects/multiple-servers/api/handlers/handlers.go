package handlers

import (
	"net/http"
  "servers/api/models"
	"servers/api/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"servers/api/utls"
	"fmt"
	"os"
	"io"
	"encoding/json"
	"context"
)

type ImageHandler struct {
	Conn *pgxpool.Pool
}

func CreateImageHandler(conn *pgxpool.Pool) *ImageHandler {
	return &ImageHandler{Conn: conn}
}

func (i *ImageHandler) HandleImages(w http.ResponseWriter, r *http.Request) {
	conn := i.Conn
	cxt := r.Context()
	switch r.Method {	
    case http.MethodGet:
			getAllImages(cxt, conn, w, r)
			break
			
		case http.MethodPost:
			postImage(cxt, conn, w, r)
			break 
    
		case http.MethodOptions:
			// code to handle cors pre-flight..
			break 
		}
}

func getAllImages(cxt context.Context, conn *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	imgs, err := db.GetAllImages(cxt, conn)

  if err != nil {
		fmt.Fprintln(os.Stderr,"error fetching images from the db: \n", err)
		errJson := models.CreateFailedJson("Cannot fetch images at the moment.", 2)
		utls.JSONError(w, errJson, http.StatusInternalServerError, utls.ResponseIndent(r.URL))
		return
	}
  
	utls.JSON(w, imgs, http.StatusOK, utls.ResponseIndent(r.URL))
	
}


func postImageErrResWrapper(err error, customErrorMessage string, w http.ResponseWriter, errCode, statusCode, responseIndent int) {
  fmt.Fprintln(os.Stderr, customErrorMessage, err)
	errJson := models.CreateFailedJson(customErrorMessage, errCode)
	utls.JSONError(w, errJson, statusCode, responseIndent)
}

func postImage(cxt context.Context, conn *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	indent := utls.ResponseIndent(r.URL)
  
	if err != nil {
		postImageErrResWrapper(err, "Unable to read request body. ", w, 3, http.StatusInternalServerError, indent)
		return
	}

	var image models.Image 
	if err := json.Unmarshal(data, &image); err != nil {
    postImageErrResWrapper(err, "Invalid Json data in request body. ", w, 3, http.StatusBadRequest, indent)
		 return
	} 

	if err := db.UploadImage(cxt, conn, image); err != nil {
		postImageErrResWrapper(err, "Upload failed. ", w, 4, http.StatusInternalServerError, indent)
		 return
	}

	var images []models.Image
	images = append(images, image)
  utls.JSON(w, images, http.StatusCreated, utls.ResponseIndent(r.URL))

}
