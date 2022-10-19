# Beneburg

## Dependencies

- Go >1.19
- Docker
- Docker Compose

## Setup
```$ git clone git@github.com:edyapups/Beneburg.git```

```$ cd Beneburg```

```$ echo MYSQL_DATABASE=your_db_name >> .env```

```$ echo MYSQL_USER=your_db_user >> .env```

```$ echo MYSQL_PASSWORD=your_db_password >> .env```

```$ docker-compose --env-file .env -f docker-compose.yml up```