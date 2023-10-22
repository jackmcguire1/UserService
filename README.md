# User Service

[![Go Report Card](https://goreportcard.com/badge/github.com/jackmcguire1/UserService)](https://goreportcard.com/report/github.com/jackmcguire1/UserService)
[![codecov](https://codecov.io/gh/jackmcguire1/UserService/graph/badge.svg?token=URT8YBBJFF)](https://codecov.io/gh/jackmcguire1/UserService)

[git]:    https://git-scm.com/
[golang]: https://golang.org/
[modules]: https://github.com/golang/go/wiki/Modules
[goLand]: https://www.jetbrains.com/go/
[golint]: https://github.com/golangci/golangci-lint
[docker]: https://www.docker.com/products/docker-desktop

## ABOUT
> This repo contains a go module that exposes a User Microservice using MongoDB as a datastore

### Prerequisites

- [Git][git]
- [Go 1.21.1][golang]+
- [Docker][docker]


### SETUP
> setup your mongo connection details in docker-compose.yaml
```yaml
environment:
    - MONGO_HOST=mongodb+srv://****
    - MONGO_DATABASE=***
    - MONGO_USERS_COLLECTION=users
  ```

#### run the docker-compose stack
```shell
docker-compose up -d && docker compose watch
```

### Environment Variables
- EVENTS_URL - external HTTP endpoint provided by interested services
- LOG_VERBOSITY - warn | error | info | debug
- MONGO_HOST - your mongo host url
- MONGO_DATABASE - your mongo database
- MONGO_USERS_COLLECTION - your mongo user's collection

## REQUIREMENTS
The service must allow you to:
- add a new User
- modify an existing User
- remove a User
- return a list of the Users, allowing for filtering by certain criteria (e.g. all Users with the
country &quot;GB&quot;)

The service must include:
- A sensible storage mechanism for the Users
- The ability to send events to notify other interested services of changes to User entities

## ASSUMPTIONS
- Service to be accessed via HTTP

- Small user base, no support pagination requests required

- interested services have an exposed HTTP callback endpoint for user events


## IMPROVEMENTS
- Support paginated events for search queries

- Password encryption

- Serve over HTTPS

## ENDPOINTS

### Healthcheck/

<details>
<summary>Healthcheck</summary>

*Healthcheck*
----

* **URL**

  > localhost:7755/healthcheck

* **Method:**
  `GET`

* **Success Response:**
  
  *Code:* 200 <br />
  *Content:*
    ```json
    {
      "logVerbosity": "debug",
      "upTime": "10s"
    }
  ```
</details>


### User/
<details>
<summary>Get a user </summary>

*Get a User*
----

* **URL**

  > localhost:7755/users?id={user-id}

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
        "_id": "100249558",
        "firstName": "Jack",
        "lastName": "McGuire",
        "countryCode": "GB",
        "nickName": "crazyjack12",
        "email": "jack@blah.com",
        "saved": "2021-04-27T17:03:40+01:00"
    }
  ```

OR <br>
   * *Code:* 200 STATUS OK <br />
    *Content:* `{"error": "user not found"}`
    
* **Error Responses:**

  * **Code:** 400 BAD REQUEST error <br />
    **Content:** `{"error":"reason"}`
    
    OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{"error":"reason"}`

* **Notes:**

 a response of `{"error": "user not found"}` will be returned if user cannot be found
 
</details>

<details>
<summary>Delete a user </summary>

*Delete a User*
----

* **URL**

  > localhost:7755/users?id={user-id}

* **Method:**
  `DELETE`
  
*  **URL Params**
   **Required:**
 
   id=[string]

* **Success Response:**
  
  *Code:* 200 <br />
  *Content:*
    ```json
    {
        "deleted": true,
        "message": "success"
    }
  ```
    
* **Error Responses:**

  *  *Code:* 200 <br />
      *Content:*
        ```json
        {"error":"reason"}
      ```
    
    OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `error reason`

</details>

<details>
<summary>Update a user </summary>

*Update a User*
----

* **URL**

  > localhost:7755/users

* **Method:**
  `POST`
  
  * **Data Params**
     **Required:**
 
     ```
        {
            "_id": "100249558",
            "firstName": "Jack",
            "lastName": "McGuire",
            "countryCode": "GB",
            "email": "jack@blah.com",
            "nickName": "skr",
        }
      ```
     **OPTIONAL:**
      ```
        {
            "nickName": "crazyjack12",
        }
      ```

* **Success Response:**
  
  *Code:* 200 STATUS OK<br />
  *Content:*
   ```json
    {
      "_id": "100249558",
      "firstName": "Jack",
      "lastName": "McGuire",
      "countryCode": "GB",
      "email": "jack@blah.com",
      "nickName": "skr"
    }
  ```
    
* **Error Responses:**

  *  *Code:* 400 BAD REQUEST <br />
      *Content:* `{"error":"reason"}`
    
    OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{"error":"reason"}`

* **Notes:**

> emails must contain '@'

> passwords must be more than 5 chars long

> country code must be ISO ALPHA-2
</details>


<details>
<summary>Create a user </summary>

*Create a User*
----

* **URL**

  > localhost:7755/users

* **Method:**
  `PUT`
  
* **Data Params**
   **Required:**
 
   ```json
   {
      "firstName": "Jack",
      "lastName": "McGuire",
      "countryCode": "GB",
      "email": "GB"
  }
    ```
  
  **OPTIONAL:**
  ```
    {
        "_id": "100249558",
        blah,
    }
  ```

* **Success Response:**
  
  *Code:* 200 STATUS OK<br />
  *Content:*
    ```json
    {
        "_id": "100249558",
        "firstName": "Jack",
        "lastName": "McGuire",
        "countryCode": "GB",
        "nickName": "crazyjack12",
        "email": "jack@blah.com",
        "saved": "2021-04-27T17:03:40+01:00"
    }
  ```
    
* **Error Responses:**

  *  *Code:* 400 BAD REQUEST <br />
      *Content:* `error reason`
    
    OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `error reason`

* **Notes:**

> the field 'ID' is optional

> emails must contain '@'

> passwords must be more than 5 chars long

> country code must be ISO ALPHA-2
</details>

### Search/Users/

<details>
<summary>All Users</summary>

*All Users*
----

* **URL**

  > localhost:7755/search/users/

* **Method:**
  `GET`

* **Success Response:**

  *Code:* 200 <br />
  *Content:*
    ```json
    {
      "users": [
    	{
            "_id": "100249558",
            "firstName": "Jack",
            "lastName": "McGuire",
            "countryCode": "GB",
            "nickName": "crazyjack12",
            "email": "jack@blah.com",
            "saved": "2021-04-27T17:03:40+01:00"
    	}
      ]
    }
  ```

OR <br>
* *Code:* 200 STATUS OK <br />
  *Content:*
  ```
  {
  "users": []
  }
  ```

* **Error Responses:**

    * **Code:** 400 BAD REQUEST error <br />
      **Content:** `{"error":"reason"}`

      OR

    * **Code:** 500 INTERNAL SERVER ERROR <br />
      **Content:** `{"error":"reason"}`

* **Notes:**

'cc' query parameter value will automatically be defaulted into uppercase

</details>

<details>
<summary>By Country</summary>

*By Country*
----

* **URL**

  > localhost:7755/search/users/by_country?cc={country-code}

* **Method:**
  `GET`
  
*  **URL Params**
   **Required:**
 
   cc=[string]

* **Success Response:**
  
  *Code:* 200 <br />
  *Content:*
    ```json
    {
      "users": [
    	{
            "_id": "100249558",
            "firstName": "Jack",
            "lastName": "McGuire",
            "countryCode": "GB",
            "nickName": "crazyjack12",
            "email": "jack@blah.com",
            "saved": "2021-04-27T17:03:40+01:00"
    	}
      ]
    }
  ```

OR <br>
   * *Code:* 200 STATUS OK <br />
    *Content:*
    ```
    {
        "users": []
    }
    ```

* **Error Responses:**

  * **Code:** 400 BAD REQUEST error <br />
    **Content:** `{"error":"reason"}`
    
    OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{"error":"reason"}`

* **Notes:**

'cc' query parameter value will automatically be defaulted into uppercase
 
</details>

## Thanks

This project exists thanks to **ME!**.

## Donations
All donations are appreciated!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](http://paypal.me/crazyjack12)
