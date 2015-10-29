var express = require('express');
var app = express();
var redis = require('redis');
var httpProxy = require('http-proxy'),

var client = redis.createClient(6379, 'redis');

var proxy = new httpProxy.RoutingProxy();

function apiProxy(host, port) {
  return function(req, res, next) {
    if(req.url.match(new RegExp('^\/api\/'))) {
      proxy.proxyRequest(req, res, {host: host, port: port});
    } else {
      next();
    }
  }
}

client.on("error", function (err) {
    console.error("Redis error", err);
});

app.get('/', function (req, res) {
    res.redirect('/index.html');
});

app.use(express.static('static'));
app.use(apiProxy('api', 8080));


var server = app.listen(80, function () {
    console.log('Web running on port 80');
});
