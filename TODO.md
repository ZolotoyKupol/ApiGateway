# Делаем api-gateway для гостиницы

1. Напиши hello world и запусти
2. 
    Добавить ручки: (пример https://go.dev/doc/tutorial/web-service-gin)
    (будем исопльзовать gin)
    В main.go без БД
    1. Создать гостя `POST` `/guest` body{"name", ...} reponse 203, {"id"}
    2. Обновить `GET` `/guest/:id` body{"name", ...} reponse 200, {"id"} 
    3. Обновить `PUT` `/guest` body{"name", ...} reponse 200, {"id"} 
    4. Удалить `DELETE` `/guest/:id` reponse 200

3. В постмане отправить запросы и проверить

4. отличие слайса от массива, структура слайса, что происходит при append
5. маппа, бакеты, миграции, что проиходит при коллизиях