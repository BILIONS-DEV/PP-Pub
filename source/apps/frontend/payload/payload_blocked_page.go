package payload

type CSV struct {
	Site                    string `csv:"Site or app" json:"site"`
	Entity                  string `csv:"Entity" json:"entity"`
	Platform                string `csv:"Platform" json:"platform"`
	AppPackageNameOrStoreID string `csv:"App package name or store ID" json:"app_package_name_or_store_id"`
	IssueLocation           string `csv:"Issue location" json:"issue_location"`
	MustFix                 string `csv:"Must fix" json:"must_fix"`
	Issues                  string `csv:"Issues" json:"issues"`
	Status                  string `csv:"Status" json:"status"`
	PropertyCodes           string `csv:"Property codes" json:"property_codes"`
	AdRequestsLast7Days     string `csv:"Ad requests - last 7 days" json:"ad_requests_last_7_days"`
	DateReported            string `csv:"Date reported" json:"date_reported"`
	LastDateFound           string `csv:"Last date found" json:"last_date_found"`
}

type BlockedPageSubmit struct {
	Id     int64    `json:"id"`
	Name   string   `json:"name"`
	Domain int64    `json:"domain"`
	Pages  []string `json:"pages"`
}
