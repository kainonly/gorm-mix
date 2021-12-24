package api

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Db *mongo.Database
}

func (x *Service) Create(ctx context.Context, name string, doc interface{}) (*mongo.InsertOneResult, error) {
	if data, ok := doc.(bson.M); ok {
		data["create_time"] = time.Now()
		data["update_time"] = time.Now()
		return x.Db.Collection(name).InsertOne(ctx, data)
	}
	return x.Db.Collection(name).InsertOne(ctx, doc)
}

func (x *Service) Find(
	ctx context.Context,
	name string,
	filter bson.M,
	sort []string,
	opts ...*options.FindOptions,
) (data []map[string]interface{}, err error) {
	option := options.Find()
	if len(sort) != 0 {
		sorts := make(bson.D, len(sort))
		for i, x := range sort {
			v := strings.Split(x, ",")
			var direction int
			if direction, err = strconv.Atoi(v[1]); err != nil {
				return
			}
			sorts[i] = bson.E{Key: v[0], Value: direction}
		}
		option.SetSort(sorts)
		option.SetAllowDiskUse(true)
	} else {
		option.SetSort(bson.M{"_id": -1})
	}
	opts = append(opts, option)
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection(name).
		Find(ctx, filter, opts...); err != nil {
		return
	}
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	return
}

func (x *Service) FindById(
	ctx context.Context,
	name string,
	id []string,
	sort []string,
) (data []map[string]interface{}, err error) {
	ids := make([]primitive.ObjectID, len(id))
	for i, v := range id {
		if ids[i], err = primitive.ObjectIDFromHex(v); err != nil {
			return
		}
	}
	return x.Find(ctx, name, bson.M{"_id": bson.M{"$in": ids}}, sort)
}

type FindResult struct {
	Total int64                    `json:"total"`
	Data  []map[string]interface{} `json:"data"`
}

func (x *Service) FindByPage(
	ctx context.Context,
	name string,
	page PaginationDto,
	filter bson.M,
	sort []string,
) (result FindResult, err error) {
	if len(filter) != 0 {
		if result.Total, err = x.Db.Collection(name).
			CountDocuments(ctx, filter); err != nil {
			return
		}
	} else {
		if result.Total, err = x.Db.Collection(name).
			EstimatedDocumentCount(ctx); err != nil {
			return
		}
	}
	option := options.Find()
	option.SetLimit(page.Size)
	option.SetSkip((page.Index - 1) * page.Size)
	if result.Data, err = x.Find(ctx, name, filter, sort, option); err != nil {
		return
	}
	return
}

func (x *Service) FindOne(ctx context.Context, name string, filter bson.M) (data map[string]interface{}, err error) {
	if err = x.Db.Collection(name).FindOne(ctx, filter).Decode(&data); err != nil {
		return
	}
	return
}

func (x *Service) FindOneById(ctx context.Context, name string, id string) (data map[string]interface{}, err error) {
	var objectId primitive.ObjectID
	if objectId, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	return x.FindOne(ctx, name, bson.M{"_id": objectId})
}

func (x *Service) UpdateMany(ctx context.Context,
	name string,
	filter bson.M,
	update interface{},
) (result *mongo.UpdateResult, err error) {
	if data, ok := update.(bson.M); ok {
		data["$set"].(map[string]interface{})["update_time"] = time.Now()
		return x.Db.Collection(name).UpdateMany(ctx, filter, data)
	}
	return x.Db.Collection(name).UpdateMany(ctx, filter, update)
}

func (x *Service) UpdateManyById(
	ctx context.Context,
	name string,
	id []string,
	update interface{},
) (result *mongo.UpdateResult, err error) {
	ids := make([]primitive.ObjectID, len(id))
	for i, v := range id {
		if ids[i], err = primitive.ObjectIDFromHex(v); err != nil {
			return
		}
	}
	return x.UpdateMany(ctx, name, bson.M{"_id": bson.M{"$in": ids}}, update)
}

func (x *Service) UpdateOne(
	ctx context.Context,
	name string,
	filter bson.M,
	update interface{},
) (result *mongo.UpdateResult, err error) {
	if data, ok := update.(bson.M); ok {
		data["$set"].(map[string]interface{})["update_time"] = time.Now()
		return x.Db.Collection(name).UpdateOne(ctx, filter, data)
	}
	return x.Db.Collection(name).UpdateOne(ctx, filter, update)
}

func (x *Service) UpdateOneById(
	ctx context.Context,
	name string,
	id string,
	update interface{},
) (result *mongo.UpdateResult, err error) {
	var objectId primitive.ObjectID
	if objectId, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	return x.UpdateOne(ctx, name, bson.M{"_id": objectId}, update)
}

func (x *Service) ReplaceOneById(
	ctx context.Context,
	name string,
	id string,
	doc interface{},
) (result *mongo.UpdateResult, err error) {
	var objectId primitive.ObjectID
	if objectId, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	filter := bson.M{"_id": objectId}
	if data, ok := doc.(bson.M); ok {
		data["create_time"] = time.Now()
		data["update_time"] = time.Now()
		return x.Db.Collection(name).ReplaceOne(ctx, filter, data)
	}
	return x.Db.Collection(name).ReplaceOne(ctx, filter, doc)
}

func (x *Service) DeleteOneById(ctx context.Context, name string, id string) (result *mongo.DeleteResult, err error) {
	var objectId primitive.ObjectID
	if objectId, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	return x.Db.Collection(name).DeleteOne(ctx, bson.M{"_id": objectId})
}