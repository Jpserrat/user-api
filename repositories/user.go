package repositories

import (
	"context"
	"log"
	"user-api/models"
	"user-api/response"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo interface {
	Save(u models.User) response.ApiError
	GetAll(limit uint64, page uint64) ([]models.User, error)
	FindByField(value interface{}, key string) (models.User, response.ApiError)
	FindById(id string) (models.User, response.ApiError)
	DeleteById(id string) response.ApiError
	UpdateByID(id string, u models.User) (apiErr response.ApiError)
}

type userMongoImpl struct {
	db  *mongo.Collection
	ctx context.Context
}

func NewUserMongo(mongoDb *mongo.Collection, ctx context.Context) UserRepo {
	return userMongoImpl{
		db:  mongoDb,
		ctx: ctx,
	}
}

func (m userMongoImpl) Save(u models.User) response.ApiError {
	_, err := m.db.InsertOne(m.ctx, u)
	if err != nil {
		log.Printf("[UserRepo] Error saving user %s", err.Error())
		return response.InternalServerError
	}
	log.Printf("[UserRepo] user inserted %v", u)
	return response.ApiError{}
}

func (m userMongoImpl) GetAll(limit uint64, page uint64) ([]models.User, error) {
	result := make([]models.User, 0)

	l := int64(limit)
	skip := int64(page*limit - limit)
	opt := options.FindOptions{Limit: &l, Skip: &skip}

	curr, err := m.db.Find(m.ctx, bson.D{}, &opt)
	if err != nil {
		return result, err
	}

	for curr.Next(m.ctx) {
		var el models.User
		if err := curr.Decode(&el); err != nil {
			log.Println(err)
			return nil, err
		}

		result = append(result, el)
	}

	log.Printf("[UserRepo] Users found %v", len(result))

	return result, nil
}

func (r userMongoImpl) FindByField(value interface{}, key string) (models.User, response.ApiError) {
	u := models.User{}
	err := r.db.FindOne(r.ctx, bson.D{{key, value}}, options.FindOne()).Decode(&u)

	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			log.Printf("[UserRepo] No document found with value %s and key %s", value, key)
			return u, response.ResourceNotFoundError
		}
		log.Printf("[UserRepo] Error getting document with value %s and key %s, error: %s", value, key, err.Error())
		return u, response.InternalServerError
	}
	return u, response.ApiError{}
}

func (r userMongoImpl) FindById(id string) (models.User, response.ApiError) {
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Printf("[UserRepo] Invalid id format %s", id)
		return models.User{}, response.BadRequestError
	}

	return r.FindByField(objID, "_id")
}

func (r userMongoImpl) DeleteById(id string) response.ApiError {
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Printf("[UserRepo] Invalid id format %s", id)
		return response.BadRequestError
	}

	_, err = r.db.DeleteOne(r.ctx, bson.D{{"_id", objID}}, options.Delete())

	if err != nil {
		log.Printf("[UserRepo] Unexpected error deleting user by id: %s", err.Error())
		return response.InternalServerError
	}

	return response.ApiError{}
}

func (r userMongoImpl) UpdateByID(id string, u models.User) (apiErr response.ApiError) {
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Printf("[UserRepo] Invalid id format %s", id)
		return response.BadRequestError
	}

	filter := bson.D{{"_id", objID}}
	update := bson.D{{"$set", bson.D{{"age", u.Age}, {"address", u.Address}, {"Name", u.Name}}}}

	_, err = r.db.UpdateOne(r.ctx, filter, update)

	if err != nil {
		log.Printf("[UserRepo] Error updating user document: %s", err.Error())
		return response.InternalServerError
	}

	return
}
