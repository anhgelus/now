package main

type Data struct {
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Person      *Person `json:"person"`
}

type Person struct {
	Name     string `json:"name"`
	Pronouns string `json:"pronouns"`
	Image    string `json:"image"`
}
