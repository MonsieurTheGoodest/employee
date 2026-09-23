#!/bin/bash
echo "разработка: "
curl 'http://127.0.0.1:9000/department/разработка'
echo "продажи: "
curl 'http://127.0.0.1:9000/department/продажи'
echo "тестирование: "
curl 'http://127.0.0.1:9000/department/тестирование'

echo создание Кроша Крошовича из отдела разработки
curl -X POST -d '{"first_name":"Крош","last_name":"Крошович","department":"разработка"}' http:/127.0.0.1:9000/employee

echo создание Кроша Крошовича из отдела тестирования
curl -X POST -d '{"first_name":"Крош","last_name":"Крошович","department":"тестирование"}' http:/127.0.0.1:9000/employee

curl http://127.0.0.1:9000/pending

sleep 30

echo создание Ёжика Ёжиковича из отдела разработки
curl -X POST -d '{"first_name":"Ёжик","last_name":"Ёжикович","department":"разработка"}' http:/127.0.0.1:9000/employee

sleep 60

echo "разработка: "
curl 'http://127.0.0.1:9000/department/разработка'
echo "продажи: "
curl 'http://127.0.0.1:9000/department/продажи'
echo "тестирование: "
curl 'http://127.0.0.1:9000/department/тестирование'

echo создание Нюши Нюшовны из отдела продаж
curl -X POST -d '{"first_name":"Нюша","last_name":"Нюшовна","department":"продажи"}' http:/127.0.0.1:9000/employee

sleep 30

curl http://127.0.0.1:9000/pending

echo "разработка: "
curl 'http://127.0.0.1:9000/department/разработка'
echo "продажи: "
curl 'http://127.0.0.1:9000/department/продажи'
echo "тестирование: "
curl 'http://127.0.0.1:9000/department/тестирование'

echo создание Бараша Барашовича из отдела тестирования
curl -X POST -d '{"first_name":"Бараш","last_name":"Барашович","department":"тестирование"}' http:/127.0.0.1:9000/employee

sleep 60

curl http://127.0.0.1:9000/pending

echo "разработка: "
curl 'http://127.0.0.1:9000/department/разработка'
echo "продажи: "
curl 'http://127.0.0.1:9000/department/продажи'
echo "тестирование: "
curl 'http://127.0.0.1:9000/department/тестирование'

echo создание Лосяш Лосяшовича из отдела разработки
curl -X POST -d '{"first_name":"Лосяш","last_name":"Лосяшович","department":"разработка"}' http:/127.0.0.1:9000/employee

sleep 60

curl http://127.0.0.1:9000/pending

echo "разработка: "
curl 'http://127.0.0.1:9000/department/разработка'
echo "продажи: "
curl 'http://127.0.0.1:9000/department/продажи'
echo "тестирование: "
curl 'http://127.0.0.1:9000/department/тестирование'

echo создание Совуньи Птицьевны из отдела тестирования
curl -X POST -d '{"first_name":"Совунья","last_name":"Птицьевна","department":"тестирование"}' http:/127.0.0.1:9000/employee

sleep 30

echo создание Кар Карыча из отдела тестирования
curl -X POST -d '{"first_name":"Кар","last_name":"Карыч","department":"тестирование"}' http:/127.0.0.1:9000/employee

sleep 30

echo создание Копатыча Копатовича из отдела продаж
curl -X POST -d '{"first_name":"Копатыч","last_name":"Копатович","department":"продажи"}' http:/127.0.0.1:9000/employee

sleep 60

curl http://127.0.0.1:9000/pending

echo "разработка: "
curl 'http://127.0.0.1:9000/department/разработка'
echo "продажи: "
curl 'http://127.0.0.1:9000/department/продажи'
echo "тестирование: "
curl 'http://127.0.0.1:9000/department/тестирование'

echo создание Пина Пиновича из отдела продаж
curl -X POST -d '{"first_name":"Пин","last_name":"Пинович","department":"продажи"}' http:/127.0.0.1:9000/employee

sleep 45

echo создание Би Би из отдела разработки
curl -X POST -d '{"first_name":"Би","last_name":"Би","department":"разработка"}' http:/127.0.0.1:9000/employee

sleep 190

curl http://127.0.0.1:9000/pending

echo "разработка: "
curl 'http://127.0.0.1:9000/department/разработка'
echo "продажи: "
curl 'http://127.0.0.1:9000/department/продажи'
echo "тестирование: "
curl 'http://127.0.0.1:9000/department/тестирование'