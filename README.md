### Веб-сервер
Рассмотрен общий случай, когда один пользователь может иметь несколько счетов.

Реализован следующий функционал:

По части данных 
- выдачи списка пользователей (с пагинацией и фильтрацией)
- выдачи данных конкретного пользователя
- добавления пользователя
- удаление пользователя
- редактирование данных пользователя

По части счётов с балансами
- выдачи списка счётов (с пагинацией и фильтрацией)
- выдачи данных конкретного счёта
- добавления счета
- удаление счёта
- редактирование данных 
- пополнение счета
- перевод определенной суммы

### Примеры
В файлах __cards.sh__ и __users.sh__ (в папке
__scritps__) можно рассмотреть некоторые примеры позитивных и негативных сценариев по всем вышеописанным действиям.

Также при помощи данных скриптов можно осуществить пополнение тестовой базы.

### Сборка и запуск
Осуществляется при помощи Docker.

Сборка версии с Postgres:
```
make build
```
Запуск:

```
make run
```