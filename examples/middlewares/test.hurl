GET http://localhost:3000/greeting
HTTP 401

GET http://localhost:3000/greeting?name=Esequiel
API-Key: fast-is-awesome
HTTP 200
[Asserts]
jsonpath "$.message" == "Hello, Esequiel!"
