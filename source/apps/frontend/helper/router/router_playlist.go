package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Playlist(app *fiber.App) {
	Playlist := new(controller.Playlist)
	app.Get(config.URIPlaylist, Playlist.Index)
	app.Post(config.URIPlaylist, Playlist.Filter)
	app.Get(config.URIPlaylistAdd, Playlist.Add)
	app.Post(config.URIPlaylistAdd, Playlist.AddPost)
	app.Get(config.URIPlaylistEdit, Playlist.Edit)
	app.Get(config.URIPlaylistView, Playlist.View)
	app.Post(config.URIPlaylistEdit, Playlist.EditPost)
	app.Post(config.URIPlaylistDel, Playlist.Delete)
	app.Post(config.URIPlaylistCollapse, Playlist.Collapse)
}
