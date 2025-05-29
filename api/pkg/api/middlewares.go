package api

import (
	"github.com/gin-gonic/gin"
	"github.com/thomas-maurice/api/go-vue/pkg/services/userservice"
)

func (a *Api) DebugCORSMiddleware(ctx *gin.Context) {
	if a.Debug {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,X-AUTH-TOKEN,X-API-KEY")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
	}

	ctx.Next()
}

func (a *Api) UserTokenMiddleware(ctx *gin.Context) {
	token := ctx.Request.Header.Get("X-AUTH-TOKEN")
	if token == "" {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})
		return
	}

	session, user, err := a.UserService.VerifySessionToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})
		return
	}

	ctx.Set("user", user)
	ctx.Set("session", session)

	ctx.Next()
}

// RequiresUserLogin checks the user is logged in, and optionally
// if it is an admin of the instance. The logic here will check if
// there is an authentication JWT for the frontend first, and then
// check if the user provided an api key. In both cases the `user`
// key will be set on the context, and if the user is loggin in with
// a jwt, the `session` key will be populated too.
func (a *Api) RequiresUserLogin(requiresAdmin bool) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var user *userservice.User
		var session *userservice.Session
		var err error

		token := ctx.Request.Header.Get("X-AUTH-TOKEN")
		if token != "" {
			session, user, err = a.UserService.VerifySessionToken(token)
			if err != nil {
				ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})
				return
			}

			if requiresAdmin && !user.Admin {
				ctx.AbortWithStatusJSON(401, gin.H{"error": "access denied"})
				return
			}

			ctx.Set("user", user)
			ctx.Set("session", session)
			ctx.Next()
			return
		} else {
			token := ctx.Request.Header.Get("X-API-KEY")
			if token != "" {
				session, user, err = a.UserService.VerifySessionToken(token)
				if err != nil {
					ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})
					return
				}

				if requiresAdmin && !user.Admin {
					ctx.AbortWithStatusJSON(401, gin.H{"error": "access denied"})
					return
				}

				ctx.Set("user", user)
				ctx.Set("session", session)
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})
		return
	}
}
