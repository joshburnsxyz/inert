# inert
Secure static file server

## Usage

```
Usage of inert:
  -c, --cert string   SSL Certificate file (default "cert.pem")
  -d, --dir string    Directory to serve as root (default "./static")
  -k, --key string    SSL Kery file (default "key.pem")
  -p, --port int      Port to bind server too (default 8334)
      --ssl           Run the server with SSL support (requires key.pem and cert.pem files)
```

If you require a Self signed SSL certificate for testing purposes you can use `make ssl_keys` to generate a `cert.pem`
and `key.pem` file that will allow you to run the server in SSL mode like so...

```
$ inert --ssl -c ./cert.pem -k ./key.pem -d ./my_static_dir
```

Or you can serve insecurely over HTTP like so...

```
$ inert -d ./my_static_dir
```

## Installation

Build the server binary with `make`

```
$ make
```