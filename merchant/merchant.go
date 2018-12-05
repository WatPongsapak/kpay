package merchant

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sethvargo/go-password/password"
)

type Merchant struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	Name        string        `bson:"name" json:"name"`
	BankAccount string        `bson:"bankaccount" json:"bankaccount"`
	Username    string        `bson:"username" json:"username"`
	Password    string        `bson:"password" json:"password"`
}

type Manager struct {
	Collection *mgo.Collection
}

func (m *Manager) Create(merchant *Merchant) error {
	var merchants []Merchant
	m.Collection.Find(bson.M{"bankaccount": merchant.BankAccount}).All(&merchants)
	if len(merchants) > 0 {
		return errors.New("bank account have not same")
	}
	merchant.ID = bson.NewObjectId()
	res, _ := password.Generate(32, 10, 10, false, false)
	merchant.Password = res
	merchant.Username = merchant.Name + strconv.Itoa(rand.Intn(10000))
	err := m.Collection.Insert(merchant)
	return err
}

func (m *Manager) Update(id string, merchant *Merchant) error {
	err := m.Collection.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"name": merchant.Name}})
	return err
}

func (m *Manager) FindByID(id string) (*Merchant, error) {
	var merchant *Merchant
	err := m.Collection.FindId(bson.ObjectIdHex(id)).One(&merchant)
	return merchant, err
}

func (m *Manager) Auth(username string, password string) error {
	var merchants []Merchant
	m.Collection.Find(bson.M{"username": username, "password": password}).All(&merchants)
	if len(merchants) == 0 {
		return errors.New("Auth fail")
	}
	return nil
}
