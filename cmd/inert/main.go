package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
)

var (
	flagPort    int
	flagSsl     bool
	flagCert    string
	flagKey     string
	flagDir     string
	flagHelp    bool
	flagVersion bool
)

const INERT_VERSION = "0.1.0"

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

	// --help, -h - Displays help message
	flag.BoolVarP(&flagHelp, "help", "h", false, "Displays help message")

	// --version, -v - Displays version information
	flag.BoolVarP(&flagVersion, "version", "v", false, "Displays version information")

	flag.Usage = func() {
		fmt.Printf("Usage: inert [OPTIONS]\n\n")
		flag.PrintDefaults()
	}
}

func main() {

	// Initialise flag values
	flag.Parse()

	if flagHelp {
		flag.Usage()
		return
	}

	if flagVersion {
		fmt.Printf("Inert v%s\n", INERT_VERSION)
		return
	}

	// Create port string for HTTP(S) listener
	finalPort := fmt.Sprintf(":%d", flagPort)
	fmt.Println("Booting sever on", finalPort)

	// Create FS Sever Handler
	fsHandle, err := makeFS(flagDir)
	if err != nil {
		log.Fatal(err)
	}

	// Boot Server with or without SSL support based on the "--ssl" flag
	if flagSsl {
		fmt.Println("Running in HTTPS mode")
		log.Fatal(http.ListenAndServeTLS(finalPort, flagCert, flagKey, fsHandle))
	} else {
		fmt.Println("Running in HTTP mode")
		log.Fatal(http.ListenAndServe(finalPort, fsHandle))
	}
}
