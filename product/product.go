package product

import (
	"math"
	"errors"
	"kpay/transection"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Product struct {
	ID            bson.ObjectId `bson:"_id" json:"-"`
	MerchantID    bson.ObjectId `bson:"merchant_id" json:"-"`
	Name          string        `bson:"name" json:"name"`
	Amount        float64           `bson:"amount" json:"amount"`
	AmountChanges float64           `bson:"amountchanges" json:"-"`
}

type Manager struct {
	Collection     *mgo.Collection
	CollectionTran *mgo.Collection
}

func (m *Manager) Create(merchantid string, product *Product) error {
	var merchantproducts []Product
	m.Collection.Find(bson.M{"merchant_id": merchantid}).All(&merchantproducts)
	if len(merchantproducts) > 4 {
		return errors.New("maximum products is 5")
	}
	product.ID = bson.NewObjectId()
	product.Amount = math.Floor(product.Amount*100)/100.0
	product.AmountChanges = product.Amount
	product.MerchantID = bson.ObjectIdHex(merchantid)
	err := m.Collection.Insert(product)
	return err
}

func (m *Manager) Update(merchantid string, productid string, product *Product) error {
	selector := bson.M{
		"_id":         bson.ObjectIdHex(productid),
		"merchant_id": bson.ObjectIdHex(merchantid),
	}
	var oldproduct Product
	product.Amount = math.Floor(product.Amount*100)/100.0
	err := m.Collection.Find(selector).One(&oldproduct)
	updater := bson.M{"amount": product.Amount, "amountchanges": oldproduct.Amount}
	err = m.Collection.Update(selector, bson.M{"$set": updater})
	return err
}

func (m *Manager) Delete(merchantid string, productid string) error {
	var transections []transection.Transection
	m.CollectionTran.Find(bson.M{"product_id": productid}).All(&transections)
	if len(transections) > 0 {
		return errors.New("cant remove")
	}
	selector := bson.M{
		"_id":         bson.ObjectIdHex(productid),
		"merchant_id": bson.ObjectIdHex(merchantid),
	}
	err := m.Collection.Remove(selector)
	return err
}

func (m *Manager) All(merchantid string) ([]Product, error) {
	selector := bson.M{
		"merchant_id": bson.ObjectIdHex(merchantid),
	}
	var products []Product
	err := m.Collection.Find(selector).All(&products)
	return products, err
}

func (m *Manager) GetByID(merchantid string, productid string) (Product, error) {
	selector := bson.M{
		"_id":         bson.ObjectIdHex(productid),
		"merchant_id": bson.ObjectIdHex(merchantid),
	}
	var product Product
	err := m.Collection.Find(selector).One(&product)
	return product, err
}
