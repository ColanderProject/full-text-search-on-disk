package main

import (
	"os"
	search_service "github.com/sxwxs/FullTextSearchOnDisk/search"
)

func main() {
	search_service.MainServer(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
}
