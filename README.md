
## API

| Ручка | Метод | Влияние |Параметры|
|-------|-------|---------|---------|
|/api/v1/rent|`POST`|Создает запись аренды, присваивает ей id|`user`|
|/api/v1/rent_start|`POST`|Помечает сессию аренды как начатую|`rentid`|
|/api/v1/rent_stop|`POST`|Помечает сессию аренды как завершенную|`rentid`|
|/api/v1/upload_photo|`POST`|Присвоить сессии аренды фото пруф|`rentid`|
|/api/v1/download_photo|`GET`|Получить фотопруфы с сессии аренды|`rentid`|

## CREDS: 

DB_CONNECTION_STRING = user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable
MINIO_HOST = localhost:9000
MINIO_USER = root
MINIO_ACCESS_KEY = password
SERVER_PORT = :8080

:warning: Обязательно бакет должен называться `testbucket`