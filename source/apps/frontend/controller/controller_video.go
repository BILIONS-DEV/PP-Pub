package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config/assign"
)

type Video struct{}

type AssignVideoIndex struct {
	assign.Schema
	LinkVideo string
}

func (t *Video) Index(ctx *fiber.Ctx) error {
	idVideo := ctx.Query("id")
	var LinkVideo string
	assigns := AssignVideoIndex{Schema: assign.Get(ctx)}
	switch idVideo {
	case "0":
		LinkVideo = "https://ci5.googleusercontent.com/proxy/FVmcn6ZUFl9VtGKLiwtTvJ0owuQxxy3yzbPMRRiTdxDyVNywVF7EJuC1VmiMOTBhqmUWKsdEm6cRxx0P6nvg9kLwV-RC7E1QDOVb=s0"
		break
	case "1":
		LinkVideo = "https://ci5.googleusercontent.com/proxy/8-mD39baUJIP-0enOlj-vWT7yQ6DKFHoRRC2PBUJ_KuuRqRPr-oeVPyBKwyydBiQVhS5qt2SQ5BHxJJo2fRLEz9pGpLLgkaBSynl=s0"
		break
	case "2":
		LinkVideo = "https://ci3.googleusercontent.com/proxy/vdVfOmZiR_uUorY2c-MDBtAlFCjcdae_s4zZ1vR8ZKYBK7jIGNAMX9P6jj9ELRmWNcAsGhhpVY6biZLKEIYfR6uYGY2SYGB3_ie0=s0"
		break
	case "4":
		LinkVideo = assigns.RootDomain + "/static/js/video/video.mp4"
		break
	}
	assigns.LinkVideo = LinkVideo
	//assigns.Title = config.TitleWithPrefix("List bidder")
	return ctx.Redirect(LinkVideo)
}
