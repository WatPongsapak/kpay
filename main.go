package main

import (
	"kpay/handler"
	"kpay/merchant"
	"kpay/product"
	"kpay/transection"
	"os"

	"github.com/globalsign/mgo"
)

// const url = "127.0.0.1:27017"
const url = "mongodb://admin:admin@cluster0-shard-00-00-l5ohv.mongodb.net:27017,cluster0-shard-00-01-l5ohv.mongodb.net:27017,cluster0-shard-00-02-l5ohv.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin"

func main() {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()
	s := &handler.Server{
		MerchantApiService: &merchant.Manager{
			Collection: session.DB("kpay").C("merchant"),
		},
		ProductApiService: &product.Manager{
			Collection: session.DB("kpay").C("product"),
			CollectionTran: session.DB("kpay").C("transection"),
		},
		TransectionApiService: &transection.Manager{
			Collection: session.DB("kpay").C("transection"),
		},
	}

	r := handler.SetupRoute(s)

	r.Run(":" + os.Getenv("PORT"))
}
