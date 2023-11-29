package view

import "source/pkg/utility"

const (
	LAYOUTMain       = "_layouts/main"
	LAYOUTLogin      = "_layouts/login"
	LAYOUTTemplate   = "_layouts/template"
	LAYOUTTemplateV2 = "_layouts/template-v2"
	LAYOUTEmpty      = "_layouts/empty"
	LAYOUTTest       = "_layouts/test"
)

func IsActiveSidebar(uri string, path ...string) (flag bool) {
	if utility.InArray(uri, path, true) {
		flag = true
	}
	return
}

func IsActiveSidebarWithGroup(uri string, path []string) (flag bool) {
	if utility.InArray(uri, path, true) {
		flag = true
	}
	return
}
