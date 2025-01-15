package main

import "github.com/tnfy-link/backend/internal"

//go:generate swag init --parseDependency -g ./main.go -o ./api

//	@title			tnfy.link backend API
//	@version		{{VERSION}}
//	@description	The backend of the tnfy.link URL shortener allows you to shorten URLs and get statistics.

//	@contact.name	tnfy.link Support
//	@contact.email	support@tnfy.link

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		api.tnfy.link
//	@BasePath	/v1
//  @schemes	https

func main() {
	internal.Run()
}
