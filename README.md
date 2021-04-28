# User Service

[![Build Status](https://travis-ci.org/jackmcguire1/UserService.svg?branch=main)](hhttps://travis-ci.org/jackmcguire1/UserService)
[![Go Report Card](https://goreportcard.com/badge/github.com/jackmcguire1/UserService)](https://goreportcard.com/report/github.com/jackmcguire1/UserService)
[![codecov](https://codecov.io/gh/jackmcguire1/UserService/branch/main/graph/badge.svg?token=URT8YBBJFF)](https://codecov.io/gh/jackmcguire1/UserService)

[git]:    https://git-scm.com/
[golang]: https://golang.org/
[modules]: https://github.com/golang/go/wiki/Modules
[goLand]: https://www.jetbrains.com/go/
[golint]: https://github.com/golangci/golangci-lint
[docker]: https://www.docker.com/products/docker-desktop

## ABOUT
> This repo contains a go module that exposes a User Microservice

### Prerequisites

- [Git][git]
- [Go 1.16][golang]+
- [Docker][docker]


### SETUP

> build the userservice docker container
```shell
docker build -t userservice .
```

> run the stack
```shell
docker-compose up -d
```

## ENDPOINTS

### User/
<details>
<summary> Get a user </summary>

*Get a User*
----

* **URL**

  > localhost:7755/user?id={user-id}

* **Method:**
  `GET`
  
*  **URL Params**
   **Required:**
 
   id=[string]

* **Success Response:**
  
  *Code:* 200 <br />
  *Content:*
    ```json
    {
        "ID": "steve",
        "FirstName": "Jack",
        "LastName": "McGuire",
        "CountryCode": "GB",
        "Saved": "2021-04-27T17:03:40+01:00"
    }
  ```

OR <br>
   * *Code:* 200 STATUS OK <br />
    *Content:* `{}`
    
* **Error Responses:**

  * **Code:** 400 BAD REQUEST error <br />
    **Content:** `error reason`
    
    OR
    
  * **Code:** 500 internal error <br />
    **Content:** `error reason`

* **Notes:**

 an empty response of `{}` will be returned if user cannot be found
 
</details>

#### DELETE 
> DELETE localhost:7755/user?id={user-id}

This endpoint will delete a user from the system

#### POST
> POST localhost:7755/user

This endpoint updates an existing user<br>
Note:- The updated user object will be returned

JSON request BODY
```json
{
    "ID": "steve",
    "FirstName": "Jack",
    "LastName": "McGuire",
    "CountryCode": "GB"
}
```

#### PUT
> PUT localhost:7755/user
 
This endpoint creates a new user<br><br>
note:- 'ID' field is optional, if not provided a userID will be auto-generated

JSON request BODY
```json
{
    "ID": "100249558", -- ID FIELD OPTIONAL
    "FirstName": "Jack",
    "LastName": "McGuire",
    "CountryCode": "GB"
}
```

## Thanks

This project exists thanks to **ME!**.

## Donations
All donations are appreciated!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](http://paypal.me/crazyjack12)