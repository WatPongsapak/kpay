package main

import (
	"kpay/merchant"
	"kpay/product"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

const url = "mongodb://admin:kbtgadmin@cluster0-shard-00-00-o2soi.mongodb.net:27017,cluster0-shard-00-01-o2soi.mongodb.net:27017,cluster0-shard-00-02-o2soi.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin"


type MerchantApiService interface {
	Insert(merchant *Merchant) error
	Update(merchant *Merchant) error
	FindByID(id string) (*Merchant, error)
	Delete(merchant *Merchant)
	All() ([]Merchant, error)
}

type ProductApiService interface {
	Insert(merchant *Product) error
	Update(merchant *Product) error
	FindByID(id string) (*Product, error)
	Delete(merchant *Product) error
	All() ([]Product, error)
}

type SecretService interface {
	Insert(s *secret.Secret) error
	FindSecretKey(s *secret.Secret) error
}

type Server struct {
	productApiService ProductApiService
	// secretService  SecretService
	merchantApiService  MerchantApiService
}

func setupRoute(s *Server) *gin.Engine {
	r := gin.Default()
	// accs := r.Group("/accs")
	admin := r.Group("/admin")
	user := r.Group("/users")
	account := r.Group("/bankAccounts")
	admin.Use(gin.BasicAuth(gin.Accounts{
		"admin": "1234",
	}))
	admin.POST("/secrets", s.CreateSecret)
	// user.Use(s.Auth)
	// account.Use(s.Auth)
	user.GET("", s.AllUsers)
	user.POST("", s.CreateUser)
	user.GET("/:id", s.GetUserByID)
	user.PUT("/:id", s.UpdateUser)
	user.DELETE("/:id", s.DeleteUser)
	user.POST("/:id/bankAccounts", s.CreateAccount)
	user.GET("/:id/bankAccounts", s.GetAccountByUserID)
	account.PUT("/:id/withdraw", s.AccountWithdraw)
	account.PUT("/:id/deposit", s.AccountDeposit)
	account.DELETE("/:id", s.DeleteAccount)

	return r
}

func main() {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()
	s := &Server{
		merchantApiService: &merchant.Manager{
			Collection: session.DB("kpay").C("merchant"),
		},
		productApiService: &product.Manager{
			Collection: session.DB("kpay").C("merchant"),
		},
	}

	r := setupRoute(s)

	r.Run(":" + os.Getenv("PORT"))