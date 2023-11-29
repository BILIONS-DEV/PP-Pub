package router

import (
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	Home(app)
	User(app)
	Inventory(app)
	LineItem(app)
	Target(app)
	AdTag(app)
	Player(app)
	Payment(app)
	Playlist(app)
	Report(app)
	Floor(app)
	Gam(app)
	Content(app)
	Bidder(app)
	//Config(app)
	Video(app)
	AdsTxt(app)
	Blocking(app)
	Support(app)
	Notification(app)
	Identity(app)
	Channels(app)
	AbTesting(app)
	History(app)
	Rule(app)
	BlockedPage(app)

	LineItemV2(app)
	InventoryV2(app)
	AdTagV2(app)
	//PlayerV2(app)
	AdBlock(app)
}
