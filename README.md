## Запуск
Сборка версии с Postgres:
```bash
$ docker build -t go_server . --target=server
```

Запуск:
```bash
$ docker run -p 10000:10000 -p 5440:5432 go_server
```