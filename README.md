## Majoo POS (Point Of Sales) <a name = "about"></a>

## Command <a name = "getting_started"></a>

### Application Lifecycle

```
$ cp .env.example .env
$ go mod download
$ go run main.go
```

### Run Apps with Compose

```
docker-compose up -d
```

### Run Integration test

```
// make sure your device already installed docker engine
$ cd integration_test
$ go test ./...
```

## Endpoint <a name = "tests"></a>

| Name    | Endpoint                 | Method   | With Token | Description                 |
| ------- | ------------------------ | -------- | ---------- | --------------------------- |
| Product | _/api/products_          | _POST_   | No         | For add product             |
|         | _/api/products_          | _GET_    | No         | For get products            |
| Cart    | _/api/carts_             | _POST_   | No         | Add product to cart         |
|         | _/api/carts_             | _GET_    | No         | For get products in cart    |
|         | _/api/carts/:product_id_ | _DELETE_ | No         | For delete product in chart |
