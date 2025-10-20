package static

import(
	"fmt"
  "flag"
	"net/http"
  "log"
  "strconv"
)

var assetsDir string
var serverPort int  
func init() {
	flag.StringVar(&assetsDir, "path", "assets", "include a path to your assets directory")
  flag.IntVar(&serverPort, "port", 8082, "include a port number for your server")
}


func Run() {
	flag.Parse()
	fmt.Printf("port: %d", serverPort)
  
	staticFileDir := fmt.Sprintf("./%s", assetsDir)

  fs := http.FileServer(http.Dir(staticFileDir))

	http.Handle("/", fs)

 portStr := fmt.Sprintf(":%s", strconv.Itoa(serverPort))
 log.Fatal(http.ListenAndServe(portStr, nil))

}
