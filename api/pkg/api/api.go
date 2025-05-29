package api

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/thomas-maurice/api/go-vue/docs"
	"github.com/thomas-maurice/api/go-vue/pkg/config"
	"github.com/thomas-maurice/api/go-vue/pkg/embeded"
	"github.com/thomas-maurice/api/go-vue/pkg/services/configservice"
	sqlconfigservice "github.com/thomas-maurice/api/go-vue/pkg/services/configservice/sql"
	"github.com/thomas-maurice/api/go-vue/pkg/services/userservice"
	sqluserservice "github.com/thomas-maurice/api/go-vue/pkg/services/userservice/sql"
	"github.com/thomas-maurice/api/go-vue/pkg/store"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Api struct {
	Router           *gin.Engine
	Config           *config.Config
	SigninigKey      *ecdsa.PrivateKey
	SigningPublicKey *ecdsa.PublicKey
	Debug            bool
	DB               *gorm.DB
	UserService      userservice.UserService
	ConfigService    configservice.ConfigService
	OidcProvider     *oidc.Provider
}

func NewAPI(cfgFile string) (*Api, error) {
	cfg, err := config.LoadFromFile(cfgFile)
	if err != nil {
		return nil, err
	}

	decodedPem, _ := pem.Decode([]byte(cfg.Security.SigninigKey))
	pKey, err := x509.ParseECPrivateKey(decodedPem.Bytes)
	if err != nil {
		return nil, err
	}

	var a Api

	a.Debug = cfg.Debug
	if !a.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := store.NewSqlStore(cfg.Storage.Driver, cfg.Storage.URL)
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	}

	a.DB = db

	us, err := sqluserservice.NewUserService(db, pKey)
	if err != nil {
		return nil, err
	}

	a.UserService = us

	cs, err := sqlconfigservice.NewConfigService(db)
	if err != nil {
		return nil, err
	}

	a.ConfigService = cs

	/*var admin models.User
	if err = db.First(&admin, &models.User{Username: "admin"}).Error; errors.Is(gorm.ErrRecordNotFound, err) {
		fmt.Println("creating admin user")
		if err = db.Create(&models.User{
			Username: "admin",
			Email:    "admin@localhost",
			Password: cfg.Security.AdminPassword,
			Admin:    true,
			Kind:     "local",
			Created:  time.Now(),
		}).Error; err != nil {
			return nil, err
		}
	}*/

	_, err = us.GetUserByUsername("admin")
	if err == userservice.ErrUserNotFound {
		_, err = us.CreateUser("admin", "admin@localhost", cfg.Security.AdminPassword, "local", true, "Admin")
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	router := gin.Default()

	router.Use(
		a.DebugCORSMiddleware,
	)

	if cfg.Security.OIDC != nil {
		for name, config := range cfg.Security.OIDC {
			_, err := cs.UpsertOIDCProvider(&configservice.OIDCProvider{
				Name:         name,
				Active:       true,
				ClientID:     config.ClientID,
				ClientSecret: config.ClientSecret,
				Issuer:       config.Issuer,
				DisplayName:  config.DisplayName,
				Scopes:       config.Scopes,
			})

			if err != nil {
				return nil, err
			}
		}
	}

	apiGroup := router.Group("/api")
	authGroup := apiGroup.Group("/auth")
	{
		authGroup.GET("/oidc/:provider", a.GenerateOIDCRedirectURL)
		authGroup.GET("/callback/:provider", a.OIDCCallback)
		authGroup.GET("/oidc/providers", a.GetAvailableOIDCProviders)
		authGroup.POST("/login", a.AuthPassword)
		authGroup.POST("/logout", a.Logout)
	}

	userGroup := apiGroup.Group("/user", a.RequiresUserLogin(false))
	{
		userGroup.GET("/profile", a.ProfileSelf)
	}

	adminGroup := apiGroup.Group("/admin", a.RequiresUserLogin(true))
	{
		adminGroup.GET("/users", a.AdminListUsers)
		adminGroup.GET("/user/:id", a.AdminGetUser)
	}

	configGroup := apiGroup.Group("/config", a.RequiresUserLogin(true))
	{
		configGroup.POST("/oidc/provider", a.CreateOIDCProvider)
	}

	apiGroup.GET("/ping", a.UserTokenMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"pong": time.Now()})
	})

	apiGroup.GET("/uuid", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"uuid": uuid.NewString()})
	})

	subbed, err := fs.Sub(embeded.UserInterfaceFS, "ui/dist")
	if err != nil {
		panic(err)
	}

	assets, err := fs.Sub(subbed, "assets")
	if err != nil {
		panic(err)
	}

	assetsGroup := router.Group("/assets", func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	})
	assetsGroup.StaticFS("", http.FS(assets))

	router.GET("/", func(ctx *gin.Context) {
		ctx.FileFromFS("/", http.FS(subbed))
	})

	router.GET("/index.html", func(ctx *gin.Context) {
		ctx.FileFromFS("/", http.FS(subbed))
	})

	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.FileFromFS("/favicon.ico", http.FS(subbed))
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// SPA fallback for client-side routing
	router.NoRoute(func(ctx *gin.Context) {
		if len(ctx.Request.URL.Path) >= len("/api") && ctx.Request.URL.Path[:len("/api")] == "/api" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}

		ctx.FileFromFS("/", http.FS(subbed))
	})

	a.Router = router
	a.SigningPublicKey = &pKey.PublicKey
	a.SigninigKey = pKey
	a.Config = cfg

	return &a, nil
}

type PingOutput struct {
	Pong string `json:"pong"`
}

// Ping
//
//	@Summary		Pings the server and check authentication
//	@Description	Pings the server and check authentication
//	@Produce		json
//	@Success		200	{object}	PingOutput
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Security		jwt
//	@Security		apikey
//	@Router			/ping [get]
func (a *Api) Ping(ctx *gin.Context) {
	ctx.JSON(200, &PingOutput{Pong: time.Now().Format(time.RFC1123Z)})
}

func (a *Api) Run() error {
	return a.Router.Run(a.Config.HTTP.Listen)
}
