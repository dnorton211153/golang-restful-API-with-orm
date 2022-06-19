# Golang RESTful API w/ mux/gorm/MySQL

```
This is a simple example of a Golang RESTful API,
using the gorm ORM library and the gorilla/mux router.

Tested with Postman (http://localhost:8000/user route)

My first attempt using Golang.  Hello Daz Productions guys!
```

Using Repository Interface per Model (best DB migration practice),
thanks to suggestions from [Praveen](https://techinscribed.com/different-approaches-to-pass-database-connection-into-controllers-in-golang/).

## Creating a user

Use Postman or curl to create a user:
```
curl -X POST -H "Content-Type: application/json" -d '{"username":"Dave","email":"dave@gmail.com","password":"asdf"}' http://localhost:3000/user
```

Assuming mySql is running on localhost:3306 (or wherever you configured in .env), the user table will be created on startup, and the user looks something like this:

```
{
    "id": 1,
    "username": "Dave",
    "email": "dave@gmail.com",
    "password": "asdf"
}
```

## About the authentication model

Sessions are tracked in memory by uuid, expiry, and username.
Upon successful authentication (on the /login route), a session 
is created and stored in memory, and a cookie with session_token
is sent to the client.  Without active session, a 401 is returned.

The only route that currently depends on the session_token cookie 
is the /user (GET) route (GetAllUsers).

To require authentication, add *validateSession* middleware to the route.  E.g.,

```
router.HandleFunc("/user", validateSession(http.HandlerFunc(userHandler.GetAll))).Methods("GET")
```

The login route currently uses cleartext for authentication;
just POST something like this in the body to the /login route:

```
{
	"username": "Dave",
	"password": "asdf"
}
```

** @author Norton 2022