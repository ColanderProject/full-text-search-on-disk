package main

import (
	search_service "github.com/ColanderProject/full-text-search-on-disk/search"
	"os"
)

func main() {
	search_service.MainServer(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
}
