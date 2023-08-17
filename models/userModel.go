package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"firstName" validate:"required,min=2,max=100"`
	Last_name     *string            `json:"lastName" validate:"required,min=2,max=100"`
	Password      string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token"`
	User_type     string            `json:"userType"`
	Refresh_token *string            `json:"refreshToken"`
	Created_at    time.Time          `json:"createdAt"`
	Updated_at    time.Time          `json:"updatedAt"`
	User_id       string             `json:"userId"`
}

type Availability struct {
	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime" `
	Date      time.Time `json:"date" bson:"date"`
}

type Ticket struct {
	ID         int       `json:"id" bson:"id"`
	StartTime  time.Time `json:"startTime" bson:"startTime" `
	EndTime    time.Time `json:"endTime" bson:"endTime" `
	Date       time.Time `json:"date" bson:"date"`
	IsAssigned bool      `bson:"isAssigned" json:"isAssigned"`
}

type Reservation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TicketID  int                `bson:"ticketId" json:"ticketId"`
	UserID    string             `bson:"userId" json:"userId"`
	StartTime time.Time          `bson:"startTime" json:"startTime"`
	EndTime   time.Time          `bson:"endTime" json:"endTime"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	Phone     string             `bson:"phone" json:"phone"`
}
