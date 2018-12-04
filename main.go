package main

import (
	"kpay/merchant"
	"kpay/product"
	"kpay/transection"
	"os"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

const url = "mongodb://admin:kbtgadmin@cluster0-shard-00-00-o2soi.mongodb.net:27017,cluster0-shard-00-01-o2soi.mongodb.net:27017,cluster0-shard-00-02-o2soi.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin"


// type MerchantApiService interface {
// 	Insert(merchant *merchant.Merchant) error
// 	Update(id string, merchant *merchant.Merchant) error
// 	FindByID(id string) (*merchant.Merchant, error)
// 	AllProducts(id string) []merchant.Product
// 	AddProduct(id string, product *merchant.Product) error
// }

// type SecretService interface {
// 	Insert(s *secret.Secret) error
// 	FindSecretKey(s *secret.Secret) error
// }

type Server struct {
	// secretService  SecretService
	merchantApiService  *merchant.Manager
	productApiService  *product.Manager
	transectionApiService  *transection.Manager
}

func (s *Server) RegisterMerchant(c *gin.Context) {
	var merchant merchant.Merchant
	err := c.ShouldBindJSON(&merchant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("json: wrong params: %s", err),
		})
		return
	}
	if err = s.merchantApiService.Create(&merchant); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, merchant)
}

func (s *Server) MerchantInformation(c *gin.Context) {
	id := c.Param("id")
	merchant, err := s.merchantApiService.FindByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, merchant)
}

func (s *Server) UpdateMerchant(c *gin.Context) {
	id := c.Param("id")
	if id == "register"{
		s.RegisterMerchant(c)
		return
	}
	var merchant merchant.Merchant
	if err := c.ShouldBindJSON(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.merchantApiService.Update(id,&merchant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, merchant)
}

// func (s *Server) ListAllProducts(c *gin.Context) {
// 	merchantid := c.Param("id")
// 	products := s.productApiService.All(merchantid)
// 	c.JSON(http.StatusOK, products)
// }

// func (s *Server) AddProduct(c *gin.Context) {
// 	merchantid := c.Param("id")
// 	var product product.Product
// 	err := c.ShouldBindJSON(&product)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"object":  "error",
// 			"message": fmt.Sprintf("json: wrong params: %s", err),
// 		})
// 		return
// 	}
// 	if err = s.productApiService.Creat(merchantid,&product); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusCreated, product)
// }

// func (s *Server) UpdateProduct(c *gin.Context) {
// 	merchantid := c.Param("id")
// 	productid := c.Param("product_id")
// 	var product product.Product
// 	err := c.ShouldBindJSON(&product)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"object":  "error",
// 			"message": fmt.Sprintf("json: wrong params: %s", err),
// 		})
// 		return
// 	}
// 	if err = s.productApiService.Update(merchantid,productid,&product); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOk, product)
// }

// func (s *Server) RemoveProduct(c *gin.Context) {
// 	merchantid := c.Param("id")
// 	productid := c.Param("product_id")
// 	if err = s.productApiService.Delete(merchantid,productid); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOk, product)
// }

// func (s *Server) SellReports(c *gin.Context) {
// 	h := map[string]string{}
// 	if err := c.ShouldBindJSON(&h); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, err)
// 		return
// 	}
// 	transections, err = s.transectionApiService.Report(h["date"]);
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOk, gin.H{
// 		date: h["date"],
// 		products: transections,
// 		accumulate: 100.25 ,
// 	})
// }

// func (s *Server) BuyProduct(c *gin.Context) {
// 	var transection transection.Transection
// 	err := c.ShouldBindJSON(&transection)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"object":  "error",
// 			"message": fmt.Sprintf("json: wrong params: %s", err),
// 		})
// 		return
// 	}
// 	if err = s.transectionApiService.Create(transection); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOk, transection)
// }



func setupRoute(s *Server) *gin.Engine {
	r := gin.Default()
	// accs := r.Group("/accs")
	// admin := r.Group("/admin")
	// user := r.Group("/users")
	// account := r.Group("/bankAccounts")
	// admin.Use(gin.BasicAuth(gin.Accounts{
	// 	"admin": "1234",
	// }))
	// admin.POST("/secrets", s.CreateSecret)
	// user.Use(s.Auth)
	// account.Use(s.Auth)
	r.GET("/merchant/:id",s.MerchantInformation)
	r.POST("/merchant/:id",s.UpdateMerchant)
	// r.GET("/merchant/:id/products",s.ListAllProducts)
	// r.POST("/merchant/:id/product",s.AddProduct)
	// r.POST("/merchant/:id/product/:product_id",s.UpdateProduct)
	// r.DELETE("/merchant/:id/product/:product_id",s.RemoveProduct)
	// r.POST("/merchant/:id/report",s.SellReports)
	// r.POST("/buy/product",s.BuyProduct)

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
		transectionApiService: &transection.Manager{
			Collection: session.DB("kpay").C("merchant"),
		},
	}

	r := setupRoute(s)

	r.Run(":" + os.Getenv("PORT"))
}