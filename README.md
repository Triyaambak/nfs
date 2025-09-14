```env
./.env
API_PORT = 3001
```

```env
./api/.env
API_PORT = 3001

NFS_DIR = ".././nfs"

SECRET_KEY = "dKCbzxjR8UIY6GO1LMnpSXUsRWvToBqZ"
```

```cmd
docker compose up
```

Routes

##### Sample JWT Auth Token - eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxIiwibmFtZSI6ImFsaWNlIiwiZ2lkIjoiMjAiLCJncm91cCI6InN0YWZmIiwiaWF0IjoxNTE2MjM5MDIyfQ.0CW30reVqoXW4diiNxK7R-WCGbhxDpWK7O-5TtSTiCw

- / : Interactive file system
- ls : http://localhost:3001/ls/{directory} || http://localhost:3001/ls/{file_path}
- mkdir : http://localhost:3001/ls/{directory}
- touch : http://localhost:3001/ls/{file_path}
- cat : http://localhost:3001/cat/{file_path}
- mv : http://localhost:3001/ls/{source_file_path}-{dst_file_path}
- echo : http://localhost:3001/echo/ {text} >> {directory}
