## Introduction

🚀🚀🚀 lightning is a lightweight and fast web framework for Go. It is designed to be easy to use and highly performant.

## Features

- Easy to use and quick to get started
- Supports middleware
- Fast routing, with routing algorithm implemented based on Trie tree
- Support for grouping routes and applying middleware to specific groups
- Customizable 404 Not Found and 500 Internal Server Error handler functions

## Getting Started

To get started with lightning, simply install it using go get:

```
go get github.com/go-labx/lightning
```

Then, create a new lightning app and start adding routes:

```go
package main

import (
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, lightning.Map{
			"message": "pong",
		})
	})

	app.Run()
}
```

To run the lightning app, run the following command:

```bash
go run app.go
```

To verify that the server has started successfully, run the following command in your terminal:

```bash
curl http://127.0.0.1:6789/ping
```


## Documentation

For more information on how to use lightning, check out the [documentation](https://go-labx.github.io/docs/introduction).

## Contributing

If you'd like to contribute to lightning, please
see [CONTRIBUTING.md](https://github.com/go-labx/lightning/blob/main/CONTRIBUTING.md) for guidelines.

## License

lightning is licensed under the [MIT License](https://github.com/go-labx/lightning/blob/main/LICENSE).


