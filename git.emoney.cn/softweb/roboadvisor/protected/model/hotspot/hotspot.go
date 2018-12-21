package model

import "html/template"

type Hotspot struct {
	Title string
	GroupName string
	GroupCode string
	TopicSummary string
	TopicContent string
	TopicPic string
	CreateTime string
	TopicSummaryHTML template.HTML
}

