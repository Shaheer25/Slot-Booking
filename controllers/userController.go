package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Shaheer25/go-auth/database"
	helper "github.com/Shaheer25/go-auth/helpers"
	"github.com/Shaheer25/go-auth/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ShowTickets() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		cursor, err := ticketCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch tickets",
			})
			return
		}
		defer cursor.Close(ctx)

		var tickets []models.Ticket
		for cursor.Next(ctx) {
			var ticket models.Ticket
			err := cursor.Decode(&ticket)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to decode ticket",
				})
				return
			}
			tickets = append(tickets, ticket)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "cursor error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tickets": tickets,
		})
	}
}

var reservationCollection *mongo.Collection = database.OpenCollection(database.Client, "reservations")

func BookTickets() gin.HandlerFunc {
	return func(c *gin.Context) {

		ticketID := c.Param("id")

		ticketIDInt, err := strconv.Atoi(ticketID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ticket ID",
			})
			return
		}

		userToken := helper.GetUserTokenFromContext(c)

		claims, errMsg := helper.ValidateToken(userToken)
		if errMsg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMsg,
			})
			return
		}

		var user models.User
		err = userCollection.FindOne(context.Background(), bson.M{"user_id": claims.Uid}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to fetch user",
				})
			}
			return
		}

		var ticket models.Ticket
		err = ticketCollection.FindOne(context.Background(), bson.M{"id": ticketIDInt, "isAssigned": false}).Decode(&ticket)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Ticket not found or already assigned",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to fetch ticket",
				})
			}
			return
		}

		reservation := models.Reservation{
			ID:        primitive.NewObjectID(),
			TicketID:  ticketIDInt,
			UserID:    claims.Uid,
			StartTime: ticket.StartTime,
			EndTime:   ticket.EndTime,
			FirstName: claims.First_name,
			LastName:  claims.Last_name,
			Phone:     claims.Phone,
		}

		_, err = reservationCollection.InsertOne(context.Background(), reservation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to make reservation",
			})
			return
		}

		// Mark the ticket as assigned
		_, err = ticketCollection.UpdateOne(context.Background(), bson.M{"id": ticketIDInt}, bson.M{"$set": bson.M{"isAssigned": true}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to mark ticket as assigned",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Success": "Reservation Done Successfully",
		})
	}
}

func GetUserReservations() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := helper.GetUserTokenFromContext(c)
		claims, errMsg := helper.ValidateToken(userToken)
		if errMsg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMsg,
			})
			return
		}

		var reservations []models.Reservation
		cursor, err := reservationCollection.Find(context.Background(), bson.M{"userId": claims.Uid})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch reservations",
			})
			return
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var reservation models.Reservation
			if err := cursor.Decode(&reservation); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to decode reservation",
				})
				return
			}
			reservations = append(reservations, reservation)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Cursor error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reservations": reservations,
		})
		return
	}
}

func GetEmptySlots() gin.HandlerFunc {
	return func(c *gin.Context) {
		emptySlots := []models.Ticket{}

		filter := bson.M{"isAssigned": false}
		cursor, err := ticketCollection.Find(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch empty slots",
			})
			return
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var ticket models.Ticket
			if err := cursor.Decode(&ticket); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to decode ticket",
				})
				return
			}
			emptySlots = append(emptySlots, ticket)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Cursor error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"empty_slots": emptySlots,
		})
	}
}

func DeleteReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ticketID := c.Param("id")

		ticketIDInt, err := strconv.Atoi(ticketID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ticket ID",
			})
			return
		}

		userToken := helper.GetUserTokenFromContext(c)
		claims, errMsg := helper.ValidateToken(userToken)
		if errMsg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMsg,
			})
			return
		}

		_, err = ticketCollection.UpdateOne(context.Background(), bson.M{"id": ticketIDInt}, bson.M{"$set": bson.M{"isAssigned": false}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update ticket",
			})
			return
		}

		res, err := reservationCollection.DeleteOne(context.Background(), bson.M{"userId": claims.Uid, "ticketId": ticketIDInt})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete reservation",
			})
			return
		}

		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Reservation not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Reservation deleted successfully",
		})
	}
}
