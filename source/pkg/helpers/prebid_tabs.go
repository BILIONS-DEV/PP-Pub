package helpers

const (
	TabBanner    = "banner"
	TabInstream  = "video_instream"
	TabOutstream = "video_outstream"
	TabNative    = "native"
)

type TabInfoStruct struct {
	Id              string
	Name            string
	Active          bool
	Available       map[string]AvailableItem
	MultiConditonal map[string]AvailableItem
}

type AvailableItem struct {
	Name string
	List []string
}

// Display, Play Zone:  Banner + Video Outstream (for size["300x250", "336x280", "970x250", "300x600"]),
// Interstitial: Banner + Video Instream + Video Outstream,
// Sticky Banner: Banner,
// Native: Native,
// Video Outstream: Video Outstream
// Video Instream:
// + Overlay: Banner + Video Instream + Video Outstream
// + Vast: Banner + Video Instream + Video Outstream
// + Normal: Video Instream + Video Outstream

// Native: 1x1, 728x90, 300x250, 160x600, 336x280, 300x600, 970x250, 970x90, 320x100, 320x50
// Video: 1x1, 640x480, 400x300
// Banner: all
var TabInfo = []TabInfoStruct{
	{Id: TabBanner, Name: "Banner", Active: true,
		Available: map[string]AvailableItem{
			"adformat": {
				Name: "Ad Format",
				List: []string{"Display", "Instream", "Sticky Banner", "Interstitial", "Play Zone"},
			},
			"adsize": {
				Name: "Ad Size",
				List: []string{"all"},
			},
		},
	},

	{Id: TabInstream, Name: "Video Instream",
		Available: map[string]AvailableItem{
			"adformat": {
				Name: "Ad Format",
				List: []string{"Instream", "Interstitial"},
			},
			"adsize": {
				Name: "Ad Size",
				List: []string{"1x1", "640x480", "400x300"},
			},
		},
	},

	{Id: TabOutstream, Name: "Video Outstream",
		Available: map[string]AvailableItem{
			"adformat": {
				Name: "Ad Format",
				List: []string{"Instream", "Outstream", "Interstitial"},
			},
			"adsize": {
				Name: "Ad Size",
				List: []string{"1x1", "640x480", "400x300"},
			},
		},
		MultiConditonal: map[string]AvailableItem{
			"adformat": {
				Name: "Ad Format",
				List: []string{"Display", "Play Zone"},
			},
			"adsize": {
				Name: "Ad Size",
				List: []string{"300x250", "336x280", "970x250", "300x600"},
			},
		},
	},

	{Id: TabNative, Name: "Native",
		Available: map[string]AvailableItem{
			"adformat": {
				Name: "Ad Format",
				List: []string{"Native"},
			},
			"adsize": {
				Name: "Ad Size",
				List: []string{"1x1", "728x90", "300x250", "160x600", "336x280", "300x600", "970x250", "970x90", "320x100", "320x50"},
			},
		},
	},
}
