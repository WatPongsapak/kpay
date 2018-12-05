package product

import (
	"errors"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Product struct {
	ID            bson.ObjectId `bson:"_id" json:"-"`
	MerchantID    bson.ObjectId `bson:"merchant_id" json:"-"`
	Name          string        `bson:"name" json:"name"`
	Amount        int           `bson:"amount" json:"amount"`
	AmountChanges int           `bson:"amountchanges" json:"-"`
	transection   []Transection `bson:"transection" json:"transection"`
}

type Transection struct {
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	SellingVolume int           `bson:"selling_volume" json:"selling_volume"`
	Datetime      time.Time     `bson:"datetime" json:"datetime"`
}

type Manager struct {
	Collection *mgo.Collection
}

func (m *Manager) Create(merchantid string, product *Product) error {
	var merchantproducts []Product
	m.Collection.Find(bson.M{"merchant_id": merchantid}).All(&merchantproducts)
	if len(merchantproducts) > 4 {
		return errors.New("maximum products is 5")
	}
	product.ID = bson.NewObjectId()
	product.MerchantID = bson.ObjectIdHex(merchantid)
	err := m.Collection.Insert(product)
	return err
}

func (m *Manager) Update(merchantid string, productid string, products *Product) error {
	selector := bson.M{
		"_id":         bson.ObjectIdHex(productid),
		"merchant_id": bson.ObjectIdHex(merchantid),
	}
	var oldproduct Product
	err := m.Collection.Find(selector).One(&oldproduct)
	updater := bson.M{"amount": products.Amount, "amountchanges": oldproduct.Amount}
	err = m.Collection.Update(selector, bson.M{"$set": updater})
	return err
}

func (m *Manager) Delete(merchantid string, productid string) error {

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

func (m *Manager) Sell(merchantid string, productid string, volume int) error {
	selector := bson.M{
		"_id":         bson.ObjectIdHex(productid),
		"merchant_id": bson.ObjectIdHex(merchantid),
	}
	PushToArray := bson.M{"$push": bson.M{"transection": bson.M{"Datetime": time.Now(), "selling_volume": volume}}}
	err := m.Collection.Update(selector, PushToArray)
	return err
}

func (m *Manager) Report(merchantid string, date string) ([]Transection, error) {
	// layout := "02/01/2006"
	// fromDate, _ := time.Parse(layout, date)
	// toDate := fromDate.AddDate(0, 0, 1)

	var transections []Transection
	// err := m.Collection.Find(bson.M{"transection": bson.M{"$elemMatch": bson.M{"datetime": bson.M{"$gt": fromDate, "$lt": toDate}}}}).All(&transections)
	err := m.Collection.Find(bson.M{
		"merchant_id": bson.ObjectIdHex(merchantid),
	}).All(&transections)
	return transections, err
}
