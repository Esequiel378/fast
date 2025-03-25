GET http://localhost:3000/admin/users
HTTP 200

[Asserts]
jsonpath "$.users" count == 2

jsonpath "$.users[0].name" == "Alice"
jsonpath "$.users[1].name" == "Bob"

GET http://localhost:3000/admin/accounts
HTTP 200

[Asserts]
jsonpath "$.accounts" count == 2

jsonpath "$.accounts[0].id" == 1
jsonpath "$.accounts[0].name" == "Alice"

jsonpath "$.accounts[1].id" == 2
jsonpath "$.accounts[1].name" == "Bob"
