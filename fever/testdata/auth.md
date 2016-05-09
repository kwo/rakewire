# Auth

Login as the test user `jeff` with the password `abcdefg`

## POST /?api
* Content-Type: application/x-www-form-urlencoded

```
api_key=edbabac8c85235bdfe6a75ae6079aa4f
```

===

* Status: 200
* Content-Type: text/json; charset=utf-8
* Data.api_version: 3
* Data.auth: 1
* Data.last_refreshed_on_time: /\d{9}/


# Auth Fail

## POST /?api
* Content-Type: application/x-www-form-urlencoded

```
api_key=bogusmd5signature
```

===

* Status: 200
* Content-Type: text/json; charset=utf-8
* Data.api_version: 3
* Data.auth: 0
