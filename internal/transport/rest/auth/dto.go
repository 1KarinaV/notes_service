package auth

type Credentials struct {
	Login    string `validate:"min=6,max=20,regexp=^(\\w|-|_|@|\\.|)*$" json:"login"`
	Password string `validate:"min=6,max=40,regexp=^(\\w|-|_|=|@|\\.|\\?|!|%|#|\\$)*$" json:"password"`
}
