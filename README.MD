# gRPC Thumbnail Service
gRPC Thumbnail Service- сервис, включающий клиент и сервер, который позволяет скачивать thumbnail'ы видео с YouTube.

## Технологии
- gRPC-Go - фреймворк для работы с gRPC на Go
- Gorm - ORM-библиотека для работы с БД
- SQLite - база данных
- Bufconn - для тестирования gRPC без запуска сервиса


## Функциональность
- Получение thumbnail'а YouTube-ролика
- Сохаранение и передача thumbnail'ов из кэша
- gRPC-клиент для скачивания и сохранения в файлы формата jpg thumbnail'ов
- Клиент принимает флаг --async для асинхронного скачивания


## Запуск
- Распакуйте архив thumbnails.tar.gz

- Скачайте зависимости:
```
go mod download
```
- Сгенерируйте pb-файлы:
```
protoc --proto_path=./proto --go_out=./proto --go-grpc_out=./proto ./proto/thumbnails.proto
```
### Запуск сервера
```
go run cmd/server/main.go
```

### Запуск клиента
```
go run cmd/client/main.go [--async async download] [Youtube URL's]
```
Ссылки на Yotube-видео должны иметь формат `https://www.youtube.com/watch?v=*unique_video_code*`

Пример:
```
go run cmd/client/main.go --async async download https://www.youtube.com/watch?v=BV7BMtBNmyI https://www.youtube.com/watch?v=ljnZpLSfo2k
```

При скачивании, рядом с `cmd/server/main.go`, создастся папка `/images` в которую будут скачаны thumbnail'ы. Имена файлов также являются уникальным кодом видео, при подставлении имени файла в конец `https://www.youtube.com/watch?v=` даст ссылку на видео, thumbnail'ом которого является файл.

Для доступа к YouTube из РФ может понадобится VPN или ПО позволящее обходить Deep Packet Inspection ([GoodbyeDpi](https://github.com/ValdikSS/GoodbyeDPI))

## Запуск тестов
```
go test -race ./...
```