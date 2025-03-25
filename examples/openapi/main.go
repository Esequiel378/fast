package main

import (
	"log"

	"github.com/esequiel378/fast"
)

func main() {
	app, err := fast.New(fast.WithExperimentalOpenAPISchema())
	if err != nil {
		log.Fatal(err)
	}

	// Pet endpoints
	app.Group("/pet").
		MustRegister("", PetHandler{}).
		MustRegister("/findByStatus", PetHandler{}).
		MustRegister("/findByTags", PetHandler{}).
		MustRegister("/{petId}", PetHandler{}).
		MustRegister("/{petId}/uploadImage", PetHandler{})

	// Store endpoints
	app.Group("/store").
		MustRegister("/inventory", StoreHandler{}).
		MustRegister("/order", StoreHandler{}).
		MustRegister("/order/{orderId}", StoreHandler{})

	// User endpoints
	app.Group("/user").
		MustRegister("", UserHandler{}).
		MustRegister("/createWithArray", UserHandler{}).
		MustRegister("/createWithList", UserHandler{}).
		MustRegister("/login", UserHandler{}).
		MustRegister("/logout", UserHandler{}).
		MustRegister("/{username}", UserHandler{})

	log.Fatal(app.Listen(":3003"))
}

type User struct {
	Name string `json:"name"`
}

type UserHandler struct{}

func (h UserHandler) HandleList() fast.Handler {
	type Out struct {
		Users []User `json:"users"`
	}

	return fast.
		Endpoint[fast.In, Out]().
		Handle(func(*fast.Context, fast.In) (Out, error) {
			output := Out{
				Users: []User{
					{Name: "Alice"},
					{Name: "Bob"},
				},
			}

			return output, nil
		})
}

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AccountsHandler struct{}

func (h AccountsHandler) HandleList() fast.Handler {
	type In struct {
		Limit int `json:"limit" validate:"omitempty,gte=10,lte=100"`
		Page  int `json:"page" validate:"omitempty,gte=0"`
	}

	type Out struct {
		Accounts []Account `json:"accounts"`
	}

	return fast.
		Endpoint[In, Out]().
		Path("/:id/test-id").
		Handle(func(_ *fast.Context, input In) (Out, error) {
			output := Out{
				Accounts: []Account{
					{ID: 1, Name: "Alice"},
					{ID: 2, Name: "Bob"},
				},
			}

			return output, nil
		})
}

type PetHandler struct{}

func (h PetHandler) HandleCreate() fast.Handler {
	type Category struct {
		ID   int64  `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	type Pet struct {
		ID        int64     `json:"id,omitempty"`
		Category  *Category `json:"category,omitempty"`
		Name      string    `json:"name"`
		PhotoUrls []string  `json:"photoUrls"`
		Status    string    `json:"status,omitempty"` // available, pending, sold
	}

	return fast.
		Endpoint[Pet, Pet]().
		Handle(func(_ *fast.Context, input Pet) (Pet, error) {
			// In a real implementation, you would save the pet to a database
			// For this example, we'll just return the same pet with an ID
			input.ID = 12345
			return input, nil
		})
}

// StoreHandler handles all store-related endpoints
type StoreHandler struct{}

func (h StoreHandler) HandleInventory() fast.Handler {
	return fast.
		Endpoint[fast.In, map[string]int32]().
		Path("/inventory").
		Handle(func(_ *fast.Context, _ fast.In) (map[string]int32, error) {
			// In a real implementation, you would query the database for inventory
			return map[string]int32{
				"available": 10,
				"pending":   5,
				"sold":      2,
			}, nil
		})
}
