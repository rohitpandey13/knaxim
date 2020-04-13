package mongo

import (
	"fmt"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Viewbase is a connection to the database with pdf view operations
type Viewbase struct {
	Database
}

// Insert adds pdf view to the database
func (vb *Viewbase) Insert(vs *types.ViewStore) error {
	chunks := chunkify(vs.ID, vs.Content)
	_, err := vb.client.Database(vb.DBName).Collection(vb.CollNames["view"]).InsertMany(
		vb.ctx,
		chunks,
		options.InsertMany().SetOrdered(false),
	)
	if err != nil {
		return srverror.New(err, 500, "Database Error V3", "unable to insert viewstore chunks")
	}
	return nil
}

// Get view from database
func (vb *Viewbase) Get(id types.StoreID) (out *types.ViewStore, err error) {
	var chunks []*contentchunk
	cursor, err := vb.client.Database(vb.DBName).Collection(vb.CollNames["view"]).Find(
		vb.ctx,
		bson.M{"id": id},
	)
	if err != nil {
		return nil, srverror.New(err, 500, "Database Error V4", "failed to get view data chunks")
	}
	if err = cursor.All(vb.ctx, &chunks); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound.Extend("no View", id.String())
		}
		return nil, srverror.New(err, 500, "Database Error V5", "failed to decode view chunks")
	}
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = srverror.New(v, 500, "Database Error V6", "unable to build chunks")
				out = nil
			default:
				err = srverror.New(fmt.Errorf("GetStore: %+#v", v), 500, "Database Error V7")
				out = nil
			}
		}
	}()
	out = new(types.ViewStore)
	out.Content = appendchunks(chunksort(chunks))
	out.ID = id
	return
}
