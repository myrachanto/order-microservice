package routes

import (
	"fmt"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/microservice/order/src/controllers"
	pubsub "github.com/myrachanto/microservice/order/src/events"
	m "github.com/myrachanto/microservice/order/src/middlewares"
	"github.com/myrachanto/microservice/order/src/repository"
	service "github.com/myrachanto/microservice/order/src/services"

	"github.com/spf13/viper"
)

//StoreAPI =>entry point to routes
type Open struct {
	Port     string `mapstructure:"PORT"`
	Key      string `mapstructure:"EncryptionKey"`
	DURATION string `mapstructure:"DURATION"`
}

// func LoadConfig(path string) (open Open, err error) {
func LoadConfig() (open Open, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&open)
	return
}
func StoreApi() {
	open, err := LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	repository.IndexRepo.InitDB()
	//check db connection//////////////////////
	fmt.Println("initialization----------------")
	controllers.NeworderController(service.NeworderService(repository.NeworderRepo()))
	//initialize the pub/sub
	pubsub.SetupProducer()
	e := echo.New()

	e.Static("/", "public")
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	JWTgroup := e.Group("/api/")
	JWTgroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(open.Key),
	}))
	e.POST("/create", controllers.OrderController.Create, m.Tokenizing)
	e.POST("/item", controllers.OrderController.Additem, m.Tokenizing)
	e.GET("/order/:id", controllers.OrderController.GetOne, m.Tokenizing)
	e.GET("/order", controllers.OrderController.GetAll, m.Tokenizing)
	JWTgroup.PUT("order/:id", controllers.OrderController.Update, m.Tokenizing)
	JWTgroup.PUT("order/item/:id", controllers.OrderController.Updateitem, m.Tokenizing)
	JWTgroup.DELETE("oder/:id", controllers.OrderController.Delete, m.Tokenizing)
	//e.DELETE("loggoutall/:id", controllers.orderController.DeleteALL) logout all accounts

	// Start server
	e.Logger.Fatal(e.Start(open.Port))
}
