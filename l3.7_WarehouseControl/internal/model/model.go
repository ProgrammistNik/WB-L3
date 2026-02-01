package model

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

type Item struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UpdatedAt string `json:"updated_at"`
}

type ItemHistory struct {
	ID        int         `json:"id"`
	ItemID    int         `json:"item_id"`
	Action    string      `json:"action"`
	ChangedBy string      `json:"changed_by"`
	OldData   interface{} `json:"old_data"`
	NewData   interface{} `json:"new_data"`
	ChangedAt string      `json:"changed_at"`
}
