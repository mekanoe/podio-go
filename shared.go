package podio

type User struct {
	UserID     int    `json:"user_id,omitempty"`
	SpaceID    int    `json:"space_id,omitempty"`
	ProfileID  int    `json:"profile_id,omitempty"`
	Name       string `json:"name,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
	Type       string `json:"type,omitempty"`
	LastSeenOn string `json:"last_seen_on,omitempty"`
	Link       string `json:"link,omitempty"`
	OrgID      int    `json:"org_id,omitempty"`
	Image      Image  `json:"image,omitempty"`
}

type Image struct {
	HostedBy              string `json:"hosted_by,omitempty"`
	HostedByHumanizedName string `json:"hosted_by_humanized_name,omitempty"`
	ThumbnailLink         string `json:"thumbnail_link,omitempty"`
	Link                  string `json:"link,omitempty"`
	FileID                int    `json:"file_id,omitempty"`
	ExternalFileID        string `json:"external_file_id,omitempty"`
	LinkTarget            string `json:"link_target,omitempty"`
}
