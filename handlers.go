package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	user.ID = uuid.New().String()

	DB.Create(&user)

	MongoUser.InsertOne(ctx, user)

	DB.Create(&AuditLog{
		Entity:   "USER",
		EntityID: user.ID,
		Action:   "CREATE",
	})

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	val, err := RDB.Get(ctx, id).Result()
	if err == nil {
		var cachedUser User
		json.Unmarshal([]byte(val), &cachedUser)
		c.JSON(http.StatusOK, cachedUser)
		return
	}

	var user User
	DB.First(&user, "id = ?", id)

	data, _ := json.Marshal(user)
	RDB.Set(ctx, id, data, 60*time.Second)

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var input User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Age = input.Age
	DB.Save(&user)

	MongoUser.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{
			"$set": bson.M{
				"name":  user.Name,
				"email": user.Email,
				"age":   user.Age,
			},
		},
	)

	RDB.Del(ctx, id)

	DB.Create(&AuditLog{
		Entity:   "USER",
		EntityID: id,
		Action:   "UPDATE",
	})

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	DB.Delete(&Task{}, "user_id = ?", id)

	if err := DB.Delete(&User{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	MongoUser.DeleteOne(ctx, bson.M{"id": id})

	RDB.Del(ctx, id)

	DB.Create(&AuditLog{
		Entity:   "USER",
		EntityID: id,
		Action:   "DELETE",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
