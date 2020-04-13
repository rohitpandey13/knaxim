package mongo

import (
	"context"
	"sync"

	"git.maxset.io/web/knaxim/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"log"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

// Database implements database.Database for mongodb
type Database struct {
	trackOwners
	URI       string            `json:"uri"`
	DBName    string            `json:"db"`
	CollNames map[string]string `json:"coll"`
	client    *mongo.Client
	ctx       context.Context
	cancel    context.CancelFunc
}

// Init tests the connection to the database
// if reset is true, empties out the databases and sets up indexex
func (d *Database) Init(ctx context.Context, reset bool) error {
	d.CollNames = initcoll(d.CollNames)
	//try connecting
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//testclient, err := mongo.Connect(ctx, d.URI)
	testclient, err := mongo.Connect(ctx, options.Client().ApplyURI(d.URI))
	if err != nil {
		return err
	}
	defer testclient.Disconnect(ctx)
	err = testclient.Ping(ctx, nil)
	if err != nil {
		return err
	}
	if reset {
		if err = testclient.Database(d.DBName).Drop(ctx); err != nil {
			return err
		}
		var wg sync.WaitGroup
		wg.Add(10)
		indexctx, cancel := context.WithCancel(ctx)
		defer cancel()
		cherr := make(chan error, 8)
		go func() {
			//user
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["user"]).Indexes()
			var err error
			if _, err = I.CreateMany(
				indexctx,
				[]mongo.IndexModel{
					mongo.IndexModel{
						Keys:    bson.M{"id": 1},
						Options: options.Index().SetUnique(true),
					},
					mongo.IndexModel{
						Keys:    bson.M{"name": 1},
						Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"name": bson.M{"$exists": true}}),
					},
				}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
			if _, err = I.CreateOne(indexctx, mongo.IndexModel{
				Keys:    bson.M{"name": 1},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"name": bson.M{"$exists": true}}),
			}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		go func() {
			//group
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["group"]).Indexes()
			var err error
			if _, err = I.CreateMany(
				indexctx,
				[]mongo.IndexModel{
					mongo.IndexModel{
						Keys:    bson.M{"id": 1},
						Options: options.Index().SetUnique(true),
					},
					mongo.IndexModel{
						Keys:    bson.M{"name": 1},
						Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"name": bson.M{"$exists": true}}),
					},
					mongo.IndexModel{
						Keys: bson.M{"own": 1},
					},
				}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		go func() {
			//file
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["file"]).Indexes()
			var err error
			if _, err = I.CreateMany(
				indexctx,
				[]mongo.IndexModel{
					mongo.IndexModel{
						Keys:    bson.M{"id": 1},
						Options: options.Index().SetUnique(true),
					},
					mongo.IndexModel{
						Keys: bson.M{"name": 1},
					},
					mongo.IndexModel{
						Keys: bson.M{"own": 1},
					},
				}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
			}
		}()
		go func() {
			//store
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["store"]).Indexes()
			var err error
			if _, err = I.CreateMany(
				indexctx,
				[]mongo.IndexModel{
					mongo.IndexModel{
						Keys:    bson.M{"id": 1},
						Options: options.Index().SetUnique(true),
					},
					mongo.IndexModel{
						Keys: bson.M{"id.hash": 1},
					},
				}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		go func() {
			//chunk
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["chunk"]).Indexes()
			var err error
			if _, err = I.CreateOne(indexctx, mongo.IndexModel{
				Keys:    bson.D{bson.E{Key: "id", Value: 1}, bson.E{Key: "idx", Value: 1}},
				Options: options.Index().SetUnique(true),
			}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		go func() {
			//tag
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["filetags"]).Indexes()
			var err error
			if _, err = I.CreateMany(
				indexctx,
				[]mongo.IndexModel{
					mongo.IndexModel{
						Keys: bson.D{
							bson.E{Key: "owner", Value: 1},
							bson.E{Key: "word", Value: 1},
							bson.E{Key: "file", Value: 1},
						},
						Options: options.Index().SetUnique(true),
					},
					mongo.IndexModel{
						Keys: bson.M{"word": 1},
					},
				}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		go func() {
			//tag
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["storetags"]).Indexes()
			var err error
			if _, err = I.CreateMany(
				indexctx,
				[]mongo.IndexModel{
					mongo.IndexModel{
						Keys:    bson.D{bson.E{Key: "store", Value: 1}, bson.E{Key: "word", Value: 1}},
						Options: options.Index().SetUnique(true),
					},
					mongo.IndexModel{
						Keys: bson.M{"word": 1},
					},
				}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		go func() {
			//group
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["lines"]).Indexes()
			var err error
			if _, err = I.CreateOne(indexctx, mongo.IndexModel{
				Keys:    bson.D{bson.E{Key: "id", Value: 1}, bson.E{Key: "position", Value: 1}},
				Options: options.Index().SetUnique(true),
			}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		go func() {
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["acronym"]).Indexes()
			var err error
			if _, err = I.CreateOne(indexctx, mongo.IndexModel{
				Keys:    bson.D{bson.E{Key: "acronym", Value: 1}, bson.E{Key: "complete", Value: 1}},
				Options: options.Index().SetUnique(true),
			}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		go func() {
			defer wg.Done()
			I := testclient.Database(d.DBName).Collection(d.CollNames["reset"]).Indexes()
			var err error
			if _, err = I.CreateMany(indexctx, []mongo.IndexModel{
				mongo.IndexModel{
					Keys:    bson.M{"user": 1},
					Options: options.Index().SetUnique(true),
				},
				mongo.IndexModel{
					Keys:    bson.M{"key": 1},
					Options: options.Index().SetUnique(true),
				},
				mongo.IndexModel{
					Keys:    bson.M{"expire": 1},
					Options: options.Index().SetExpireAfterSeconds(0),
				},
			}); err != nil {
				select {
				case cherr <- err:
				case <-indexctx.Done():
				}
				return
			}
		}()
		wg.Wait()
		cherr <- nil
		err := <-cherr
		if err != nil {
			return err
		}
	}

	return testclient.Disconnect(ctx)
}

func initcoll(c map[string]string) map[string]string {
	if c == nil {
		c = make(map[string]string)
	}
	if _, ok := c["user"]; !ok {
		c["user"] = "user"
	}
	if _, ok := c["group"]; !ok {
		c["group"] = "group"
	}
	if _, ok := c["file"]; !ok {
		c["file"] = "file"
	}
	if _, ok := c["store"]; !ok {
		c["store"] = "store"
	}
	if _, ok := c["lines"]; !ok {
		c["lines"] = "lines"
	}
	if _, ok := c["chunk"]; !ok {
		c["chunk"] = "chunk"
	}
	if _, ok := c["filetags"]; !ok {
		c["filetags"] = "filetags"
	}
	if _, ok := c["storetags"]; !ok {
		c["storetags"] = "storetags"
	}
	if _, ok := c["acronym"]; !ok {
		c["acronym"] = "acronym"
	}
	if _, ok := c["reset"]; !ok {
		c["reset"] = "reset"
	}
	if _, ok := c["view"]; !ok {
		c["view"] = "view"
	}
	return c
}

func (d *Database) inittracking(buildform *Database) {
	if d.gotten == nil {
		if buildform.gotten == nil {
			d.trackOwners = newTrackOwners()
		} else {
			d.trackOwners = buildform.trackOwners
		}
	}
}

func (d *Database) initclient(c context.Context) {
	if c != nil {
		var err error
		d.ctx, d.cancel = context.WithCancel(c)
		d.client, err = mongo.Connect(d.ctx, options.Client().ApplyURI(d.URI))
		if err != nil {
			panic(err)
		}
	}
}

// Owner opens a new connection to the database if provided a context and returns Ownerbase type
// if provided context is nil the resulting Ownerbase will reuse the existing connection
func (d *Database) Owner(c context.Context) database.Ownerbase {
	n := new(Ownerbase)
	n.Database = *d
	n.inittracking(d)
	n.initclient(c)
	return n
}

// File opens a new connection to the database if provided a context and returns Filebase type
// if provided context is nil it will reuse the existing connection
func (d *Database) File(c context.Context) database.Filebase {
	n := new(Filebase)
	n.Database = *d
	n.inittracking(d)
	n.initclient(c)
	return n
}

// Store opens a new connection to the database if provided a context and returns Storebase type
// if provided context is nil it will reuse the existing connection
func (d *Database) Store(c context.Context) database.Storebase {
	n := new(Storebase)
	n.Database = *d
	n.initclient(c)
	return n
}

// Content opens a new connection to the database if provided a context and returns Contentbase type
// if provided context is nil it will reuse the existing connection
func (d *Database) Content(c context.Context) database.Contentbase {
	n := new(Contentbase)
	n.Database = *d
	n.initclient(c)
	return n
}

// Tag opens a new connection to the database if provided a context and returns Tagbase type
// if provided context is nil it will reuse the existing connection
func (d *Database) Tag(c context.Context) database.Tagbase {
	n := new(Tagbase)
	n.Database = *d
	n.initclient(c)
	return n
}

// Acronym opens a new connection to the database if provided a context and returns Acronymbase type
// if provided context is nil it will reuse the existing connection
func (d *Database) Acronym(c context.Context) database.Acronymbase {
	n := new(Acronymbase)
	n.Database = *d
	n.initclient(c)
	return n
}

// View opens a new connection to the database if provided a context and returns Viewbase type
// if provided context is nil it will reuse the existing connection
func (d *Database) View(c context.Context) database.Viewbase {
	n := new(Viewbase)
	n.Database = *d
	n.initclient(c)
	return n
}

// Close stops the active connection necessary or else there can be memory leak from unclosed connections
func (d *Database) Close(ctx context.Context) error {
	defer func() {
		d.client = nil
		d.cancel = nil
		d.ctx = nil
	}()
	if d.cancel != nil {
		d.cancel()
	}
	if ctx == nil {
		ctx = d.ctx
	}
	if d.client != nil {
		if err := d.client.Disconnect(ctx); err != nil {
			return srverror.New(err, 500, "Database Error 101", "unable to close connection to Database")
		}
	}
	return nil
}

// GetContext returns context of the current open connection
func (d *Database) GetContext() context.Context {
	return d.ctx
}
