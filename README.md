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


### Environment Variables
- EVENTS_URL - external HTTP endpoint provided by interested services
- LOG_VERBOSITY - warn | error | info | debug


## REQUIREMENTS
The service must allow you to:
- add a new User
- modify an existing User
- remove a User
- return a list of the Users, allowing for filtering by certain criteria (e.g. all Users with the
country &quot;UK&quot;)

The service must include:
- A sensible storage mechanism for the Users
- The ability to send events to notify other interested services of changes to User entities

## ASSUMPTIONS
- Service to be accessed via HTTP

- Small user base, no support pagination requests required

- Passwords do not have to be encrypted

- interested services have an exposed HTTP callback endpoint for user events


## IMPROVEMENTS
- Support paginated events for elastic search queries

- Password encryption

- Serve over HTTPS

## ENDPOINTS

### Healthcheck/

<details>
<summary>Healthcheck</summary>

*By Country*
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
      "LogVerbosity": "info"
    }
  ```
</details>


### User/
<details>
<summary>Get a user </summary>

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
        "ID": "100249558",
        "FirstName": "Jack",
        "LastName": "McGuire",
        "CountryCode": "GB",
        "NickName": "crazyjack12",
        "Email": "jack@blah.com",
        "Password": "blah",
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
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `error reason`

* **Notes:**

 an empty response of `{}` will be returned if user cannot be found
 
</details>

<details>
<summary>Delete a user </summary>

*Delete a User*
----

* **URL**

  > localhost:7755/user?id={user-id}

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
        "Delete": true,
        "Message": "success"
    }
  ```
    
* **Error Responses:**

  *  *Code:* 200 <br />
      *Content:*
        ```json
        {
            "Delete": false,
            "Message": "error info"
        }
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

  > localhost:7755/user

* **Method:**
  `POST`
  
* **Data Params**
   **Required:**
 
   ```
      {
          "ID": "100249558",
          "FirstName": "Jack",
          "LastName": "McGuire",
          "CountryCode": "GB",
          "Email": "jack@blah.com",
          "Password": "blah1",
      }
    ```
   **OPTIONAL:**
    ```
      {
          "NickName": "crazyjack12",
      }
    ```

* **Success Response:**
  
  *Code:* 200 STATUS OK<br />
  *Content:*
   ```json
    {
        "ID": "100249558",
        "FirstName": "Jack",
        "LastName": "McGuire",
        "CountryCode": "GB",
        "NickName": "crazyjack12",
        "Email": "GB",
        "Password": "BLAH1",
        "Saved": "2021-04-27T17:03:40+01:00"
    }
  ```
    
* **Error Responses:**

  *  *Code:* 400 BAD REQUEST <br />
      *Content:* `error reason`
    
    OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `error reason`

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

  > localhost:7755/user

* **Method:**
  `PUT`
  
* **Data Params**
   **Required:**
 
   ```
   {
           "FirstName": "Jack",
           "LastName": "McGuire",
           "CountryCode": "GB",
           "Email": "GB",
           "Password": "GB",
       }
    ```
  
  **OPTIONAL:**
  ```
    {
        "ID": "100249558",
        "NickName": "100249558",
    }
  ```

* **Success Response:**
  
  *Code:* 200 STATUS OK<br />
  *Content:*
    ```json
    {
        "ID": "100249558",
        "FirstName": "Jack",
        "LastName": "McGuire",
        "CountryCode": "GB",
        "NickName": "crazyjack12",
        "Email": "jack@blah.com",
        "Password": "blah1",
        "Saved": "2021-04-27T17:03:40+01:00"
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
      "Users": [
    	{
            "ID": "100249558",
            "FirstName": "Jack",
            "LastName": "McGuire",
            "CountryCode": "GB",
            "NickName": "crazyjack12",
            "Email": "jack@blah.com",
            "Password": "blah1",
            "Saved": "2021-04-27T17:03:40+01:00"
    	}
      ]
    }
  ```

OR <br>
   * *Code:* 200 STATUS OK <br />
    *Content:*
    ```
    {
        "Users": []
    }
    ```

* **Error Responses:**

  * **Code:** 400 BAD REQUEST error <br />
    **Content:** `error reason`
    
    OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `error reason`

* **Notes:**

'cc' query parameter will auto be defaulted into uppercase
 
</details>

## Thanks

This project exists thanks to **ME!**.

## Donations
All donations are appreciated!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](http://paypal.me/crazyjack12)