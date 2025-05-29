package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type UserAdmin struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Admin       bool      `json:"admin"`
	Active      bool      `json:"active"`
	Id          string    `json:"id"`
	Created     time.Time `json:"created"`
	LastLogin   time.Time `json:"last_login"`
}

type UserListAdmin struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

// AdminListUsers List users
//
//	@Summary		Lists the Ids of the users
//	@Description	List Users
//	@Tags			Admin
//	@Produce		json
//	@Success		200	{object}	[]UserListAdmin
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Security		jwt
//	@Security		apikey
//	@Router			/admin/users [get]
func (a *Api) AdminListUsers(ctx *gin.Context) {
	users, err := a.UserService.ListUsers()
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": fmt.Sprintf("failed to list users: %s", err)})
		return
	}

	userList := make([]UserListAdmin, 0)
	for _, u := range users {
		userList = append(userList, UserListAdmin{
			Username: u.Username,
			Id:       u.Id,
		})
	}

	ctx.JSON(200, userList)
}

// AdminGetUser Get a specific user
//
//	@Summary		Returns a user
//	@Description	Returns a user
//	@Tags			Admin
//	@Produce		json
//	@Param			id	path		string	true	"Id of the user"
//	@Success		200	{object}	UserAdmin
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Security		jwt
//	@Security		apikey
//	@Router			/admin/user/{id} [get]
func (a *Api) AdminGetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "no id provided"})
	}

	u, err := a.UserService.GetUserById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": fmt.Sprintf("failed to get user: %s", err)})
		return
	}

	ctx.JSON(200, &UserAdmin{
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Admin:       u.Admin,
		Id:          u.Id,
		Created:     u.Created,
		LastLogin:   u.LastLogin,
		Active:      u.Active,
	})
}
