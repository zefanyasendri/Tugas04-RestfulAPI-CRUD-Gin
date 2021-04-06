package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Tugas04-RestfulAPI-CRUD-Gin/controllers"
	"github.com/Tugas04-RestfulAPI-CRUD-Gin/models"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var identityKey = "id"

func main() {
	r := gin.Default()

	// JWT Middleware Gin
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Zefa zone",
		Key:         []byte("160301"),
		Timeout:     time.Hour, // Token Berlaku hanya 1 Jam
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Student); ok {
				return jwt.MapClaims{
					identityKey: v.NIM,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.Student{
				NIM: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals models.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &models.Student{
					NIM:  userID,
					Name: "Wu",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.Student); ok && v.NIM == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	// Pengecekan kalau tidak ada page yang ditemukan
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// Model
	db := models.ConnectDB()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Untuk Login
	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/auth")
	// Untuk Refresh Token
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	// Untuk Logout
	auth.GET("/logout", authMiddleware.LogoutHandler)

	// CRUD Auth
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "Success Login with JWT"})
		})

		// CRUD
		auth.GET("/students", controllers.ReadDataStudent)
		auth.GET("/student/:name", controllers.ReadDataOneStudent)
		auth.GET("/login/:name", controllers.Login)
		auth.POST("/student", controllers.CreateDataStudent)
		auth.PUT("/student/:nim", controllers.UpdateDataStudent)
		auth.DELETE("/student/:nim", controllers.DeleteDataStudent)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Restful API - Gin"})
	})

	r.Run()
}
