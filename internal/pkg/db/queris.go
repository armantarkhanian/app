package db

var (
	SelectUser = "SELECT userID, username, email FROM users WHERE userID=?"
	InsertUser = "INSERT INTO users (userID, username, email) VALUES (?, ?, ?)"
)
