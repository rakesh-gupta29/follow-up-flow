// internal/repository/admin_repository.go
package repository

import (
	"context"
	"fmt"

	"github.com/shingo/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	collection *mongo.Collection
}

func NewAdminRepository(db *mongo.Client) *AdminRepository {
	return &AdminRepository{
		collection: db.Database("nudgebuddy_db").Collection("admins"),
	}
}

func (r *AdminRepository) GetAdmin(ctx context.Context, email string) (*models.Admin, error) {
	var user models.Admin
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	fmt.Println("admin foiudn", user, email)
	return &user, err
}
