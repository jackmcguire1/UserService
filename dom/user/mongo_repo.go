package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackmcguire1/UserService/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	BaseRepository

	Collection *mongo.Collection
}

type MongoRepoParams struct {
	Host           string
	Database       string
	CollectionName string
}

func NewMongoRepo(ctx context.Context, params *MongoRepoParams) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(params.Host))
	if err != nil {
		return nil, err
	}
	database := client.Database(params.Database)

	collection := database.Collection(params.CollectionName)

	return &MongoRepository{Collection: collection}, nil
}

func (repo *MongoRepository) GetUser(userId string) (*User, error) {
	return repo.GetUserByAttr("_id", userId)
}

func (repo *MongoRepository) GetUserByEmail(email string) (*User, error) {
	return repo.GetUserByAttr("email", email)
}

func (repo *MongoRepository) GetUserByAttr(attr, value string) (*User, error) {
	filter := bson.M{attr: value}
	res := repo.Collection.FindOne(context.Background(), filter, nil)
	if res.Err() != nil {
		if strings.Contains(res.Err().Error(), "no documents in result") {
			return nil, utils.ErrNotFound
		}
		return nil, res.Err()
	}

	var user *User
	err := res.Decode(&user)
	if err != nil {
		err = fmt.Errorf("failed to umarshal bson user document err:%w", err)
		return nil, err
	}

	return user, nil
}

func (repo *MongoRepository) GetUsersByCountry(cc string) ([]*User, error) {
	filter := bson.M{"countryCode": cc}
	return repo.searchUsers(filter)
}

func (repo *MongoRepository) PutUser(u *User) error {

	filter := bson.M{"_id": u.ID}

	data, err := bson.Marshal(u)
	if err != nil {
		return err
	}
	opts := options.Replace().SetUpsert(true) // Set upsert options if needed

	//repo.Collection.ReplaceOne(context.Background(), filter, data)
	_, err = repo.Collection.ReplaceOne(context.Background(), filter, data, opts)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoRepository) DeleteUser(id string) error {
	filter := bson.M{"_id": id}
	res, err := repo.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if res.DeletedCount != 1 {
		err = fmt.Errorf("failed to remove user from repo count:%d %w", res.DeletedCount, utils.ErrNotFound)
		return err
	}

	return nil
}

func (repo *MongoRepository) GetAllUsers() ([]*User, error) {
	return repo.searchUsers(bson.M{})
}

func (repo *MongoRepository) searchUsers(filter bson.M) ([]*User, error) {
	cursor, err := repo.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if cursor.Err() != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users = []*User{}
	for cursor.Next(context.Background()) {
		var tmpUsers []*User
		err = cursor.All(context.Background(), &users)
		if err != nil {
			return nil, err
		}
		users = append(users, tmpUsers...)
	}

	return users, err
}
