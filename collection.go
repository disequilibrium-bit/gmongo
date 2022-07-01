package gmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

// Collection is a collection from the MongoDB deployment.
type Collection struct {
	col *mongo.Collection
}

// InsertOne Saves the document into the Collection.
// The document parameter must be defined by tag bson.
func (c *Collection) InsertOne(document any) (*primitive.ObjectID, error) {
	result, err := c.col.InsertOne(context.Background(), &document)
	if err != nil {
		return nil, err
	}
	objectId := result.InsertedID.(primitive.ObjectID)
	return &objectId, nil
}

// UpdateOne updates the document into the Collection.
// The filter parameter must be a document containing query operators and can be used to select the document to be
// updated. It cannot be nil.
// The update parameter must be defined by tag bson and the data to be updated is stored in the struct.
// It cannot be nil.
// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
func (c *Collection) UpdateOne(filter, update any, opts ...*options.UpdateOptions) error {
	vfMap := make(map[string]any, 0)
	obtainValidField(update, "", vfMap)
	update = bson.M{"$set": vfMap}
	_, err := c.col.UpdateOne(context.Background(), filter, update, opts...)
	if err != nil {
		return err
	}
	return nil
}

const bs = "bson"

func obtainValidField(material any, parentKey string, upMap map[string]any) {
	v := reflect.ValueOf(material)
	t := reflect.TypeOf(material)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			ve := v.Field(i)
			te := t.Field(i)

			kind := ve.Type().Kind()
			switch kind {
			case reflect.Struct:
				if !ve.IsNil() {
					pk := parentKey + te.Tag.Get(bs)
					structField := ve.Interface()
					obtainValidField(structField, pk, upMap)
				}
			case reflect.Slice:
				if !ve.IsNil() {
					pk := parentKey + "." + te.Tag.Get(bs) + ".$"
					for i := 0; i < ve.Len(); i++ {
						sliceField := ve.Index(i).Interface()
						obtainValidField(sliceField, pk, upMap)
					}
				}
			default:
				if !ve.IsZero() {
					pk := parentKey + "." + te.Tag.Get(bs)
					first := strings.Index(pk, ".")
					upMap[pk[first+1:]] = ve.Interface()
				}
			}
		}
	}
}
