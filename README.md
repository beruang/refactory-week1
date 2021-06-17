# Notes API Documentation
Notes API is backend application created for Refactory.id RSP recruitment test.

## Technology Used
|Name|Description|
|----|-----------|
|[Golang](https://golang.org)| Golang programing language|
|[Echo](https://github.com/labstack/echo) | Web framework for Golang |
|[Redis/ Cache](https://github.com/go-redis) | Using local cache and store at redis |
|[Viper](https://github.com/spf13/viper)| Library for handle reading configuration |
|[Migrate](https://github.com/golang-migrate/migrate)| Database migration |
|[Sqlx](https://github.com/jmoiron/sqlx)|Library for provide database connection|
|[Casbin](https://github.com/casbin/casbin)|Access-control library|
|[Swaggo](https://github.com/swaggo/swag)|API Documentation library provide swagger|
|[GoMail](https://gopkg.in/gomail.v2)|library for sending email|
|[GoValidator](github.com/go-playground/validator)|Additional library for validating request|

## Feature
- Sending email verification with queue system
- Access-Control
- API Web Documentation
- Migrations

## API Documentation
```
{{host}}/api/swagger/index.html
```
## Environment
```
web_address=
web_read_timeout=
web_write_timeout=
web_shutdown_timeout=
pg_host=
pg_port=
pg_user=
pg_password=
pg_name=
rdr_host=
rdr_port=
rdr_db=
rdr_pool=
mail_host=
mail_port=
mail_user=
mail_password=
is_dev=true
```

## Contacts
If you have any issue or feature request, please contact me. PR are welcome
- [https://github.com/beruang](https://github.com/beruang)
- [zackymughnii@gmail.com](zackymughnii@gmail.com)
