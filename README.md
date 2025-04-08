## chirpy app

chirpy is a RESTful api that manages users and their chirps, using go lang http package

for authentication and authorization is used refresh, access token flow

## setup 

you need go (at least 1.24), postgresql, goose, sqlc

> .env
- DB_URL 
    - postgres://user:password@localhost:5432/chirpy
- JWT_SECRET
- POLKA_KEY 
    - api key

#### run
- go run .

#### test
- go test ./...


## API

> /api/chirps
- GET /api/chirps displays all the chirps sorted by creation date in ascending order
- GET /api/chirps/{id} displays chirp by id
- GET /api/chirps?author_id={id} displays chirps by author id
- GET /api/chirps?sort=desc displays all the chirps sorted by creation date in descending order
- POST /api/chirps creates a new chirp for an authorized user
- DELETE /api/chirps/{id} deletes a chirp by id for an authorized user


> /api/users
- POST /api/users creates a new user with provided email and password, the password is hashed before storing
- PUT /api/users updates the users email and password


> /admin/
- GET /admin/metrics returns a HTML with server hits value
- POST /admin/reset resets the database
- GET /admin/users returns all the users

> /api/healthz
- GET /api/healthz returns the status of the service

# auth
- POST /api/login logs a user in, returning user with token and a refresh token in response
- POST /api/refresh refreshes the JWT token provided a refresh token
- POST /api/revoke revokes a refresh token
