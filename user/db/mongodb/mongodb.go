package mongodb

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/microservices-demo/user/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	name     string
	password string
	host     string
	dbName   = "users"
	//ErrInvalidHexID represents a entity id that is not a valid bson ObjectID
	ErrInvalidHexID = errors.New("Invalid Id Hex")
)

func init() {
	flag.StringVar(&name, "mongo-user", os.Getenv("MONGO_USER"), "Mongo user")
	flag.StringVar(&password, "mongo-password", os.Getenv("MONGO_PASS"), "Mongo password")
	flag.StringVar(&host, "mongo-host", os.Getenv("mongo"), "Mongo host")
}

// Mongo meets the Database interface requirements
type Mongo struct {
	Client *mongo.Client
}

func (m *Mongo) collection(name string) *mongo.Collection {
	return m.Client.Database(dbName).Collection(name)
}

// Init MongoDB
func (m *Mongo) Init() error {
	u := getURL()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var err error
	m.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(u.String()))
	if err != nil {
		return err
	}
	return m.EnsureIndexes()
}

// MongoUser is a wrapper for the users
type MongoUser struct {
	users.User `bson:",inline"`
	ID         primitive.ObjectID   `bson:"_id"`
	AddressIDs []primitive.ObjectID `bson:"addresses"`
	CardIDs    []primitive.ObjectID `bson:"cards"`
}

// New Returns a new MongoUser
func New() MongoUser {
	u := users.New()
	return MongoUser{
		User:       u,
		AddressIDs: make([]primitive.ObjectID, 0),
		CardIDs:    make([]primitive.ObjectID, 0),
	}
}

// AddUserIDs adds userID as string to user
func (mu *MongoUser) AddUserIDs() {
	if mu.User.Addresses == nil {
		mu.User.Addresses = make([]users.Address, 0)
	}
	for _, id := range mu.AddressIDs {
		mu.User.Addresses = append(mu.User.Addresses, users.Address{
			ID: id.Hex(),
		})
	}
	if mu.User.Cards == nil {
		mu.User.Cards = make([]users.Card, 0)
	}
	for _, id := range mu.CardIDs {
		mu.User.Cards = append(mu.User.Cards, users.Card{ID: id.Hex()})
	}
	mu.User.UserID = mu.ID.Hex()
}

// MongoAddress is a wrapper for Address
type MongoAddress struct {
	users.Address `bson:",inline"`
	ID            primitive.ObjectID `bson:"_id"`
}

// AddID ObjectID as string
func (m *MongoAddress) AddID() {
	m.Address.ID = m.ID.Hex()
}

// MongoCard is a wrapper for Card
type MongoCard struct {
	users.Card `bson:",inline"`
	ID         primitive.ObjectID `bson:"_id"`
}

// AddID ObjectID as string
func (m *MongoCard) AddID() {
	m.Card.ID = m.ID.Hex()
}

// CreateUser Insert user to MongoDB, including connected addresses and cards, update passed in user with Ids
func (m *Mongo) CreateUser(u *users.User) error {
	ctx := context.TODO()
	id := primitive.NewObjectID()
	mu := New()
	mu.User = *u
	mu.ID = id
	var carderr error
	var addrerr error
	mu.CardIDs, carderr = m.createCards(u.Cards)
	mu.AddressIDs, addrerr = m.createAddresses(u.Addresses)
	c := m.collection("customers")
	_, err := c.ReplaceOne(ctx, bson.M{"_id": mu.ID}, mu, options.Replace().SetUpsert(true))
	if err != nil {
		m.cleanAttributes(mu)
		return err
	}
	mu.User.UserID = mu.ID.Hex()
	if carderr != nil || addrerr != nil {
		return fmt.Errorf("%v %v", carderr, addrerr)
	}
	*u = mu.User
	return nil
}

func (m *Mongo) createCards(cs []users.Card) ([]primitive.ObjectID, error) {
	ctx := context.TODO()
	ids := make([]primitive.ObjectID, 0)
	for k, ca := range cs {
		id := primitive.NewObjectID()
		mc := MongoCard{Card: ca, ID: id}
		c := m.collection("cards")
		_, err := c.ReplaceOne(ctx, bson.M{"_id": mc.ID}, mc, options.Replace().SetUpsert(true))
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
		cs[k].ID = id.Hex()
	}
	return ids, nil
}

func (m *Mongo) createAddresses(as []users.Address) ([]primitive.ObjectID, error) {
	ctx := context.TODO()
	ids := make([]primitive.ObjectID, 0)
	for k, a := range as {
		id := primitive.NewObjectID()
		ma := MongoAddress{Address: a, ID: id}
		c := m.collection("addresses")
		_, err := c.ReplaceOne(ctx, bson.M{"_id": ma.ID}, ma, options.Replace().SetUpsert(true))
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
		as[k].ID = id.Hex()
	}
	return ids, nil
}

func (m *Mongo) cleanAttributes(mu MongoUser) error {
	ctx := context.TODO()
	c := m.collection("addresses")
	_, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": mu.AddressIDs}})
	c = m.collection("cards")
	_, err = c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": mu.CardIDs}})
	return err
}

func (m *Mongo) appendAttributeId(attr string, id primitive.ObjectID, userid string) error {
	ctx := context.TODO()
	oid, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}
	c := m.collection("customers")
	_, err = c.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$addToSet": bson.M{attr: id}})
	return err
}

func (m *Mongo) removeAttributeId(attr string, id primitive.ObjectID, userid string) error {
	ctx := context.TODO()
	oid, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}
	c := m.collection("customers")
	_, err = c.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$pull": bson.M{attr: id}})
	return err
}

// GetUserByName Get user by their name
func (m *Mongo) GetUserByName(name string) (users.User, error) {
	ctx := context.TODO()
	c := m.collection("customers")
	mu := New()
	err := c.FindOne(ctx, bson.M{"username": name}).Decode(&mu)
	mu.AddUserIDs()
	return mu.User, err
}

// GetUser Get user by their object id
func (m *Mongo) GetUser(id string) (users.User, error) {
	ctx := context.TODO()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return users.New(), errors.New("Invalid Id Hex")
	}
	c := m.collection("customers")
	mu := New()
	err = c.FindOne(ctx, bson.M{"_id": oid}).Decode(&mu)
	mu.AddUserIDs()
	return mu.User, err
}

// GetUsers Get all users
func (m *Mongo) GetUsers() ([]users.User, error) {
	ctx := context.TODO()
	c := m.collection("customers")
	cursor, err := c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var mus []MongoUser
	if err := cursor.All(ctx, &mus); err != nil {
		return nil, err
	}
	us := make([]users.User, 0)
	for _, mu := range mus {
		mu.AddUserIDs()
		us = append(us, mu.User)
	}
	return us, nil
}

// GetUserAttributes given a user, load all cards and addresses connected to that user
func (m *Mongo) GetUserAttributes(u *users.User) error {
	ctx := context.TODO()
	ids := make([]primitive.ObjectID, 0)
	for _, a := range u.Addresses {
		oid, err := primitive.ObjectIDFromHex(a.ID)
		if err != nil {
			return ErrInvalidHexID
		}
		ids = append(ids, oid)
	}
	c := m.collection("addresses")
	cursor, err := c.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return err
	}
	var ma []MongoAddress
	if err := cursor.All(ctx, &ma); err != nil {
		return err
	}
	na := make([]users.Address, 0)
	for _, a := range ma {
		a.Address.ID = a.ID.Hex()
		na = append(na, a.Address)
	}
	u.Addresses = na

	ids = make([]primitive.ObjectID, 0)
	for _, ca := range u.Cards {
		oid, err := primitive.ObjectIDFromHex(ca.ID)
		if err != nil {
			return ErrInvalidHexID
		}
		ids = append(ids, oid)
	}
	cc := m.collection("cards")
	cursor, err = cc.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return err
	}
	var mc []MongoCard
	if err := cursor.All(ctx, &mc); err != nil {
		return err
	}
	nc := make([]users.Card, 0)
	for _, ca := range mc {
		ca.Card.ID = ca.ID.Hex()
		nc = append(nc, ca.Card)
	}
	u.Cards = nc
	return nil
}

// GetCard Gets card by objects Id
func (m *Mongo) GetCard(id string) (users.Card, error) {
	ctx := context.TODO()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return users.Card{}, errors.New("Invalid Id Hex")
	}
	c := m.collection("cards")
	mc := MongoCard{}
	err = c.FindOne(ctx, bson.M{"_id": oid}).Decode(&mc)
	mc.AddID()
	return mc.Card, err
}

// GetCards Gets all cards
func (m *Mongo) GetCards() ([]users.Card, error) {
	ctx := context.TODO()
	c := m.collection("cards")
	cursor, err := c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var mcs []MongoCard
	if err := cursor.All(ctx, &mcs); err != nil {
		return nil, err
	}
	cs := make([]users.Card, 0)
	for _, mc := range mcs {
		mc.AddID()
		cs = append(cs, mc.Card)
	}
	return cs, nil
}

// CreateCard adds card to MongoDB
func (m *Mongo) CreateCard(ca *users.Card, userid string) error {
	if userid != "" {
		if _, err := primitive.ObjectIDFromHex(userid); err != nil {
			return errors.New("Invalid Id Hex")
		}
	}
	ctx := context.TODO()
	c := m.collection("cards")
	id := primitive.NewObjectID()
	mc := MongoCard{Card: *ca, ID: id}
	_, err := c.ReplaceOne(ctx, bson.M{"_id": mc.ID}, mc, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	if userid != "" {
		err = m.appendAttributeId("cards", mc.ID, userid)
		if err != nil {
			return err
		}
	}
	mc.AddID()
	*ca = mc.Card
	return err
}

// GetAddress Gets an address by object Id
func (m *Mongo) GetAddress(id string) (users.Address, error) {
	ctx := context.TODO()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return users.Address{}, errors.New("Invalid Id Hex")
	}
	c := m.collection("addresses")
	ma := MongoAddress{}
	err = c.FindOne(ctx, bson.M{"_id": oid}).Decode(&ma)
	ma.AddID()
	return ma.Address, err
}

// GetAddresses gets all addresses
func (m *Mongo) GetAddresses() ([]users.Address, error) {
	ctx := context.TODO()
	c := m.collection("addresses")
	cursor, err := c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var mas []MongoAddress
	if err := cursor.All(ctx, &mas); err != nil {
		return nil, err
	}
	as := make([]users.Address, 0)
	for _, ma := range mas {
		ma.AddID()
		as = append(as, ma.Address)
	}
	return as, nil
}

// CreateAddress Inserts Address into MongoDB
func (m *Mongo) CreateAddress(a *users.Address, userid string) error {
	if userid != "" {
		if _, err := primitive.ObjectIDFromHex(userid); err != nil {
			return errors.New("Invalid Id Hex")
		}
	}
	ctx := context.TODO()
	c := m.collection("addresses")
	id := primitive.NewObjectID()
	ma := MongoAddress{Address: *a, ID: id}
	_, err := c.ReplaceOne(ctx, bson.M{"_id": ma.ID}, ma, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	if userid != "" {
		err = m.appendAttributeId("addresses", ma.ID, userid)
		if err != nil {
			return err
		}
	}
	ma.AddID()
	*a = ma.Address
	return err
}

// Delete removes an entity from MongoDB
func (m *Mongo) Delete(entity, id string) error {
	ctx := context.TODO()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid Id Hex")
	}
	c := m.collection(entity)
	if entity == "customers" {
		u, err := m.GetUser(id)
		if err != nil {
			return err
		}
		aids := make([]primitive.ObjectID, 0)
		for _, a := range u.Addresses {
			aoid, _ := primitive.ObjectIDFromHex(a.ID)
			aids = append(aids, aoid)
		}
		cids := make([]primitive.ObjectID, 0)
		for _, ca := range u.Cards {
			coid, _ := primitive.ObjectIDFromHex(ca.ID)
			cids = append(cids, coid)
		}
		ac := m.collection("addresses")
		ac.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": aids}})
		cc := m.collection("cards")
		cc.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": cids}})
	} else {
		cc := m.collection("customers")
		cc.UpdateMany(ctx, bson.M{}, bson.M{"$pull": bson.M{entity: oid}})
	}
	_, err = c.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func getURL() url.URL {
	ur := url.URL{
		Scheme: "mongodb",
		Host:   host,
		Path:   dbName,
	}
	if name != "" {
		u := url.UserPassword(name, password)
		ur.User = u
	}
	return ur
}

// EnsureIndexes ensures username is unique
func (m *Mongo) EnsureIndexes() error {
	ctx := context.TODO()
	c := m.collection("customers")
	_, err := c.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return m.Client.Ping(ctx, nil)
}
