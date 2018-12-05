package handler

import (
	"fmt"
	"kpay/merchant"
	"kpay/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRoute(s *Server) *gin.Engine {
	r := gin.Default()

	r.POST("/register", s.RegisterMerchant)
	r.POST("/buy/product", s.BuyProduct)

	merchant := r.Group("/merchant")
	merchant.Use(s.Auth)
	merchant.GET("/:id", s.MerchantInformation)
	merchant.POST("/:id", s.UpdateMerchant)
	merchant.GET("/:id/products", s.ListAllProducts)
	merchant.POST("/:id/product", s.AddProduct)
	merchant.POST("/:id/product/:product_id", s.UpdateProduct)
	merchant.DELETE("/:id/product/:product_id", s.RemoveProduct)
	merchant.POST("/:id/report", s.SellReports)

	return r
}

type Server struct {
	// secretService  SecretService
	MerchantApiService *merchant.Manager
	ProductApiService  *product.Manager
}

func (s *Server) Auth(c *gin.Context) {
	h := map[string]string{}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	if err := s.MerchantApiService.Auth(h["username"], h["password"]); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}
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
	if err = s.MerchantApiService.Create(&merchant); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, merchant)
}

func (s *Server) MerchantInformation(c *gin.Context) {
	id := c.Param("id")
	merchant, err := s.MerchantApiService.FindByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, merchant)
}

func (s *Server) UpdateMerchant(c *gin.Context) {
	id := c.Param("id")
	var merchant merchant.Merchant
	if err := c.ShouldBindJSON(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.MerchantApiService.Update(id, &merchant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, merchant)
}

func (s *Server) ListAllProducts(c *gin.Context) {
	merchantid := c.Param("id")
	products, err := s.ProductApiService.All(merchantid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, products)
}

func (s *Server) AddProduct(c *gin.Context) {
	merchantid := c.Param("id")
	var product product.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("json: wrong params: %s", err),
		})
		return
	}
	if err = s.ProductApiService.Create(merchantid, &product); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (s *Server) UpdateProduct(c *gin.Context) {
	merchantid := c.Param("id")
	productid := c.Param("product_id")
	var product product.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("json: wrong params: %s", err),
		})
		return
	}
	if err = s.ProductApiService.Update(merchantid, productid, &product); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (s *Server) RemoveProduct(c *gin.Context) {
	merchantid := c.Param("id")
	productid := c.Param("product_id")
	if err := s.ProductApiService.Delete(merchantid, productid); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"mode":   "remove",
		"status": "success",
	})
}

func (s *Server) SellReports(c *gin.Context) {
	h := map[string]string{}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	transections, err := s.ProductApiService.Report(h["merchantid"], h["date"])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"date":       h["date"],
		"products":   transections,
		"accumulate": 100.25,
	})
}

func (s *Server) BuyProduct(c *gin.Context) {
	h := map[string]string{}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	volume, _ := strconv.Atoi(h["volume"])
	if err := s.ProductApiService.Sell(h["merchantid"], h["productid"], volume); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
