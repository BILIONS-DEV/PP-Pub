package aerospike

type SetAutoAd struct {
	AutoAds []AutoAd `as:"autoAds"`
}

type AutoAd struct {
	UrlType   string           `json:"url_type" as:"urlType"`
	Url       string           `json:"url" as:"url"`
	Root      string           `json:"root" as:"root"`
	Selectors []AutoAdSelector `json:"selectors" as:"selectors"`
}

type AutoAdSelector struct {
	DeviceType            string                     `json:"device_type" as:"deviceType"`
	DeviceConfig          AutoAdSelectorDeviceConfig `json:"device_config" as:"deviceConfig"`
	Selector              string                     `json:"selector" as:"selector"`
	AdPosition            string                     `json:"ad_position" as:"adPosition"`
	AdContent             AutoAdSelectorAdContent    `json:"ad_content" as:"adContent"`
	Scan                  string                     `json:"scan" as:"scan"`
	AdSpace               int                        `json:"ad_space" as:"adSpace"`
	LimitAd               int                        `json:"limit_ad" as:"limitAd"`
	RenderFirstImpression bool                       `json:"renderFirstImpression" as:"rfImpression"`
	Advertisement         bool                       `json:"advertisement" as:"advertisement"`
	AdStyle               AutoAdSelectorStyle        `json:"ad_style" as:"adStyle"`
}

type AutoAdSelectorDeviceConfig struct {
	Min    []int    `json:"min" as:"min"`
	Device []string `json:"device" as:"device"`
}

type AutoAdSelectorAdContent struct {
	Type    string `json:"type" as:"type"`
	Tag     string `json:"tag" as:"tag"`
	Content string `json:"content" as:"content"`
}

type AutoAdSelectorStyle struct {
	PaddingLeft   int    `json:"padding_left" as:"paddingLeft"`
	PaddingTop    int    `json:"padding_top" as:"paddingTop"`
	PaddingRight  int    `json:"padding_right" as:"paddingRight"`
	PaddingBottom int    `json:"padding_bottom" as:"paddingBottom"`
	Align         string `json:"align" as:"align"`
	Sticky4k      string `json:"sticky_4k" as:"sticky4k"`
}
