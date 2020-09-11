package main

type UserInfo struct {
	Id        string      `json:"sub"`
	Issuer    string      `json:"iss"`
	Email     string      `json:"email"`
	Audiences []string    `json:"aud"`
	Claims    interface{} `json:"claims"`
}
