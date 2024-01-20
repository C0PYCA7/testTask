# Реализация сервиса обогащения данных о людях

## Описание

Сервис, разработанный на языке программирования Go, 
предназначен для получения информации о людях через API, 
обогащения ответа данными о возрасте, поле и национальности, 
а затем сохранения полученных данных в базе данных PostgreSQL.

## Использованные библиотеки
- [Chi](https://github.com/go-chi/chi/)
- [Gorm](https://github.com/go-gorm/gorm)
- [Validator](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Chi-render](https://github.com/go-chi/render)
- [Cleanenv](https://github.com/ilyakaznacheev/cleanenv)

### Получение данных
- `GET /people`: Получение списка людей с различными фильтрами и пагинацией.

### Управление данными
- `POST /people`: Добавление нового человека в формате JSON.
- `PUT /people/{id}`: Изменение информации о человеке по идентификатору.
- `DELETE /people/{id}`: Удаление человека по идентификатору.

### Пример использования сервиса
- `GET curl --location 'http://localhost:8080/list'`
- `GET(с фильтрами и пагинацией) curl --location 'http://localhost:8080/list?size=10&nationality=RU&page=1' \
  --data ''`
- `POST curl --location 'http://localhost:8080/create' \
  --header 'Content-Type: text/plain' \
  --header 'Authorization: Basic cm9vdDpyb290' \
  --data '{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"
  }'`
- `PUT curl --location --request PUT 'http://localhost:8080/update/30' \
  --header 'Content-Type: application/json' \
  --data '{
  "name": "Sergey",
  "surname": "Ushakov",
  "patronymic": "Vasilevich",
  "gender": "male",
  "nationality": "RU",
  "age": 19
  }
  '`
- `DELETE curl --location --request DELETE 'http://localhost:8080/delete/1'`

## Обогащение данных

При добавлении нового человека или изменении существующей записи, сервис обогащает информацию о возрасте, поле и национальности, используя следующие внешние API:
1. Возраст - [Agify API](https://api.agify.io/?name=Dmitriy).
2. Пол - [Genderize API](https://api.genderize.io/?name=Dmitriy).
3. Национальность - [Nationalize API](https://api.nationalize.io/?name=Dmitriy).

Обогащенные данные сохраняются в базе данных PostgreSQL.

## Логирование

Код сервиса покрыт debug- и info-логами.

## Конфигурационные данные

Для безопасности конфигурационные данные вынесены в файл `.env`.
