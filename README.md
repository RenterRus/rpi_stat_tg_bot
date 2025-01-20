# rpi_stat_tg_bot

Реализация телеграмм-бота, управляющий RAID массивом.
Функции:
- Монтирование RAID массива
- Вывод справочной информации по массиву
- Скачивание видео с YouTube, VK
- Складывание видео в RAID
- Вывод информации по скачиваемым видео
- Хранение списков/очередей в БД SQLsite3

# Настройка ftp сервера
https://lumpics.ru/how-to-create-ftp-server-in-linux/

# Добавление программы в автозагузку Linux

Написание sh скрипта https://stackoverflow.com/questions/3953886/linux-run-a-binary-in-a-script

Написание скрипта автозапуска https://losst.pro/avtozagruzka-linux


Посмотреть примеры скриптов можно тут в example_run_script и systemd_service_example

# Создание программного RAID-массива на mdadm

https://www.dmosk.ru/miniinstruktions.php?mini=mdadm

https://losst.pro/programmnyj-raid-v-linux




# Общий флоу запуска бота
0. Создать бота в телеграмм
1. Нужно создать RAID-массив, если нет возможности, то кнопка AutoConnect работать не будет
2. Создайте конфиг-файл, в котором укажите:
    - Токен бота в тг
    - Время на реакцию чата. Можно поставить больше 60
    - Свой чат id c ботом. Узнать его можно после запуска, введите любое сообщение в чат, бот напишет ваш с ним чат id
    - FTP пользователя, если используете RAID + sftp
    - Признак имени RAID, например, md
    - Путь, куда складывать скачанные видео
3. Создайте скрипт через который демон будет запускать сервис. Если запуск сводится к простому вызову бинарника и передачи пути к конфигу, то можно пропустить этот этап. Пример: example_run_script
4. Создайте демона в systemd. Обратите внимание, что запуск будет производиться от пользователя. Пример: systemd_service_example
5. Перезапустите systemd через deamon-reload или перезагрузите систему целиком.

Готово (?)

# Стэк
- go1.22.0
- sqlite3
- github.com/go-playground/validator/v10
- github.com/go-telegram-bot-api/telegram-bot-api/v5 
- github.com/lrstanley/go-ytdlp 
- github.com/spf13/viper 

# Скрипт создания таблицы ЫЙДшеу3
create table links (link text unique, status text, name text;

# Установка и настройка sensors

```sudo apt install lm-sensors``` установка
```sudo sensors-detect``` настройка. Везьде Yes

# П.с.

Добавьте файлы cookie от браузера хром с авторизованными сервисами для упрощения скачивания по пути: <пользователь от имени которого происходит запуск сервиса>/.config/google-chrome