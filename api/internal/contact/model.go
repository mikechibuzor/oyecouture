package contact

type Submission struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Topic   string `json:"topic"`
	Rating  int    `json:"rating"`
	Message string `json:"message"`
}

type Response struct {
	Received bool `json:"received"`
}
