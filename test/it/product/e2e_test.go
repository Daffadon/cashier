package product_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"tiga-putra-cashier-be/cmd"
	"tiga-putra-cashier-be/database"
	"tiga-putra-cashier-be/di"
	"tiga-putra-cashier-be/entity"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type e2eProductTestSuite struct {
	suite.Suite
	dbConn *gorm.DB
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eProductTestSuite{})
}

func (e *e2eProductTestSuite) SetupSuite() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPass, dbName, dbPort)
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	e.NoError(err)

	e.dbConn = db

	err = database.MigrateUp(e.dbConn)
	e.Require().NoError(err)

	container := di.BuildContainer()
	serverReady := make(chan bool)
	server := &cmd.Server{
		Container:   container,
		ServerReady: serverReady,
	}

	go server.Start()
	<-serverReady
}

func (e *e2eProductTestSuite) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	_ = p.Signal(syscall.SIGINT)
}

func (e *e2eProductTestSuite) SetupTest() {
	err := database.MigrateUp(e.dbConn)
	e.Require().NoError(err)
}

func (e *e2eProductTestSuite) TearDownTest() {
	e.NoError(database.MigrateDown(e.dbConn))
}

func (e *e2eProductTestSuite) Test_E2EProduct_GetProduct() {
	product := &entity.Product{
		BarcodeId:   "1",
		Image:       "img-1",
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	e.NoError(e.dbConn.Create(product).Error)
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/product?page=1", nil)
	e.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)
	e.NoError(err)

	e.Equal(http.StatusOK, response.StatusCode)
	byteBody, err := io.ReadAll(response.Body)
	e.NoError(err)
	e.Equal(`{"status_code":200,"message":"Success Get All product","data":{"products":[{"barcode_id":"1","image":"img-1","title":"title-1","price":"1000","description":"desc-1"}],"page_meta_data":{"page":1,"prev_page":1,"next_page":1,"total_page":1}}}`, strings.Trim(string(byteBody), "\n"))
	response.Body.Close()
}
func (e *e2eProductTestSuite) Test_E2EProduct_GetProductDetail() {
	product := &entity.Product{
		BarcodeId:   "1",
		Image:       "img-1",
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	e.NoError(e.dbConn.Create(product).Error)
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/product/1", nil)
	e.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)
	e.NoError(err)

	e.Equal(http.StatusOK, response.StatusCode)
	byteBody, err := io.ReadAll(response.Body)
	e.NoError(err)
	e.Equal(`{"status_code":200,"message":"Success Get Product Detail","data":{"barcode_id":"1","image":"img-1","title":"title-1","price":"1000","description":"desc-1"}}`, strings.Trim(string(byteBody), "\n"))
	response.Body.Close()
}

func (e *e2eProductTestSuite) Test_E2EProduct_SearchProduct() {
	product := &entity.Product{
		BarcodeId:   "1",
		Image:       "img-1",
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	e.NoError(e.dbConn.Create(product).Error)
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/product/search?barcode_id=1", nil)
	e.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)
	e.NoError(err)

	e.Equal(http.StatusOK, response.StatusCode)
	byteBody, err := io.ReadAll(response.Body)
	e.NoError(err)
	e.Equal(`{"status_code":200,"message":"Success Get All product","data":[{"barcode_id":"1","image":"img-1","title":"title-1","price":"1000","description":"desc-1"}]}`, strings.Trim(string(byteBody), "\n"))
	response.Body.Close()
}

func (e *e2eProductTestSuite) Test_E2EProduct_AddProduct() {
	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("barcode_id", "1")
	_ = formWriter.WriteField("title", "title-1")
	_ = formWriter.WriteField("price", "1000")
	_ = formWriter.WriteField("description", "desc-1")

	fileWriter, _ := formWriter.CreateFormFile("image", "test.jpg")
	_, err := fileWriter.Write([]byte("fake image data"))
	e.NoError(err)
	formWriter.Close()

	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/product", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())
	e.NoError(err)

	client := http.Client{}
	response, err := client.Do(request)
	e.NoError(err)

	e.Equal(http.StatusOK, response.StatusCode)
	byteBody, err := io.ReadAll(response.Body)
	e.NoError(err)
	e.Equal(`{"status_code":200,"message":"Success Add Product","data":null}`, strings.Trim(string(byteBody), "\n"))
	response.Body.Close()
}

func (e *e2eProductTestSuite) Test_E2EProduct_UpdateProduct() {
	product := &entity.Product{
		BarcodeId:   "1",
		Image:       "img-1",
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	e.NoError(e.dbConn.Create(product).Error)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("title", "title-Update")
	err := formWriter.Close()
	e.NoError(err)

	request, err := http.NewRequest(http.MethodPatch, "http://localhost:8080/v1/product/1", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())
	e.NoError(err)

	client := http.Client{}
	response, err := client.Do(request)
	e.NoError(err)

	e.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	e.NoError(err)

	e.Equal(`{"status_code":200,"message":"Success Update Product","data":null}`, strings.Trim(string(byteBody), "\n"))
	response.Body.Close()
}

func (e *e2eProductTestSuite) Test_E2EProduct_DeleteProduct() {
	product := &entity.Product{
		BarcodeId:   "1",
		Image:       "img-1",
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	e.NoError(e.dbConn.Create(product).Error)

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/v1/product/1", nil)
	e.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)
	e.NoError(err)

	e.Equal(http.StatusOK, response.StatusCode)
	byteBody, err := io.ReadAll(response.Body)
	e.NoError(err)
	e.Equal(`{"status_code":200,"message":"Success Delete Product","data":null}`, strings.Trim(string(byteBody), "\n"))
	response.Body.Close()
}
