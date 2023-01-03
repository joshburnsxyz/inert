import Koa from 'koa';
import serve from 'koa-static';
import path from "path";

const app = new Koa;

// Mount middleware
app.use(serve(path.join(__dirname)));

// Start server and listen for connections
app.listen(3000);