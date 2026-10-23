// Package docs отдаёт Swagger UI. Спека openapi/api.swagger.json генерируется buf generate из proto.
package docs

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
)

//go:embed openapi/api.swagger.json
var spec []byte

// Swagger UI с того же origin, что и API, — браузер сам шлёт cookie access_token/refresh_token.
const page = `<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <title>Chess API</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({ url: '/docs/api.swagger.json', dom_id: '#swagger-ui', persistAuthorization: true });
  </script>
</body>
</html>`

func Register(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(page)
	})
	app.Get("/docs/api.swagger.json", func(c *fiber.Ctx) error {
		c.Type("json")
		return c.Send(spec)
	})
}
