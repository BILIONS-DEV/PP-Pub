package mysql2

type DataObject struct {
	DomainName            string             `json:"domainName"`
	RootDomain            string             `json:"rootDomain"`
	TagDomain             string             `json:"tagDomain"`
	DomainID              int64              `json:"domainID"`
	UId                   string             `json:"uId"`
	UAMPubId              string             `json:"UAMPubId"`
	CaPub                 string             `json:"caPub"`
	GaId                  string             `json:"gaId"`
	PrebidJs              string             `json:"prebidJs"`
	PrebidTimeout         int                `json:"prebidTimeout"`
	AdRefreshTime         int                `json:"adRefreshTime"`
	GranularityMultiplier int                `json:"granularityMultiplier"`
	AdLoadType            string             `json:"adloadType"`
	SafeFrame             string             `json:"safeFrame"`
	Currency              string             `json:"currency"`
	ReloadMode            string             `json:"reloadMode"`
	BidAdjustment         map[string]float64 `json:"bidAdjustment"`
	BlockAdDomains        map[string]string  `json:"blockAdDomains"`
	AliasBidders          map[string]string  `json:"aliasBidders"`
	S2sBidders            []string           `json:"s2sBidders"`
	CreativeIds           []string           `json:"creativeIds"`
	CMP                   CMP                `json:"CMP"`
	UserIds               UserIds            `json:"userIds"`
	AdsTags               []AdsTags          `json:"adsTags"`
}

type AdsTags struct{}

type CMP struct{}

type UserIds struct{}
