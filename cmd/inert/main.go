package main

import (
  "fmt"
  "log"
  "net/http"
  flag "github.com/spf13/pflag"
)

var (
  flagPort int
  flagSsl bool
  flagCert string
  flagKey string
  flagDir string
)

func init() {
  // --ssl - Enables SSL mode (requires --cert and --key to be present)
  flag.BoolVar(&flagSsl, "ssl", false, "Run the server with SSL support (requires key.pem and cert.pem files)")
  
  // --port, -p - Sets server port
  flag.IntVarP(&flagPort, "port", "p", 8334, "Port to bind server too")
  
  // --cert, -c - Sets the path to the certificate file to read for SSL mode
  flag.StringVarP(&flagCert, "cert", "c", "cert.pem", "SSL Certificate file")

  // --key, -k - Sets the path to the key file to read for SSL mode
  flag.StringVarP(&flagKey, "key", "k", "key.pem", "SSL Kery file")

  // --dir, -d - Sets the path to the directory that will act as the root directory
  flag.StringVarP(&flagDir, "dir", "d", "./static", "Directory to serve as root")
}

func main() {

  // Initialise flag values
  flag.Parse()

  // Create port string for HTTP(S) listener
  finalPort := fmt.Sprintf(":%d", flagPort)
  fmt.Println("Booting sever on", finalPort)

	// Create FS Sever Handler
	fsHandle := makeFS(flagDir)


  // Boot Server with or without SSL support based on the "--ssl" flag
  if flagSsl {
    fmt.Println("SSL Mode enabled")
    log.Fatal(http.ListenAndServeTLS(finalPort, flagCert, flagKey, fsHandle))
  } else {
    fmt.Println("SSL Mode disabled")
    log.Fatal(http.ListenAndServe(finalPort, fsHandle))
  }
}
