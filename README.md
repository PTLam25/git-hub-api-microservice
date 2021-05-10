# GOLANG Microservice with GitHub API integration

## Все вокруг Domains(Models)
Самая главная часть сервиса это domain, модели, данные, с которым мы работаем: запрашиваем, передаем или изменяем.

Поэтому всегда надо начинать с проектирования domains(models). При этому domains - это не только модели\сущности в БД, но и те сущности, которые не храняться в БД, но мы с ними работаем. Например, сущность Request/Response в разные АПИ, которые имеют разные структуру, поэтому тоже надо создать domain(models) для них. Или же структуру для ошибок.


##  Один микросервис на одну сущность в БД
Обычно каждый микросервис обслуживает один domain(model) в БД. Например, для User отдельный сервер, Product другой, и Order третий.


## Наше приложуха
![alt text](http://i.imgur.com/0KYmHpv.png "diagram")
Клиент будет отправлять запросы на наш микросервис, мы будем делать запросы на GitHub.


## Структура текущего проекта
**client** - тут храним код с помощи, которого можно делать HTTP запросы;
**config** - тут храним конфигурации сервера;
**domain** - тут храним модели/представления нужны нам объектов;
**providers** - тут храним код для работы с внешними сервера. Тут используется код из client;
**services** - отвечат за бизнес логики обработки запроса и возвращения ответа. Тут вызывается код их providers для получения данных и структуры в domain для их сериализации. Также другие utils.
**utils** - сдесь кастомные ошибки для сервера, переиспользуемые функции
