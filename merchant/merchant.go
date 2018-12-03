package merchant

import (
	"kpay/product"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Merchant struct {
	ID          bson.ObjectId 		`bson:"_id" json:"_id"`
	Name        string        		`bson:"name" json:"name"`
	BankAccount string        		`bson:"bankaccount" json:"bankaccount"`
	Username    string        		`bson:"username" json:"username"`
	Password    string        		`bson:"password" json:"password"`
	Product     []product.Product   `bson:"product" json:"product"`
}

type Manager struct {
	Collection *mgo.Collection
}

func (m *Manager) Insert(merchant *Merchant) error {
	merchant.ID = bson.NewObjectId()
	err := m.Collection.Insert(merchant)
	return err
}

func (m *Manager) Update(merchant *Merchant) error {
	err := m.Collection.UpdateId(merchant.ID, &merchant)
	return err
}

func (m *Manager) FindByID(id string) (*Merchant, error) {
	var merchant *Merchant
	err := m.Collection.FindId(bson.ObjectIdHex(id)).One(&merchant)
	return merchant, err
}

func (m *Manager) Delete(merchant *Merchant) error {
	err := m.Collection.Remove(&merchant)
	return err
}

func (m *Manager) All() ([]Merchant, error) {
	var merchants []Merchant
	err := m.Collection.Find(bson.M{}).All(&merchants)
	return merchants, err
}
