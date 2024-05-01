package main

import (
  "fmt"
  "log"
  "net/http"
  flag "github.com/spf13/pflag"
)

var (
  flagPort int
  flagCert string
  flagKey string
  flagDir string
)

func init() {
  flag.IntVarP(&flagPort, "port", "p", 8334, "Port to bind server too")
  flag.StringVarP(&flagCert, "cert", "c", "cert.pem", "SSL Certificate file")
  flag.StringVarP(&flagKey, "key", "k", "key.pem", "SSL Kery file")
  flag.StringVarP(&flagDir, "dir", "d", "./static", "Directory to serve as root")
}

func main() {
  flag.Parse()
  finalPort := fmt.Sprintf(":%d", flagPort)
  log.Fatal(http.ListenAndServeTLS(finalPort, flagCert, flagKey, http.FileServer(http.Dir(flagDir))))
}
