package helper

import (
	"errors"

	"github.com/Shaheer25/go-auth/models"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
func ConvertToInterfaceSlice(tickets []models.Ticket) []interface{} {
	result := make([]interface{}, len(tickets))
	for i, ticket := range tickets {
		result[i] = ticket
	}
	return result
}
