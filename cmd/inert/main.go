package main

import (
  "log"
  "net/http"
)

func main() {
  log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", http.FileServer(http.Dir("./static"))))
}
