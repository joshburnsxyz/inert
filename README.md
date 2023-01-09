# inert
Simple server to handle and provide static files, and assets over a URL.

# Usage & Installation

```
$ git clone https://github.com/joshburnsxyz/inert
$ cd inert
$ yarn
$ node index.js [OPTIONS]
```

Store all files you wish to serve under the `./public` directory. Then visit `localhost:3000/foo.png`.

## Options

- `-p/--port` - Set port to bind server too (default is 3000)
- `-d/--dir` - Set the directory to serve from (default is ./public)

# License

Free and Open-Source software under the MIT license. 
