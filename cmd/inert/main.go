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
  flag.BoolVar(&flagSsl, "ssl", false, "Run the server with SSL support (requires key.pem and cert.pem files)")
  flag.IntVarP(&flagPort, "port", "p", 8334, "Port to bind server too")
  flag.StringVarP(&flagCert, "cert", "c", "cert.pem", "SSL Certificate file")
  flag.StringVarP(&flagKey, "key", "k", "key.pem", "SSL Kery file")
  flag.StringVarP(&flagDir, "dir", "d", "./static", "Directory to serve as root")
}

func main() {

  // Initialise flag values
  flag.Parse()

  // Create port string for HTTP(S) listener
  finalPort := fmt.Sprintf(":%d", flagPort)
  fmt.Println("Booting sever on", finalPort)

  // Boot Server with or without SSL support based on the "--ssl" flag
  if flagSsl {
    fmt.Println("SSL Mode enabled")
    log.Fatal(http.ListenAndServeTLS(finalPort, flagCert, flagKey, http.FileServer(http.Dir(flagDir))))
  } else {
    fmt.Println("SSL Mode disabled")
    log.Fatal(http.ListenAndServe(finalPort, http.FileServer(http.Dir(flagDir))))
  }
}
