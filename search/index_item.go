package search_service

type IndexChange struct {
	Add []int
	Remove map[int] Empty
}