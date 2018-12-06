package transection

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Transection struct {
	ID            bson.ObjectId `bson:"_id" json:"-"`
	ProductID     bson.ObjectId `bson:"product_id" json:"-"`
	MerchantID    bson.ObjectId `bson:"merchant_id" json:"-"`
	ProductName   string        `bson:"product_name" json:"product_name"`
	SellingVolume int           `bson:"selling_volume" json:"selling_volume"`
	Amount        float64       `bson:"amount" json:"-"`
	Date          string        `bson:"date" json:"-"`
}

type Manager struct {
	Collection *mgo.Collection
}

func (m *Manager) Sell(merchantid string, productid string, volume int, amount float64, name string) error {
	var transection Transection
	transection.ID = bson.NewObjectId()
	transection.MerchantID = bson.ObjectIdHex(merchantid)
	transection.ProductID = bson.ObjectIdHex(productid)
	transection.SellingVolume = volume
	transection.Amount = amount
	transection.ProductName = name
	currentTime := time.Now()
	transection.Date = currentTime.Format("02/01/2006")
	err := m.Collection.Insert(transection)
	return err
}

func (m *Manager) Report(merchantid string, date string) ([]Transection, error) {
	var transections []Transection
	err := m.Collection.Find(bson.M{
		"date":        date,
		"merchant_id": bson.ObjectIdHex(merchantid),
	}).All(&transections)
	return transections, err
}
