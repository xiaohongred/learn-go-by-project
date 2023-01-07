package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	if userType != role {
		return errors.New("unauthorized to access this resource")
	}
	return nil
}

func MatchUserTypeToUid(c *gin.Context, userId string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	if userType == "USER" && uid != userId {
		return errors.New("unauthorized to access this resource")
	}
	err := CheckUserType(c, userType)
	return err
}
