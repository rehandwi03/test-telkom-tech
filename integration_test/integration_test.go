package integrationtest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"interview-telkom-6/handler"
	"interview-telkom-6/repository/persistence"
	"interview-telkom-6/request"
	"interview-telkom-6/response"
	"interview-telkom-6/service"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
)

var (
	db                    *sqlx.DB
	dbName                = "postgres"
	dbHost                = "localhost"
	dbUser                = "postgres"
	dbPass                = "password"
	dbPort                = "5431"
	internalPort          = "5432"
	containerNamePostgres = "postgres-integration-testing"
	r                     *gin.Engine
	productID             uuid.UUID
	qty                   = 1
	qtyAdd                = 3
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       containerNamePostgres,
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_USER=" + dbUser,
			"POSTGRES_PASSWORD=" + dbPass,
			"POSTGRES_DB=" + dbName,
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{internalPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(internalPort): {
				{HostIP: "0.0.0.0", HostPort: dbPort},
			},
		},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	err = resource.Expire(30)
	if err != nil {
		log.Fatalf("could not set expire time: %v", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			dbHost, dbUser, dbPass, dbName, dbPort,
		)
		log.Println(dsn)
		db, err = sqlx.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("error connect to db: %v", err)
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	file, err := ioutil.ReadFile("../order-on-table.sql")
	if err != nil {
		log.Fatalf("error opening file sql: %v", err)
	}

	_, err = db.Exec(string(file))
	if err != nil {
		log.Fatalf("error create table: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", db)
	r = gin.New()
	rGroup := r.Group("/api")

	productRepo := persistence.NewProductRepository(db)
	cartRepo := persistence.NewCartRepository(db)
	cartProductRepo := persistence.NewCartProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	cartSvc := service.NewCartService(ctx, cartRepo, cartProductRepo, productRepo)

	handler.NewProductHandler(rGroup, productSvc)
	handler.NewCartHandler(rGroup, cartSvc)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestInsertProduct(t *testing.T) {
	requestData := request.ProductAddRequest{
		Name:        "Mie Goreng",
		Price:       20000,
		Description: "Mie Goreng enak",
		IsDiscount:  false,
	}
	b, err := json.Marshal(requestData)
	w := httptest.NewRecorder()
	reader := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, "/api/products", reader)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFindProduct(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/products", nil)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	type resStruct struct {
		Data []response.ProductResponse `json:"data"`
	}
	data := resStruct{}

	err = json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	productID = data.Data[0].ID
}

func TestAddProductToCart(t *testing.T) {
	requestData := request.CartAddRequest{
		FullName: "Rehan Dwi",
		Product: request.CartAddProductRequest{
			ProductID: productID,
			Quantity:  qty,
		},
	}

	b, err := json.Marshal(requestData)
	w := httptest.NewRecorder()
	reader := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, "/api/carts", reader)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddProductToCartMultipleQuantity(t *testing.T) {
	requestData := request.CartAddRequest{
		FullName: "Rehan Dwi",
		Product: request.CartAddProductRequest{
			ProductID: productID,
			Quantity:  qtyAdd,
		},
	}

	b, err := json.Marshal(requestData)
	w := httptest.NewRecorder()
	reader := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, "/api/carts", reader)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type resStruct struct {
		Data response.CartResponse `json:"data"`
	}
	data := resStruct{}

	err = json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Equal(t, qty+qtyAdd, data.Data.Products[0].Quantity)
}

func TestFindCart(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/carts?full_name=Rehan Dwi", nil)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type resStruct struct {
		Data response.CartResponse `json:"data"`
	}
	data := resStruct{}

	err = json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Len(t, data.Data.Products, 1)
}

func TestFindCartWithSearch(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/carts?full_name=Rehan Dwi&product_name=Mie", nil)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type resStruct struct {
		Data response.CartResponse `json:"data"`
	}
	data := resStruct{}

	err = json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Len(t, data.Data.Products, 1)
}

func TestDeleteProductFromCart(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/api/carts/"+productID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
