package gmongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// Client is a multi-database with the pool of connections to multi-MongoDB deployment.
type Client struct {
	dbs map[string]*Database
}

// Database is a multi-collection with the pool of connections to a MongoDB deployment.
type Database struct {
	db          *mongo.Database
	collections map[string]*Collection
	address     string
	opts        *options.ClientOptions
}

// GetCollection from a database of MongoDB deployment.
func (d *Database) GetCollection(name string) (*Collection, error) {
	collection := d.collections[name]

	if collection == nil {

		collection = d.newCollection(name)

		if collection == nil {
			return nil, fmt.Errorf("the collection isn't found")
		}
		d.collections[name] = collection

	}

	return collection, nil
}

// newCollection initializes a collection.
func (d *Database) newCollection(name string) *Collection {
	return &Collection{d.db.Collection(name)}
}

// NewDatabase initializes a database and saves in local cache of client.
// It can customize database or its pool by setting one or more DatabaseOption,
// also verifies that the database was created successfully.
func (c *Client) NewDatabase(address, name string, opts ...DatabaseOption) (*Database, error) {
	database := &Database{
		collections: make(map[string]*Collection, 0),
		opts:        new(options.ClientOptions),
	}

	for _, opt := range opts {
		opt(database)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(address), database.opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = conn.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	database.db = conn.Database(name)
	c.dbs[name] = database
	return database, nil
}

// NewClient initializes a databases' local cache.
func NewClient() *Client {
	return &Client{
		dbs: make(map[string]*Database, 0),
	}
}

// DatabaseOption customizes database or its pool.
type DatabaseOption func(*Database)

// SetMaxPoolSize for pool.
func SetMaxPoolSize(maxPoolSize uint64) DatabaseOption {
	return func(database *Database) {
		database.opts.MaxPoolSize = &maxPoolSize
	}
}

// SetMinPoolSize for pool.
func SetMinPoolSize(minPoolSize uint64) DatabaseOption {
	return func(database *Database) {
		database.opts.MinPoolSize = &minPoolSize
	}
}

// SetMaxConnIdleTime for pool.
func SetMaxConnIdleTime(maxConnIdleTime time.Duration) DatabaseOption {
	return func(database *Database) {
		database.opts.MaxConnIdleTime = &maxConnIdleTime
	}
}
