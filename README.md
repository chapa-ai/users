Тестовое задание GO

1. Описать proto файл с сервисом из 3 методов: добавить пользователя, удалить пользователя, список пользователей
2. Реализовать gRPC сервис на основе proto файла на Go
3. Для хранения данных использовать PostgreSQL
4. на запрос получения списка пользователей данные будут кешироваться в redis на минуту и брать из редиса
5. При добавлении пользователя делать лог в clickHouse
6. Добавление логов в clickHouse делать через очередь Kafka


**Installation**<br>

° Start server in grpc-server<br>
° Start client in grpc-client<br>
° See saved data in postgres and logs in clickhouse

**Requirements**<br>
° Golang = 1.17 <br>
° Redis = 7.0.4 <br>
° Kafka = 3.2 <br>
° Clickhouse = 22.6.9 <br>






