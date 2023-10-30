# Сервис для обогащения данных о людях 🚀

Этот проект представляет собой REST-сервис, написанный на Golang, который принимает ФИО, обогащает данные о возрасте,
поле и национальности, сохраняет их в базу данных PostgreSQL и предоставляет методы для получения, обновления и 
удаления данных.

## Требования ⚙️
Для запуска этого проекта вам потребуется следующее:
- Docker
- Создать и заполнить по данному шаблону .env файл:
```
#POSTGRES ENVIRONMENTS
PGUSER=user
PGPASSWORD=password
PGHOST=postgres
PGPORT=5432
PGDATABASE=db
PGSSLMODE=disable

#SERVER ENVIRONMENTS
# debug or release
ENVIRONMENT=debug

HTTP_PORT=8080

#API
GET_AGE_API=https://api.agify.io
GET_COUNTRY_API=https://api.nationalize.io
GET_GENDER_API=https://api.genderize.io

```

## Запуск 🛠️
Для запуска выполните в терминале команду ```make compose-up```, после чего сервер будет запущен на localhost на указанном
вами порту.
Для остановки сервера нужно прописать команду ```make compose-down```


## REST Методы 📎
Сервис предоставляет следующие REST-методы:

### Получение данных с различными фильтрами и пагинацией:

- Метод: GET
- URL: /people
- Пример запроса:

```
localhost:8080/people?limit=100&offset=0&min_age=46&max_age=47&gender=male&name=Muslim
```
Поддерживает данные фильтрации:
- min_age
- max_age
- gender
- name
- (limit, offset) - пагинация
### Удаление по идентификатору:

- Метод: DELETE
- URL: /people
- Пример запроса:

```
curl --location --request DELETE 'localhost:8080/people' \
--header 'Content-Type: application/json' \
--data '{
    "peopleId": 1
}'
```

### Изменение сущности:

- Метод: PUT
- URL: /people
- Пример запроса:

```
curl --location --request PUT 'localhost:8080/people' \
--header 'Content-Type: application/json' \
--data '{
    "id": 10,
    "name": "Имя",
    "surname": "Фамилия",
    "patronymic": "Отчество"
}'
```

### Добавление новых людей:

- Метод: POST
- URL: /people
- Пример запроса:

```
curl --location --request POST 'localhost:8080/people' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Ivan",
    "surname": "Ivanov",
    "patronymic": "Ivanovich"
}'
```

Обогащение данных о возрасте, поле и национальности будет автоматически выполнено и сохранено в базу данных PostgreSQL.
