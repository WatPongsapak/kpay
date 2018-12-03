package main

import (
	"kpay/merchant"

	"github.com/globalsign/mgo"
)

const url = "mongodb://admin:kbtgadmin@cluster0-shard-00-00-o2soi.mongodb.net:27017,cluster0-shard-00-01-o2soi.mongodb.net:27017,cluster0-shard-00-02-o2soi.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin"

type Person struct {
	Name  string
	Phone string
}

func main() {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	merchantManager := &merchant.Manager{
		Collection: session.DB("kpay").C("merchant"),
	}

	// err = merchantManager.Insert(&merchant.Merchant{
	// 	Name:        "wat",
	// 	BankAccount: "wat",
	// 	Username:    "wat",
	// 	Password:    "wat",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	merchant, _ := merchantManager.FindByID("5c04f12f8b932b2b04cd708a")
	merchant.Name = "wwww"
	merchantManager.Update(merchant)
	// fmt.Println()
}
