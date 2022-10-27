# Beneburg

## Dependencies

- Go >1.19
- Docker
- Docker Compose

## Setup
```bash
$ git clone git@github.com:edyapups/Beneburg.git

$ cd Beneburg

$ echo MYSQL_DATABASE=your_db_name >> .env.local
$ echo MYSQL_USER=your_db_user >> .env.local
$ echo MYSQL_PASSWORD=your_db_password >> .env.local
$ echo MYSQL_HOST=db >> .env.local
$ echo BOT_TOKEN=your_bot_token >> .env.local
$ echo SERVER_PORT=your_server_port >> .env.local

$ docker-compose --env-file .env.local -f docker-compose.yml up
```