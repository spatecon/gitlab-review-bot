package ds

import "time"

type Project struct {
	ID        int       `bson:"id"`
	Name      string    `bson:"name"`
	URL       string    `bson:"url"`
	CreatedAt time.Time `bson:"created_at"`
}
