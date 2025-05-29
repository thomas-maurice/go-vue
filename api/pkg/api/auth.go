package api

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thomas-maurice/api/go-vue/pkg/services/configservice"
	"github.com/thomas-maurice/api/go-vue/pkg/services/userservice"
	"golang.org/x/oauth2"
)

type Error struct {
	Error string `json:"error"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutOutput struct {
	Ok bool `json:"ok"`
}

type OIDCURLOutput struct {
	Url string `json:"url"`
}

type OIDCCallbackOutput struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type OIDCProvider struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type NewOIDCProvider struct {
	Name         string   `json:"name"`
	DisplayName  string   `json:"display_name"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Issuer       string   `json:"issuer"`
	Scopes       []string `json:"scopes"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

// AuthPassword
//
//	@Summary		Logs a local user in
//	@Description	Logs a local user in
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginInput	true	"Input data for the login"
//	@Success		200		{object}	LoginOutput
//	@Failure		400		{object}	Error
//	@Failure		500		{object}	Error
//	@Router			/auth/login [post]
func (a *Api) AuthPassword(ctx *gin.Context) {
	var login LoginInput

	if err := ctx.BindJSON(&login); err != nil {
		ctx.JSON(400, gin.H{"error": fmt.Errorf("could not bind the request to object: %w", err).Error()})
		return
	}

	var user *userservice.User
	var err error
	if user, err = a.UserService.Authenticate(login.Username, login.Password); err == nil {
		token, err := a.UserService.GenerateSessionToken(user)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "failed to generate session token"})
		}
		ctx.JSON(200, &LoginOutput{
			Token: token,
		})
		return
	}

	ctx.JSON(401, gin.H{"error": "invalid credentials"})
}

// GenerateOIDCRedirectURL
//
//	@Summary		Generates a redirect oidc login url
//	@Description	Generates a redirect oidc login url
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			name	path		int	true	"Name of the provider"
//	@Success		200		{object}	OIDCURLOutput
//	@Failure		400		{object}	Error
//	@Failure		500		{object}	Error
//	@Router			/auth/oidc/{name} [get]
func (a *Api) GenerateOIDCRedirectURL(ctx *gin.Context) {
	if a.Config.Security.OIDC == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}

	provider, err := a.ConfigService.GetOIDCProvider(ctx.Param("provider"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	state := uuid.NewString()
	ctx.SetCookie("oidc-state", state, 600, "", strings.Split(ctx.Request.Host, ":")[0], ctx.Request.TLS != nil, true)

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	redirect := fmt.Sprintf("%s://%s/auth/callback/%s", scheme, ctx.Request.Host, provider.Name)
	if os.Getenv("OIDC_REDIRECT_BASE_URL") != "" {
		redirect = fmt.Sprintf("%s/auth/callback/%s", os.Getenv("OIDC_REDIRECT_BASE_URL"), provider.Name)
	}

	prv, err := oidc.NewProvider(ctx, provider.Issuer)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
	}

	oauthConfig := oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  redirect,
		Endpoint:     prv.Endpoint(),
		Scopes:       provider.Scopes,
	}

	ctx.JSON(200, &OIDCURLOutput{Url: oauthConfig.AuthCodeURL(state)})
}

// OIDCCallback
//
//	@Summary		Exchanges an OIDC token and log in
//	@Description	Exchanges an OIDC token and log in
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			name	path		int		true	"Name of the provider"
//	@Param			state	query		string	true	"State OIDC parameter"
//	@Param			code	query		string	true	"Code OIDC parameter"
//	@Success		200		{object}	OIDCCallbackOutput
//	@Failure		400		{object}	Error
//	@Failure		500		{object}	Error
//	@Router			/auth/callback/{name} [get]
func (a *Api) OIDCCallback(ctx *gin.Context) {
	if !a.Debug {
		state, err := ctx.Cookie("oidc-state")
		if err != nil {
			ctx.JSON(400, gin.H{"error": fmt.Errorf("missing state cookie: %w", err)})
			return
		}

		if ctx.Query("state") != state {
			ctx.JSON(400, gin.H{"error": fmt.Errorf("bad state cookie: %w", err)})
			return
		}
	}

	provider, err := a.ConfigService.GetOIDCProvider(ctx.Param("provider"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	redirect := fmt.Sprintf("%s://%s/auth/callback", scheme, ctx.Request.Host)
	if os.Getenv("OIDC_REDIRECT_BASE_URL") != "" {
		redirect = fmt.Sprintf("%s/auth/callback/%s", os.Getenv("OIDC_REDIRECT_BASE_URL"), provider.Name)
	}

	prvd, err := a.ConfigService.GetOIDCProvider(ctx.Param("provider"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	prv, err := oidc.NewProvider(ctx, prvd.Issuer)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
	}

	fmt.Println(provider.Scopes)

	oauthConfig := oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  redirect,
		Endpoint:     prv.Endpoint(),
		Scopes:       provider.Scopes,
	}

	oauth2Token, err := oauthConfig.Exchange(ctx, ctx.Query("code"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("failed to exchange token: %s", err)})
		return
	}

	verifier := prv.Verifier(&oidc.Config{
		ClientID: provider.ClientID,
	})

	idToken, err := verifier.Verify(ctx, oauth2Token.AccessToken)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err})
		return

	}

	type customClaims struct {
		Name   string   `json:"name"`
		Email  string   `json:"email"`
		Groups []string `json:"groups"`
	}

	var c customClaims

	err = idToken.Claims(&c)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err})
		return

	}

	if user, err := a.UserService.GetUserByUsername(c.Email); errors.Is(err, userservice.ErrUserNotFound) {
		user, err := a.UserService.CreateUser(c.Email, c.Email, "", "oidc", false, c.Name)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"error": fmt.Errorf("failed to create new user: %w", err)})
			return
		}

		token, err := a.UserService.GenerateSessionToken(user)
		if err != nil {
			ctx.JSON(500, gin.H{"error": fmt.Errorf("failed to generate session token: %w", err)})
			return
		}

		ctx.JSON(200, &OIDCCallbackOutput{
			Token:    token,
			Username: c.Email,
		})
		return
	} else if user != nil {
		if user.Kind != "oidc" {
			ctx.JSON(400, gin.H{"error": "user already exists and isn't of kind oidc"})
			return
		}

		err = a.UserService.UpdateUser(user.Id, c.Email, false, c.Name)
		if err != nil {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("failed to update user %s", err)})
			return
		}

		token, err := a.UserService.GenerateSessionToken(user)
		if err != nil {
			ctx.JSON(500, gin.H{"error": fmt.Errorf("failed to generate session token: %w", err)})
			return
		}

		ctx.JSON(200, &OIDCCallbackOutput{
			Token:    token,
			Username: c.Email,
		})
		return
	} else if err != nil {
		ctx.JSON(500, gin.H{"error": "something failed, not sure why"})
		return
	}

	ctx.JSON(500, gin.H{"error": "not sure what happened but we failed to log you in"})
	return
}

// Logout
//
//	@Summary		Logs a user out and invalidates the session
//	@Description	Logs a user out and invalidates the session
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	LoginOutput
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Security		jwt
//	@Router			/auth/logout [post]
func (a *Api) Logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("X-AUTH-TOKEN")
	if token == "" {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})
		return
	}

	if err := a.UserService.LogoutFromToken(token); err == nil {
		ctx.JSON(200, &LogoutOutput{
			Ok: true,
		})
		return
	} else {
		ctx.JSON(401, gin.H{"error": fmt.Sprintf("failed to logout: %s", err)})
	}
}

// GetAvailableOIDCProviders
//
//	@Summary		Gets available OIDC providers
//	@Description	Gets available OIDC providers
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]OIDCProvider
//	@Failure		500	{object}	Error
//	@Router			/auth/oidc/providers [get]
//	@Security		jwt
//	@Security		apikey
func (a *Api) GetAvailableOIDCProviders(ctx *gin.Context) {
	providers, err := a.ConfigService.GetOIDCProviders()
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	provs := make([]OIDCProvider, 0)
	for _, p := range providers {
		provs = append(provs, OIDCProvider{
			Name:        p.Name,
			DisplayName: p.DisplayName,
		})
	}

	ctx.JSON(200, provs)
}

// CreateOIDCProvider
//
//	@Summary		Creates an OIDC provider
//	@Description	Creates an OIDC provider
//	@Tags			Config
//	@Accept			json
//	@Produce		json
//	@Param			request	body		NewOIDCProvider	true	"New OIDC provider config"
//	@Success		200		{object}	OIDCProvider
//	@Failure		500		{object}	Error
//	@Router			/config/oidc/provider [post]
func (a *Api) CreateOIDCProvider(ctx *gin.Context) {
	var input NewOIDCProvider

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": fmt.Errorf("could not bind the request to object: %w", err).Error()})
		return
	}

	if match, _ := regexp.MatchString("^[a-zA-z0-9]+$", input.Name); !match {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "invalid name"})
	}

	prov, err := a.ConfigService.CreateOIDCProvider(
		&configservice.OIDCProvider{
			Name:         input.Name,
			DisplayName:  input.DisplayName,
			ClientID:     input.ClientID,
			ClientSecret: input.ClientSecret,
			Issuer:       input.Issuer,
			Scopes:       input.Scopes,
		},
	)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
	}

	ctx.JSON(200, prov)
}
