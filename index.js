const Koa = require("koa")
const serve = require("koa-static")
const path = require("path")
const compress = require("koa-compress")

const app = new Koa;

// Mount middleware
app.use(serve(path.join(__dirname)));
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
app.listen(3000);
