package main 

import(
	"net/http"
  "flag"
	"log"
	"strconv"
	"fmt"
)

var serverPort int
func init() {

  flag.IntVar(&serverPort, "port", 80, "Port number needs to be configured for server")

}


func main() {
  flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world."))
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello!"))
	})

	log.Println("server is running on port ", serverPort)

	http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(serverPort)), nil)

}
