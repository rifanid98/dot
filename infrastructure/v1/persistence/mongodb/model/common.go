package model

import "go.mongodb.org/mongo-driver/bson/primitive"

func GetObjectId(id string) primitive.ObjectID {
	objId, _ := primitive.ObjectIDFromHex(id)
	return objId
}

func GetObjectIdHex(objId primitive.ObjectID) string {
	return objId.Hex()
}
