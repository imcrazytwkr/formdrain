package sendinblue

type EmailContact struct {
	Name    string `bson:"name,omitempty" json:"name,omitempty"`
	Address string `bson:"email" json:"email"`
}
