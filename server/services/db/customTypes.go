package db

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	ID     string   `json:"id"`
	Sender string   `json:"sender"`
	Text   string   `json:"text"`
	Media  []string `json:"media"`
	Date   int64    `json:"date"` // func (t time.Time) UnixMilli() int64
}

type Chat struct {
	ID       string   `json:"id"`
	Users    []string `json:"users"`
	Messages []string `json:"messages"`
	Active   int64    `json:"active"`
}
type Event struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	AssignedTo         []string `json:"assigned_to"`
	StartDate          int64    `json:"start_date"`
	EndDate            int64    `json:"end_date"`
	LocationName       string   `json:"location_name"`
	LocationAddress    string   `json:"location_address"`
	LocationLong       int64    `json:"location_long"`
	LocationLat        int64    `json:"location_lat"`
	Notes              string   `json:"notes"`
	FirstNotification  int64    `json:"first_notification"`
	SecondNotification int64    `json:"second_notification"`
	Active             bool     `json:"active"`
	CreatedAt          int64    `json:"created_at"`
	UpdatedAt          int64    `json:"updated_at"`
}

type Item struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Price       int64    `json:"price"`
	Inventory   int      `json:"inventory"`
	Active      bool     `json:"active"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

type Order struct {
	ID        string   `json:"id"`
	User      string   `json:"user"`
	Items     []string `json:"items"`
	Total     int64    `json:"total"`
	Status    string   `json:"status"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
}
