package helpers

import "ecommerce/pkg/models"

func SterializeUser(user models.User) models.User{
	var returnedUser models.User

	returnedUser.ID = user.ID
	returnedUser.Address = user.Address
	returnedUser.CreatedAt= user.CreatedAt
	returnedUser.UpdatedAt = user.UpdatedAt
	returnedUser.Email = user.Email
	returnedUser.Username = user.Username
	returnedUser.Orders = user.Orders
	returnedUser.UserCart = user.UserCart
	returnedUser.UserType = user.UserType

	return returnedUser
}