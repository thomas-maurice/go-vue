package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thomas-maurice/api/go-vue/pkg/services/userservice"
)

type ProfileOutput struct {
	Id          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name,omitempty"`
	Email       string    `json:"email,omitempty"`
	Kind        string    `json:"kind"`
	Admin       bool      `json:"admin"`
	Created     time.Time `json:"created"`
	LastLogin   time.Time `json:"last_login"`
}

// ProfileSelf returns the profile of a user
//
//	@Summary		Returns the profile of a user
//	@Description	Returns the profile of a user
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	ProfileOutput
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Security		jwt
//	@Security		apikey
//	@Router			/user/profile [get]
func (a *Api) ProfileSelf(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "unknown user"})
		return
	}

	self, ok := user.(*userservice.User)
	if !ok {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "unknown user"})
		return
	}

	ctx.JSON(200, ProfileOutput{
		Id:          self.Id,
		Admin:       self.Admin,
		Username:    self.Username,
		DisplayName: self.DisplayName,
		Email:       self.Username,
		Kind:        string(self.Kind),
		Created:     self.Created,
		LastLogin:   self.LastLogin,
	})
}
