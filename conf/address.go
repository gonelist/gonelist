package conf

import (
	"golang.org/x/oauth2"
)

type Remote struct {
	Name     string `json:"name"`
	EndPoint oauth2.Endpoint
	ROOTUrl  string `json:"root_url"`
	UrlBegin string `json:"url_begin"`
	UrlEnd   string `json:"url_end"`
}

var (
	// Onedrive 国际版
	OneDrive = Remote{
		Name: "onedrive",
		EndPoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
		ROOTUrl:  "https://graph.microsoft.com/v1.0/me/drive",
		UrlBegin: "https://graph.microsoft.com/v1.0/me/drive/root:",
		UrlEnd:   ":/children",
	}
	// 世纪互联
	ChinaCloud = Remote{
		Name: "chinacloud",
		EndPoint: oauth2.Endpoint{
			AuthURL:  "https://login.partner.microsoftonline.cn/organizations/oauth2/v2.0/authorize",
			TokenURL: "https://login.partner.microsoftonline.cn/organizations/oauth2/v2.0/token",
		},
		ROOTUrl:  "https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive",
		UrlBegin: "https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive/root:",
		UrlEnd:   ":/children",
	}
	// TODO SharePoint, ...
)
