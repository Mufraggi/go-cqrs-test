package views

type ArticleViewModel struct {
	ID         string
	Title      string
	Content    string
	ClapsCount int
}

type UserViewModel struct {
	ID        string
	FirstName string
	LastName  string
	Articles  []ArticleViewModel
}
