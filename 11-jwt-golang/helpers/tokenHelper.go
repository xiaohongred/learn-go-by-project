package helpers

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"jwt-golang/database"
	"log"
	"os"
	"time"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstname string, lastname string,
	userType string, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstname,
		Last_name:  lastname,
		Uid:        userId,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).
		SignedString(key)
	if err != nil {
		log.Panic(err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).
		SignedString(key)
	if err != nil {
		log.Panic(err)
	}
	return token, refreshToken, nil
}

func UpdateAllTokens(token string, refreshToken string, userId string) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelFunc()

	var updatObj primitive.D
	updatObj = append(updatObj, bson.E{"token", token})
	updatObj = append(updatObj, bson.E{"refresh_token", refreshToken})
	Update_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatObj = append(updatObj, bson.E{"updated_at", Update_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{"$set", updatObj}}, &opt)
	if err != nil {
		log.Panic(err)
	}
	return
}

func ValidateToken(clientToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(clientToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		return
	}
	return claims, msg
}
