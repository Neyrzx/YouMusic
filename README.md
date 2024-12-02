## <img align="left" hspace="20" src="https://github.com/nestoris/Win98SE/blob/main/SE98/mimes/64/audio-x-mod.png?raw=true" width="64" alt="YouMusic"/> YouMusic
📚 Онлайн Библиотека Песен 🎶

Привет! 👋 Вы находитесь в репозитории YouMusic — проекте, разработанном в рамках тестового задания для EffectiveMobile. Цель проекта — реализовать онлайн библиотеку песен с возможностью получения, добавления, изменения и удаления треков. 

### 🚀 Установка и запуск
1. Клонируйте репозиторий:

```bash
git clone https://github.com/Neyrzx/YouMusic.git
cd YouMusic
```

> [!IMPORTANT]
> 2. Создайте .env файл и настройте под себя если необходимо:
> ```bash
> cp .env.example .env
> ```

3. Установите зависимости и инструменты, выполните:

    > Все записимости проекта: **swag**, **golangci-lint** и др. будут установлены в папку `bin/` в каталоге репозитория.

```bash
make init
```

4. Запустите проект:
```bash
make compose-dev
```

## 🔍 Swagger документация
Swagger спецификация будет доступна по адресу: http://localhost:9090/docs/index.html (после запуска сервиса).

## 🛠 Makefile команды
* `make install` - Установить все необходимые инструменты.
* `make lint` - Проверить код на соответствие стандартам.
* `make test` - Запустить тесты.
* `make migration-up` - Применить миграции.
* `make migration-down` - Откатить миграции.
* `make compose-down-clean` - Остановка контейнеров с флагом -v.
* и др. [Makefile](./Makefile)

## 🎉 Примененные технологии
* **Go** - основной язык для реализации сервиса.
* **PostgreSQL** - для хранения данных.
* **Swagger** - для документирования API.
* **Docker** - для контейнеризации приложения.

Спасибо, что заглянули! Наслаждайтесь использованием YouMusic! 🎵

#### **Телеграм**: [@neyrzx](https://t.me/neyrzx)
