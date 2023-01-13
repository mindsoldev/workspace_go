package main

import "os"

func main() {

	os.Setenv("APP_DB_USERNAME", "gouser")
	os.Setenv("APP_DB_PASSWORD", "gopassword")
	os.Setenv("APP_DB_PASSADDRESS", "127.0.0.1:3306")
	os.Setenv("APP_DB_NAME", "gotest")
	os.Setenv("APP_DB_ADDRESS", "127.0.0.1:3306")

	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_ADDRESS"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8010")
}
