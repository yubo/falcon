server {
  listen 7000;
  location ^~ / {
    proxy_pass http://localhost:8002;
  }

  location ^~ /v1.0 {
    proxy_pass http://localhost:8001;
  }
}
