package common

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoSession wraps mgo.Session
type Session interface {
	DB(name string) DataLayer
	Close()
}

func NewMongoSession(host string) (Session, error) {
	session, err := mgo.Dial(host)
	return MongoSession{session}, err
}

type MongoSession struct {
	*mgo.Session
}

func (s MongoSession) DB(name string) DataLayer {
	return MongoDatabase{s.Session.DB(name)}
}

func (s MongoSession) Close() {
	s.Session.Close()
}

// MongoDatabase wraps mgo.Database
type DataLayer interface {
	C(name string) Collection
}

type MongoDatabase struct {
	*mgo.Database
}

func (d MongoDatabase) C(name string) Collection {
	return MongoCollection{d.Database.C(name)}
}

// MongoCollection wraps mgo.Collection
type Collection interface {
	Find(query interface{}) Query
	Insert(docs ...interface{}) error
}

type MongoCollection struct {
	*mgo.Collection
}

func (c MongoCollection) Find(query interface{}) Query {
	return MongoQuery{c.Collection.Find(query)}
}

func (c MongoCollection) Insert(docs ...interface{}) error {
	return c.Collection.Insert(docs)
}

// MongoQuery wraps mgo.Query
type Query interface {
	One(result interface{}) error
}

type MongoQuery struct {
	*mgo.Query
}

func (q MongoQuery) One(result interface{}) error {
	return q.Query.One(result)
}

// MongoAdapter fulfills Adapter interface
func NewMongoAdapter(host string, database string) (Adapter, error) {
	if session, err := NewMongoSession(host); err != nil {
		return nil, err
	} else {
		return &MongoAdapter{session, database}, nil
	}
}

type MongoAdapter struct {
	session  Session
	database string
}

func (a *MongoAdapter) db() DataLayer {
	return a.session.DB(a.database)
}

func (a *MongoAdapter) Find(collection string, id string) (*Document, error) {
	c := a.db().C(collection)
	doc := &Document{}
	err := c.Find(bson.M{"id": id}).One(doc)
	return doc, err
}

func (a *MongoAdapter) Create(collection string, doc *Document) error {
	c := a.db().C(collection)
	return c.Insert(doc)
}