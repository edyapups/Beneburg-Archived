# Beneburg

## Dependencies

- Go >1.19
- Docker
- Docker Compose

## Setup
```bash
$ git clone git@github.com:edyapups/Beneburg.git

$ cd Beneburg

$ echo MYSQL_DATABASE=your_db_name >> .env
$ echo MYSQL_USER=your_db_user >> .env
$ echo MYSQL_PASSWORD=your_db_password >> .env
$ echo MYSQL_HOST=db >> .env
$ echo BOT_TOKEN=your_bot_token >> .env
$ echo SERVER_PORT=your_server_port >> .env

$ docker-compose --env-file .env -f docker-compose.yml up
```