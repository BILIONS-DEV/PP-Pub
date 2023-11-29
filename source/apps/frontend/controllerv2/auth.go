package controllerv2

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"source/internal/errors"
	"source/pkg/utility"
	"strings"
)

type authConfig struct {
	Secret      string
	TokenLookup string
	ContextKey  string
	CookieName  string
}

var skipURLs = []string{
	"/login",
	"/register",
	"/user/login2",
}

func Auth(app fiber.Router, cfg authConfig) {
	// Casbin
	authorization, err := casbin.NewEnforcer("./controllerv2/auth_model.conf", "./controllerv2/auth_policy.csv")
	if err != nil {
		panic(err)
	}
	// JWT Middleware
	app.Use(
		jwtware.New(jwtware.Config{
			SigningKey:  []byte(cfg.Secret),
			TokenLookup: cfg.TokenLookup,
			ContextKey:  cfg.ContextKey,
			//=> danh sách những page không check login
			Filter: func(ctx *fiber.Ctx) bool {
				return skip(ctx)
			},
		}),
		func(ctx *fiber.Ctx) error {
			if skip(ctx) {
				return ctx.Next()
			}
			//=> hàm này để parse token JWT rồi storage vào Locals
			var is bool
			jwtToken, is := ctx.Locals(cfg.ContextKey).(*jwt.Token)
			if !is {
				return ctx.Next()
			}
			claims, is := jwtToken.Claims.(jwt.MapClaims)
			if !is {
				return ctx.Next()
			}
			userLoginID, is := claims["id"].(float64)
			roles, is := claims["roles"]
			fmt.Println(roles)
			if !is {
				return errors.New(`Tokens are useless`)
			}
			accountManagerID, _ := claims["amid"].(float64)
			ctx.Locals("userLoginID", int64(userLoginID))
			ctx.Locals("accountManagerID", int64(accountManagerID))
			ctx.Locals("roles", roles)
			return ctx.Next()
		},
		func(ctx *fiber.Ctx) error {
			if skip(ctx) {
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
}

// skip : loại bỏ các url ra khỏi các quá trình check login (JWT) + check permission
func skip(ctx *fiber.Ctx) bool {
	return utility.InArray(ctx.Path(), skipURLs, true)
}
