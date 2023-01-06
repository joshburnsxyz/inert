const Koa = require("koa")
const serve = require("koa-static")
const path = require("path")
const compress = require("koa-compress")
const fs = require("fs")
const yargs = require("yargs")

// Process command line
const argv = yargs
      .option('port', {
	  alias: 'p',
	  description: 'set server port to bind too',
	  type: 'string'
      })
      .option('dir', {
	  alias: 'd',
	  description: 'set directory to serve from (./public by default)',
	  type: 'string'
      })
      .help()
      .alias('help', 'h').argv;

const app = new Koa;

let myDir = path.join(__dirname, "./public");
if (argv.dir) {
    myDir = argv.dir;
}

let port = 3000;
if (argv.port) {
    port = argv.port;
}

// Ensure public directory exists before booting server
try{
     fs.lstatSync(myDir).isDirectory()
}catch(e){
   if(e.code == 'ENOENT'){
       fs.mkdir(myDir, () => {});
   }
}

// Mount middleware
app.use(serve(myDir));
app.use(function(req, res, next){
  console.log(req);
  next();
});
app.use(function(req, res, next){
  res.status(200);
  next();
});
app.use(compress({
  br: true
}));

// Start server and listen for connections
app.listen(port, () => {
    console.log("Starting Inert - static asset delivery service.");
    console.log("SERVER BOUND TO PORT "+port);
});
