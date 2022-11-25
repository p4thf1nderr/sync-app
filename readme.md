- Запуск приложениия:

    ```make start```. 
    Запустит синхронизацию директорий folders/origin и folders/copy.
    Для запуска с кастомными директориями нужно их указать вторым и третьим аргументом.

    ```go run cmd/syncd/main.go folders/origin folders/copy```

- Тесты + coverage: ```make test```

- Бенчмарки: ```make bench```