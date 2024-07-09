package fast

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Context struct {
	ctx *fiber.Ctx
}

func newContext(ctx *fiber.Ctx) Context {
	return Context{
		ctx: ctx,
	}
}

// GetReqHeaders returns the HTTP request headers.
func (c Context) Headers() http.Header {
	return http.Header(c.ctx.GetReqHeaders())
}

// Query returns the query string parameter in the url.
// Defaults to empty string "" if the query doesn't exist.
// If a default value is given, it will return that value if the query doesn't exist.
//
// WARN: Returned value is only valid within the handler. Do not store any references. Make copies of the values instead.
func (c Context) Query(key string, defaultValue ...string) string {
	return c.ctx.Query(key, defaultValue...)
}

// Queries returns a map of query parameters and their values.
//
// GET /?name=alex&wanna_cake=2&id=
// Queries()["name"] == "alex"
// Queries()["wanna_cake"] == "2"
// Queries()["id"] == ""
//
// GET /?field1=value1&field1=value2&field2=value3
// Queries()["field1"] == "value2"
// Queries()["field2"] == "value3"
//
// GET /?list_a=1&list_a=2&list_a=3&list_b[]=1&list_b[]=2&list_b[]=3&list_c=1,2,3
// Queries()["list_a"] == "3"
// Queries()["list_b[]"] == "3"
// Queries()["list_c"] == "1,2,3"
//
// GET /api/search?filters.author.name=John&filters.category.name=Technology&filters[customer][name]=Alice&filters[status]=pending
// Queries()["filters.author.name"] == "John"
// Queries()["filters.category.name"] == "Technology"
// Queries()["filters[customer][name]"] == "Alice"
// Queries()["filters[status]"] == "pending"
func (c Context) Queries() map[string]string {
	return c.ctx.Queries()
}

// QueryInt returns integer value of key string parameter in the url.
// Default to empty or invalid key is 0.
//
//	GET /?name=alex&wanna_cake=2&id=
//	QueryInt("wanna_cake", 1) == 2
//	QueryInt("name", 1) == 1
//	QueryInt("id", 1) == 1
//	QueryInt("id") == 0
func (c Context) QueryInt(key string, defaultValue ...int) int {
	return c.ctx.QueryInt(key, defaultValue...)
}

// QueryBool returns bool value of key string parameter in the url.
// Default to empty or invalid key is false.
//
//	Get /?name=alex&want_pizza=false&id=
//	QueryBool("want_pizza") == false
//	QueryBool("want_pizza", true) == false
//	QueryBool("name") == false
//	QueryBool("name", true) == true
//	QueryBool("id") == false
//	QueryBool("id", true) == true
func (c Context) QueryBool(key string, defaultValue ...bool) bool {
	return c.ctx.QueryBool(key, defaultValue...)
}

// QueryFloat returns float64 value of key string parameter in the url.
// Default to empty or invalid key is 0.
//
//	GET /?name=alex&amount=32.23&id=
//	QueryFloat("amount") = 32.23
//	QueryFloat("amount", 3) = 32.23
//	QueryFloat("name", 1) = 1
//	QueryFloat("name") = 0
//	QueryFloat("id", 3) = 3
func (c Context) QueryFloat(key string, defaultValue ...float64) float64 {
	return c.ctx.QueryFloat(key, defaultValue...)
}

// Params is used to get the route parameters.
// Defaults to empty string "" if the param doesn't exist.
// If a default value is given, it will return that value if the param doesn't exist.
// WARN: Returned value is only valid within the handler. Do not store any references. Make copies of the values instead.
func (c Context) Params(key string, defaultValue ...string) string {
	return c.ctx.Params(key)
}

// ParamsInt is used to get an integer from the route parameters
// it defaults to zero if the parameter is not found or if the
// parameter cannot be converted to an integer
// If a default value is given, it will return that value in case the param
// doesn't exist or cannot be converted to an integer
func (c Context) ParamsInt(key string, defaultValue ...int) (int, error) {
	return c.ctx.ParamsInt(key, defaultValue...)
}
