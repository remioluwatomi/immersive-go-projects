package db

import (
	"fmt"
	"context"
	"github.com/jackc/pgx/v5"
	"os"
	"strconv"
	"servers/api/models"
)

func InitializeDB() (*pgx.Conn, error) {
		var (
		 dbHost      = os.Getenv("DATABASE_HOST")
		 dbPortStr   = os.Getenv("DATABASE_PORT")
		 dbUser      = os.Getenv("DATABASE_USER")
		 dbName      = os.Getenv("DATABASE_NAME")
		 dbPassword  = os.Getenv("DATABASE_PASSWORD")
	  )
	
		 if dbHost == "" || dbUser == "" || dbName == "" {
        err := fmt.Errorf("Ensure you include all database auth parameters as environment variables")
     return nil, err
	   }
   
		 var dbPort int
		 var err error
	   dbPort, err = strconv.Atoi(dbPortStr)
       if err != nil {
	     return nil, fmt.Errorf("err: DB_PORT contains an invalid DATABASE_PORT")
     }	

		 connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

		 conn, err := pgx.Connect(context.Background(), connStr)

		 if err != nil {
			 err := fmt.Errorf("err: Unable to connect to the db.. %v", err)
			 return nil, err
		 }

		 return conn, nil   
} 


func GetAllImages(conn *pgx.Conn) ([]models.Image, error) {
	var images []models.Image

	rows, err := conn.Query(context.Background(), "SELECT title, alt_text, url FROM public.images;")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var title, url, altText string

		err := rows.Scan(&title, &altText, &url)

		if err != nil {return nil, err}

		images = append(images, models.Image{Title: title, AltText: altText, URL: url})
	}
  
	return images, nil
}

func UploadImage(conn *pgx.Conn, image models.Image) error {
	cmdTag, err := conn.Exec(context.Background(), "INSERT INTO public.images(title, url, alt_text) VALUES($1, $2, $3);", image.Title, image.URL, image.AltText)
  if err != nil {
		return err 
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no row was inserted into the images table")
	}

	return nil
}
