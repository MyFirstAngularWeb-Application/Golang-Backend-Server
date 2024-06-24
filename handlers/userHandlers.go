// handlers/userHandler.go
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend-server/config"
	"backend-server/models"
)

// CreateUser creates a new user
// func CreateUser(c *fiber.Ctx) error {
// 	var user struct {
// 		ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
// 		Name     string             `json:"fullName" bson:"fullName"`
// 		Email    string             `json:"email" bson:"email"`
// 		Password string             `json:"password" bson:"password"`
// 		// ConfirmPassword string             `json:"confirmPassword" bson:"-"`
// 	}

// 	// Parse the request body
// 	if err := c.BodyParser(&user); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
// 	}

// 	// Validate passwords
// 	// if user.Password != user.ConfirmPassword {
// 	// 	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Passwords do not match"})
// 	// }

// 	// Hash the password before storing
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
// 	}
// 	user.Password = string(hashedPassword)
// 	user.ID = primitive.NewObjectID()

// 	collection := config.DB.Collection("users")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_, err = collection.InsertOne(ctx, user)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
// 	}

//		return c.Status(http.StatusCreated).JSON(user)
//	}
func CreateUser(c *fiber.Ctx) error {
	var user struct {
		ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		Name     string             `json:"fullName" bson:"fullName"`
		Email    string             `json:"email" bson:"email"`
		Password string             `json:"password" bson:"password"`
	}

	// Parse the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	user.ID = primitive.NewObjectID()

	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(http.StatusCreated).JSON(user)
}

// GetUser handles fetching a user by email
func GetUser(c *fiber.Ctx) error {
	collection := config.DB.Collection("users")
	email := c.Params("email")

	// Validate email
	if email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	// Define the user model
	var user struct {
		ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		Name     string             `json:"fullName" bson:"fullName"`
		Email    string             `json:"email" bson:"email"`
		Password string             `json:"password" bson:"password"`
	}

	// Find the user by email
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user"})
	}

	// Return the user details
	return c.JSON(fiber.Map{
		"id":       user.ID.Hex(), // Convert ObjectID to string
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password, // Return the hashed password
	})
}

// UpdateUser updates a user's details
func UpdateUser(c *fiber.Ctx) error {
	collection := config.DB.Collection("users")
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": user})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(user)
}

// DeleteUser deletes a user by ID
func DeleteUser(c *fiber.Ctx) error {
	collection := config.DB.Collection("users")
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.SendStatus(http.StatusNoContent)
}

// GetUsers retrieves all users and returns their name, email, and ID
func GetUsers(c *fiber.Ctx) error {
	collection := config.DB.Collection("users")
	cur, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}
	defer cur.Close(context.TODO())

	var users []struct {
		ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		FullName string             `json:"fullName" bson:"fullName"`
		Email    string             `json:"email" bson:"email"`
	}
	for cur.Next(context.TODO()) {
		var user struct {
			ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
			FullName string             `json:"fullName" bson:"fullName"`
			Email    string             `json:"email" bson:"email"`
		}
		err := cur.Decode(&user)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode user"})
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cursor error"})
	}

	return c.JSON(users)
}
