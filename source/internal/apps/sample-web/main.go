package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"source/internal/errors"
	"source/pkg/utility"
	"strings"
	"time"
)

var skipURLs = []string{
	"/login",
	"/register",
}

const (
	JwtSecret      = "my_secret"
	JwtContextKey  = "jwtToken"
	JwtCookieName  = "myCookieLogin"
	JwtTokenLookup = "cookie:" + JwtCookieName
)

func main() {
	app := fiber.New()

	// Casbin
	authorization, _ := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")

	// JWT Middleware
	app.Use(
		jwtware.New(jwtware.Config{
			SigningKey:  []byte(JwtSecret),
			TokenLookup: JwtTokenLookup,
			ContextKey:  JwtContextKey,
			Filter: func(ctx *fiber.Ctx) bool {
				//=> danh sách những page không check login
				return utility.InArray(ctx.Path(), skipURLs, true)
			},
		}),
		func(ctx *fiber.Ctx) error {
			//=> hàm này để parse token JWT rồi storage vào Locals
			var is bool
			jwtToken, is := ctx.Locals(JwtContextKey).(*jwt.Token)
			if !is {
				return ctx.Next()
			}
			claims := jwtToken.Claims.(jwt.MapClaims)
			id, is := claims["id"]
			roles, is := claims["roles"]
			if !is {
				return errors.New(`Tokens are useless`)
			}
			ctx.Locals("userLoginID", id)
			ctx.Locals("roles", roles)
			return ctx.Next()
		},
		func(ctx *fiber.Ctx) error {
			if utility.InArray(ctx.Path(), skipURLs, true) {
				return ctx.Next()
			}
			//=> Sau khi parse JWT token -> lấy ra roles ở trên thì sẽ sử dụng roles đó ở đây để kiểm tra permission
			roles := ctx.Locals("roles").([]interface{})
			allowed := false
			for _, role := range roles {
				result, err := authorization.Enforce(strings.ToLower(fmt.Sprint(role)), ctx.Path(), ctx.Method())
				if result && err == nil {
					allowed = true
				}
			}
			if !allowed {
				return errors.New(`permission denied`)
			}
			return ctx.Next()
		},
	)

	// Login route
	app.Get("/login", login)
	app.Get("/inventory", func(ctx *fiber.Ctx) error {
		return ctx.SendString("inventory")
	})
	app.Get("/payment", func(ctx *fiber.Ctx) error {
		return ctx.SendString("payment")
	})
	// Unauthenticated route
	app.Get("/", accessible)
	app.Get("/accessible", accessible)
	// Restricted Routes
	app.Get("/restricted", restricted)
	app.Listen(":3000")
}

func login(c *fiber.Ctx) error {
	user := c.Query("user")
	pass := c.Query("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var id int64
	id = 10
	// Create the Claims
	claims := jwt.MapClaims{
		"id":      id,
		"sale_id": 8,
		"name":    "John",
		//"roles": []string{"admin"},
		"roles": []string{"member", "sale"},
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	cookie := fiber.Cookie{
		Name:    JwtCookieName,
		Value:   t,
		Expires: time.Now().Add(24 * time.Hour),
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"token": t})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals(JwtContextKey).(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Printf("\n claims: %+v \n", claims)
	name := claims["name"].(string)
	id := int64(claims["id"].(float64))

	typ := reflect.TypeOf(id).Kind()
	fmt.Printf("\n typ: %+v \n", typ)

	return c.SendString("Welcome " + name)
}
