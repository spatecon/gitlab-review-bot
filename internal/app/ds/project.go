package ds

type Project struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
	URL  string `bson:"url"`
}
