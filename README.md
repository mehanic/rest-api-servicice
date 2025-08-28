# rest-api-servicice

curl -X POST http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -d '{"url":"https://google.com"}'
{"status":"OK","alias":"MIbBwd"}

-----
CONFIG_PATH=./config/config.yaml go run ./cmd/url-shortener/main.go
&{local ./storage/storage.db {localhost:8082 4s 1m0s}}
[23:27:27.940] INFO: starting url-shortener {
  "env": "local"
}
[23:27:27.940] DEBUG: debug messages are enabled 
[23:27:27.940] INFO: logger middleware enabled {
  "component": "middleware/logger"
}
[23:27:27.940] INFO: starting server {
  "address": "localhost:8082"
}
[00:21:21.319] INFO: request body decoded {
  "op": "handlers.url.save.New",
  "request": {
    "url": "https://google.com"
  },
  "request_id": "SkyNet/rGbRrN1LTY-000001"
}
[00:21:21.332] INFO: url added {
  "id": 6,
  "op": "handlers.url.save.New",
  "request_id": "SkyNet/rGbRrN1LTY-000001"
}
[00:21:21.333] INFO: request completed {
  "bytes": 33,
  "duration": "14.037634ms",
  "method": "POST",
  "path": "/url",
  "remote_agent": "127.0.0.1:34728",
  "request_it": "SkyNet/rGbRrN1LTY-000001",
  "status": 200,
  "user_agent": "curl/8.5.0"
}
2025/08/29 00:29:21 [SkyNet/rGbRrN1LTY-000001] "POST http://localhost:8082/url HTTP/1.1" from 127.0.0.1:34728 - 200 33B in 14.076301ms


