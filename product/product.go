package product

import (
	"time"
	"kpay/product"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Product struct {
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	Name          string        `bson:"name" json:"name"`
	Amount        int           `bson:"amount" json:"amount"`
	AmountChanges []transection `bson:"amoutchanges" json:"amoutchanges"`
}

type transection struct {
	exchangetype string    `bson:"type" json:"type"`
	Amount       int       `bson:"amount" json:"amount"`
	Datetime     time.Time `bson:"datetime" json:"datetime"`
}

type Manager struct {
	Collection *mgo.Collection
}

func (m *Manager) Insert(product *Product) error {
	product.ID = bson.NewObjectId()
	err := m.Collection.Insert(product)
	return err
}

func (m *Manager) Update(product *Product) error {
	err := m.Collection.UpdateId(product.ID, &product)
	return err
}

func (m *Manager) FindByID(id string) (*Product, error) {
	var product *Product
	err := m.Collection.FindId(bson.ObjectIdHex(id)).One(&product)
	return product, err
}

func (m *Manager) Delete(product *Product) error {
	err := m.Collection.Remove(&product)
	return err
}

func (m *Manager) All() ([]Product, error) {
	var products []Product
	err := m.Collection.Find(bson.M{}).All(&products)
	return products, err
}

func (m *Manager) sell() ([]Product, error) {
	var products []Product
	err := m.Collection.Find(bson.M{}).All(&products)
	return products, err
}

func (m *Manager) buy(product *Product, merchant *merchant.Merchant) ([]Product, error) {
	var products []Product
	err := m.Collection.Find(bson.M{}).All(&products)
	return products, err
}
