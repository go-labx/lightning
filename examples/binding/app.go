package main

import "github.com/go-labx/lightning"

// User contains user information
type User struct {
	FirstName      string     `json:"firstName" validate:"required"`
	LastName       string     `json:"lastName" validate:"required"`
	Age            uint8      `json:"age" validate:"gte=0,lte=130"`
	Email          string     `json:"email" validate:"required,email"`
	FavouriteColor string     `json:"favouriteColor" validate:"iscolor"`           // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `json:"addresses" validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `json:"street" validate:"required"`
	City   string `json:"city" validate:"required"`
	Planet string `json:"planet" validate:"required"`
	Phone  string `json:"phone" validate:"required"`
}

func main() {
	app := lightning.DefaultApp()

	app.Post("/post", func(ctx *lightning.Context) {
		user := &User{}
		err := ctx.Bind(user)
		if err != nil {
			ctx.Fail(-1, err.Error())
			return
		}

		ctx.Success(user)
	})

	app.Run()
}
