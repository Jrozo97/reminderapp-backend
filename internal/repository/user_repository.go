package repository

import (
	"context"
	"os"
	"time"

	"github.com/Jrozo97/reminderapp-backend/internal/config"
	"github.com/Jrozo97/reminderapp-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRepository {

	nameDb := os.Getenv("MONGO_DB")
	db := config.MongoClient.Database(nameDb)
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now().Unix()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
