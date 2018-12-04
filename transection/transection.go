package transection

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Transection struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	MerchantID   bson.ObjectId `bson:"merchant_id" json:"merchant_id"`
	ProductID    bson.ObjectId `bson:"product_id" json:"product_id"`
	exchangetype string        `bson:"type" json:"type"`
	Amount       int           `bson:"amount" json:"amount"`
	Datetime     time.Time     `bson:"datetime" json:"datetime"`
}

type Manager struct {
	Collection *mgo.Collection
}

func (m *Manager) Create(merchantid string, transection *Transection) error {
	transection.ID = bson.NewObjectId()
	transection.MerchantID = bson.ObjectIdHex(merchantid)
	err := m.Collection.Insert(transection)
	return err
}

// func (m *Manager) Update(merchantid string, productid string, transection *Transection) error {
// 	err := m.Collection.Update(bson.M{"_id": bson.ObjectIdHex(productid), "merchant_id": bson.ObjectIdHex(merchantid)}, bson.M{"$set": bson.M{"name": transection.Name, "amount": transection.Amount}})
// 	return err
// }

func (m *Manager) Delete(merchantid string, productid string, transection *Transection) error {
	err := m.Collection.Remove(bson.M{"merchant_id": bson.ObjectIdHex(merchantid), "_id": bson.ObjectIdHex(productid)})
	return err
}
