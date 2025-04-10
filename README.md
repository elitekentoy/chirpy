# General Note

- This is a Guided Project from Boot.Dev
- The code is not written so clean because of skill issue.
- Documentation is not yet complete and will continue sometime in the future


# Go Installation
Executing the following in the Terminal:
```
curl -sS https://webi.sh/golang | sh
```

# PostgreSQL Insllation
1. Execute the following in the Terminal
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```
2. (For Linux) Set the password:
```
sudo passwd postgres
```
3. After installation, start the Postgres Server in the background:
```
sudo service postgresql start
```
4. Login into the PostgreSQL shell:
```
sudo -u postgres psql
```
5. Create Database called Chirpy
```
CREATE DATABASE chirpy
```
6. Switch to chirpy database:
```
\c chirpy
```
7. (For Linux) update the user password
```
ALTER USER postgres PASSWORD 'postgres'
```

# Library Installation

1. Download the necessary 3rd party packages needed for this project by executing in the terminal:
```
go get -u github.com/go-http-utils/headers
go get -u github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
go get github.com/joho/godotenv
go get github.com/lib/pq
```