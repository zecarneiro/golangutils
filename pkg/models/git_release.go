package models

type GitRelease struct {
	TagName string `json:"tag_name,omitempty"`
	Version string
}
