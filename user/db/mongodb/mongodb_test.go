package mongodb

import (
	"testing"

	"github.com/microservices-demo/user/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNew(t *testing.T) {
	m := New()
	if m.AddressIDs == nil || m.CardIDs == nil {
		t.Error("Expected non nil arrays")
	}
}

func TestAddUserIDs(t *testing.T) {
	m := New()
	uid := primitive.NewObjectID()
	cid := primitive.NewObjectID()
	aid := primitive.NewObjectID()
	m.ID = uid
	m.AddressIDs = append(m.AddressIDs, aid)
	m.CardIDs = append(m.CardIDs, cid)
	m.AddUserIDs()
	if len(m.Addresses) != 1 || len(m.Cards) != 1 {
		t.Error("Expected one card and one address added.")
	}
	if m.Addresses[0].ID != aid.Hex() {
		t.Error("Expected matching Address Hex")
	}
	if m.Cards[0].ID != cid.Hex() {
		t.Error("Expected matching Card Hex")
	}
	if m.UserID != uid.Hex() {
		t.Error("Expected matching User Hex")
	}
}

func TestAddressAddId(t *testing.T) {
	m := MongoAddress{Address: users.Address{}}
	id := primitive.NewObjectID()
	m.ID = id
	m.AddID()
	if m.Address.ID != id.Hex() {
		t.Error("Expected matching Address Hex")
	}
}

func TestCardAddId(t *testing.T) {
	m := MongoCard{Card: users.Card{}}
	id := primitive.NewObjectID()
	m.ID = id
	m.AddID()
	if m.Card.ID != id.Hex() {
		t.Error("Expected matching Card Hex")
	}
}

func TestGetURL(t *testing.T) {
	name = "test"
	password = "password"
	host = "thishostshouldnotexist:3038"
	u := getURL()
	if u.String() != "mongodb://test:password@thishostshouldnotexist:3038/users" {
		t.Error("expected url mismatch")
	}
}
