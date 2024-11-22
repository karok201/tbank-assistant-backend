# Финансовый помощник (серверная часть)

## Стек

- PHP 8.1
- Swoole 5.1
- Phalcon 5.1
- Elasticsearch 8
- Mysql 5.7

### Пояснение к стеку

### Swoole

Высокопроизводительный программный сервер для PHP с AsyncIO и корутинами (сопрограммами). Он написан на языке C,
благодаря чему достигается высокая эффективность выполнения программ.

По сути, Swoole превращает PHP в реального конкурента языку GO. Благодаря ему PHP становится web-сервером, в отличие от
наиболее популярного пакета FPM.

### Phalcon

Веб-фреймворк с открытым исходным кодом, поставляемый в виде расширения C для языка PHP, обеспечивающего высокую
производительность и меньшее потребление ресурсов.

### Elasticsearch

Поисковой движок, позволяющий быстро выполнять поиск и агрегацию большого количества данных.

Он выбран для хранения информации о тратах и расходах пользователей. Именно этих данных будет больше всего и именно их
необходимо наиболее быстрым путем отображать пользователю.

### Mysql

СУБД с открытым исходным кодом. Одна из самых популярных баз данных.

Выбрана для хранения информации о пользователях и категориях.

