package db

import(
	"fmt"
	"context"
  "os"
  "github.com/jackc/pgx/v5"
	"strconv"
	"server-database/types"
)

func InitPostgresDB() (*pgx.Conn, error)  {
	var (
		dbHost      = os.Getenv("DB_HOST")
		dbPortStr        = os.Getenv("DB_PORT")
		dbUser      = os.Getenv("DB_USER")
		dbName      = os.Getenv("DB_NAME")
		dbPassword  = os.Getenv("DB_PASSWORD")
	)
	
	dbPort := 5432 
	if dbPortStr != "" {
    if port, err := strconv.Atoi(dbPortStr); err == nil {
			dbPort = port
		}else {
			return nil, fmt.Errorf("Invalid DB_PORT value: %v", err)
		}
	}

	if dbHost == "" || dbUser == "" || dbName == "" {
     err := fmt.Errorf("Ensure you include all database auth parameters as environment variables")
     return nil, err
	 }

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		err := fmt.Errorf("unable to connect to the db..", err)
		return nil, err 
	}

	return conn, nil
}

func FetchImages(conn *pgx.Conn) ([]types.Image, error) {
	 var images []types.Image
		
		rows, err :=  conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
		defer rows.Close()

		if err != nil {
			return nil, fmt.Errorf("Error fetching images from the db.. \n err: %v", err)
		}
		
		for rows.Next() {
       var title, url, altText string
       err = rows.Scan(&title, &url, &altText)
      
			 images = append(images, types.Image{Title: title, AltText: altText, URL: url})
		}
    
		return images, nil
}

func UploadImage(conn *pgx.Conn, image types.Image) error {
	cmdTag, err := conn.Exec(context.Background(), "INSERT INTO public.images(title, url, alt_text) VALUES($1, $2, $3)", image.Title, image.URL, image.AltText)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
    return fmt.Errorf("no row was inserted into images table")
	}
	return nil
}


