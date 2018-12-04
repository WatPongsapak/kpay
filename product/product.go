package product

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Product struct {
	ID         bson.ObjectId `bson:"_id" json:"_id"`
	MerchantID bson.ObjectId `bson:"merchant_id" json:"merchant_id"`
	Name       string        `bson:"name" json:"name"`
	Amount     int           `bson:"amount" json:"amount"`
}

type Manager struct {
	Collection *mgo.Collection
}

func (m *Manager) Create(merchantid string, products *Product) error {
	products.ID = bson.NewObjectId()
	products.MerchantID = bson.ObjectIdHex(merchantid)
	err := m.Collection.Insert(products)
	return err
}

func (m *Manager) Update(merchantid string, productid string, products *Product) error {
	err := m.Collection.Update(bson.M{"_id": bson.ObjectIdHex(productid), "merchant_id": bson.ObjectIdHex(merchantid)}, bson.M{"$set": bson.M{"name": products.Name, "amount": products.Amount}})
	return err
}

func (m *Manager) Delete(merchantid string, productid string, products *Product) error {
	err := m.Collection.Remove(bson.M{"merchant_id": bson.ObjectIdHex(merchantid), "_id": bson.ObjectIdHex(productid)})
	return err
}

func (m *Manager) All(merchantid string) ([]Product, error) {
	var products []Product
	err := m.Collection.Find(bson.M{"merchant_id": bson.ObjectIdHex(merchantid)}).All(&products)
	return products, err
}
