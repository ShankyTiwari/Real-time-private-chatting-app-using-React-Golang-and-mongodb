package handlers

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	config "private-chat/config"
	constants "private-chat/constants"
	utils "private-chat/utils"
)

// UpdateUserOnlineStatusByUserID will update the online status of the user
func UpdateUserOnlineStatusByUserID(userID string, status string) error {
	docID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil
	}

	collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_, queryError := collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.M{"$set": bson.M{"online": status}})
	defer cancel()

	if queryError != nil {
		return errors.New(constants.ServerFailedResponse)
	}
	return nil
}

// GetUserByUsername function will return user datails based username
func GetUserByUsername(username string) UserDetailsStruct {
	var userDetails UserDetailsStruct

	collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_ = collection.FindOne(ctx, bson.M{
		"username": username,
	}).Decode(&userDetails)

	defer cancel()

	return userDetails
}

// GetUserByUserID function will return user datails based username
func GetUserByUserID(userID string) UserDetailsStruct {
	var userDetails UserDetailsStruct

	docID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return UserDetailsStruct{}
	}

	collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_ = collection.FindOne(ctx, bson.M{
		"_id": docID,
	}).Decode(&userDetails)

	defer cancel()

	return userDetails
}

// IsUsernameAvailableQueryHandler function will check username from the database
func IsUsernameAvailableQueryHandler(username string) bool {
	userDetails := GetUserByUsername(username)
	if userDetails == (UserDetailsStruct{}) {
		return true
	}
	return false
}

// LoginQueryHandler function will check username from the database
func LoginQueryHandler(userDetailsRequestPayload UserDetailsRequestPayloadStruct) (UserDetailsResponsePayloadStruct, error) {
	if userDetailsRequestPayload.Username == "" {
		return UserDetailsResponsePayloadStruct{}, errors.New(constants.UsernameCantBeEmpty)
	} else if userDetailsRequestPayload.Password == "" {
		return UserDetailsResponsePayloadStruct{}, errors.New(constants.PasswordCantBeEmpty)
	} else {
		userDetails := GetUserByUsername(userDetailsRequestPayload.Username)
		if userDetails == (UserDetailsStruct{}) {
			return UserDetailsResponsePayloadStruct{}, errors.New(constants.UserIsNotRegisteredWithUs)
		}

		if isPasswordOkay := utils.ComparePasswords(userDetailsRequestPayload.Password, userDetails.Password); isPasswordOkay != nil {
			return UserDetailsResponsePayloadStruct{}, errors.New(constants.LoginPasswordIsInCorrect)
		}

		if onlineStatusError := UpdateUserOnlineStatusByUserID(userDetails.ID, "Y"); onlineStatusError != nil {
			return UserDetailsResponsePayloadStruct{}, errors.New(constants.LoginPasswordIsInCorrect)
		}

		return UserDetailsResponsePayloadStruct{
			UserID:   userDetails.ID,
			Username: userDetails.Username,
		}, nil
	}
}

// RegisterQueryHandler function will check username from the database
func RegisterQueryHandler(userDetailsRequestPayload UserDetailsRequestPayloadStruct) (string, error) {
	if userDetailsRequestPayload.Username == "" {
		return "", errors.New(constants.UsernameCantBeEmpty)
	} else if userDetailsRequestPayload.Password == "" {
		return "", errors.New(constants.PasswordCantBeEmpty)
	} else {
		newPasswordHash, newPasswordHashError := utils.CreatePassword(userDetailsRequestPayload.Password)
		if newPasswordHashError != nil {
			return "", errors.New(constants.ServerFailedResponse)
		}
		collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		registrationQueryResponse, registrationError := collection.InsertOne(ctx, bson.M{
			"username": userDetailsRequestPayload.Username,
			"password": newPasswordHash,
			"online":   "N",
		})
		defer cancel()

		registrationQueryObjectID, registrationQueryObjectIDError := registrationQueryResponse.InsertedID.(primitive.ObjectID)

		if onlineStatusError := UpdateUserOnlineStatusByUserID(registrationQueryObjectID.Hex(), "Y"); onlineStatusError != nil {
			return " ", errors.New(constants.ServerFailedResponse)
		}

		if registrationError != nil || !registrationQueryObjectIDError {
			return "", errors.New(constants.ServerFailedResponse)
		}

		return registrationQueryObjectID.Hex(), nil
	}
}

// GetAllOnlineUsers function will return the all online users
func GetAllOnlineUsers(userID string) []UserDetailsResponsePayloadStruct {
	var onlineUsers []UserDetailsResponsePayloadStruct

	docID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return onlineUsers
	}

	collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, queryError := collection.Find(ctx, bson.M{
		"online": "Y",
		"_id": bson.M{
			"$ne": docID,
		},
	})
	defer cancel()

	if queryError != nil {
		return onlineUsers
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var singleOnlineUser UserDetailsStruct
		err := cursor.Decode(&singleOnlineUser)

		if err == nil {
			onlineUsers = append(onlineUsers, UserDetailsResponsePayloadStruct{
				UserID:   singleOnlineUser.ID,
				Online:   singleOnlineUser.Online,
				Username: singleOnlineUser.Username,
			})
		}
	}

	return onlineUsers
}

// StoreNewChatMessages is used for storing a new message
func StoreNewChatMessages(messagePayload MessagePayloadStruct) bool {
	collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_, registrationError := collection.InsertOne(ctx, bson.M{
		"fromUserID": messagePayload.FromUserID,
		"message":    messagePayload.Message,
		"toUserID":   messagePayload.ToUserID,
	})
	defer cancel()

	if registrationError == nil {
		return false
	}
	return true
}

// GetConversationBetweenTwoUsers will be used to fetch the conversation between two users
func GetConversationBetweenTwoUsers(toUserID string, fromUserID string) []ConversationStruct {
	var conversations []ConversationStruct

	collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	queryCondition := bson.M{
		"$or": []bson.M{
			{
				"$and": []bson.M{
					{
						"toUserID": toUserID,
					},
					{
						"fromUserID": fromUserID,
					},
				},
			},
			{
				"$and": []bson.M{
					{
						"toUserID": fromUserID,
					},
					{
						"fromUserID": toUserID,
					},
				},
			},
		},
	}

	cursor, queryError := collection.Find(ctx, queryCondition)
	defer cancel()

	if queryError != nil {
		return conversations
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var conversation ConversationStruct
		err := cursor.Decode(&conversation)

		if err == nil {
			conversations = append(conversations, ConversationStruct{
				ID:         conversation.ID,
				FromUserID: conversation.FromUserID,
				ToUserID:   conversation.ToUserID,
				Message:    conversation.Message,
			})
		}
	}
	return conversations
}
