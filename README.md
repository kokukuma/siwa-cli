# siwa-cli

## How to use
### Apple developer center
+ Create ServiceID
+ Setting RedirectURI
  + set verification file to your domain
  + set local redirector to the redirectURI
+ Download Apple Key

### Use it
```
go run .
```

### local redirector
``` php
<?php
$url = "http://localhost:8080";
$url .= "?code=" .$_GET['code'];
$url .= "&id_token=" .$_GET['id_token'];
$url .= "&state=" .$_GET['state'];
$url .= "&user=" .$_GET['user'];
header("location: $url");
```

``` python
import BaseHTTPServer, SimpleHTTPServer
import ssl

path = "/etc/letsencrypt/live/hostname"

class ServerHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def do_GET(self):
        SimpleHTTPServer.SimpleHTTPRequestHandler.do_GET(self)

    def do_POST(self):
        self.data_string = self.rfile.read(int(self.headers['Content-Length']))
        if self.path == '/localredirect':
            url = 'http://localhost:8080?'
        self.send_response(302)
        self.send_header('Location', url + self.data_string)
        self.end_headers()

# httpd = BaseHTTPServer.HTTPServer(('10.146.0.36', 443), SimpleHTTPServer.SimpleHTTPRequestHandler)
httpd = BaseHTTPServer.HTTPServer(('10.146.0.36', 443), ServerHandler)
httpd.socket = ssl.wrap_socket (httpd.socket, keyfile= path + 'privkey.pem', certfile=path + 'fullchain.pem', server_side=True)
httpd.serve_forever()
```
