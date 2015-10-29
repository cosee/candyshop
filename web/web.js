var express = require('express');
var app = express();
var redis = require('redis');
var proxy = require('express-http-proxy');

var client = redis.createClient(6379, 'redis');



client.on("error", function (err) {
    console.error("Redis error", err);
});

app.get('/', function (req, res) {
    res.redirect('/index.html');
});

app.use(express.static('static'));
app.use('/api', proxy('api:8080', {
  forwardPath: function(req, res) {
    return require('url').parse(req.url).path;
  }
}));
var server = app.listen(80, function () {
    console.log('Web running on port 80');
});
