package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shaheer25/go-auth/database"
	helper "github.com/Shaheer25/go-auth/helpers"
	"github.com/Shaheer25/go-auth/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var timeCollection *mongo.Collection = database.OpenCollection(database.Client, "Availability time")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while listing user items",
			})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func AvailabilityTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Check if the user is an ADMIN
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		var availability models.Availability
		if err := c.ShouldBindJSON(&availability); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()})
			return
		}

		validationErr := validate.Struct(availability)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		ticket := models.Availability{
			StartTime: availability.StartTime,
			EndTime:   availability.EndTime,
		}

		_, err := timeCollection.InsertOne(ctx, ticket)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create ticket",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "availability time added successfully",
		})
	}
}

var ticketCollection *mongo.Collection = database.OpenCollection(database.Client, "ticket")

func GenerateTickets() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var request struct {
			Date      time.Time `json:"date" binding:"required"`
			StartTime time.Time `json:"start_time" binding:"required"`
			EndTime   time.Time `json:"end_time" binding:"required"`
		}

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		existingTickets, err := ticketCollection.Find(ctx, bson.M{"date": request.Date})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to check existing tickets",
			})
			return
		}

		if existingTickets.Next(ctx) {
			c.JSON(http.StatusBadRequest, gin.H{
				"tickets": "Tickets already generated for the specified date",
			})
			return
		}

		var tickets = []models.Ticket{}
		currentTime := request.StartTime
		for currentTime.Before(request.EndTime) {
			ticket := models.Ticket{
				ID:         len(tickets) + 1,
				Date:       request.Date,
				StartTime:  currentTime,
				EndTime:    currentTime.Add(15 * time.Minute),
				IsAssigned: false,
			}
			tickets = append(tickets, ticket)

			currentTime = currentTime.Add(15 * time.Minute)
		}

		_, err = ticketCollection.InsertMany(ctx, helper.ConvertToInterfaceSlice(tickets))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate tickets",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tickets": tickets,
		})
	}
}

func DeleteTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		var request struct {
			TicketIDs []int `json:"ticket_ids" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request format",
			})
			return
		}

		filter := bson.M{"id": bson.M{"$in": request.TicketIDs}}

		result, err := ticketCollection.DeleteMany(ctx, filter)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete tickets",
			})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "tickets not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("deleted %d tickets", result.DeletedCount),
		})
	}
}

func GetAllReservations() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := helper.GetUserTokenFromContext(c)
		claims, errMsg := helper.ValidateToken(userToken)
		if errMsg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			return
		}

		if claims.User_type != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		var reservations []models.Reservation
		cursor, err := reservationCollection.Find(context.Background(), bson.M{})
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
	}
}
