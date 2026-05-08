// internal/repository/auth_repository.go
package repository

import (
	"context"

	"github.com/shingo/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	collection *mongo.Collection
}

func NewAdminRepository(db *mongo.Client) *AdminRepository {
	return &AdminRepository{
		collection: db.Database("nudgebuddy").Collection("admins"),
	}
}

func (r *AdminRepository) CreateAdmin(ctx context.Context, user *models.Admin) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *AdminRepository) FindByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var user models.Admin
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
