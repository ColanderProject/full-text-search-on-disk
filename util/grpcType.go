package util

import "google.golang.org/protobuf/types/known/timestamppb"

type SearchRequest struct {
	Keyword    string
	StartIndex int64
}

type SearchResponse struct {
	Err          string
	CurrentIndex int64
	TotalIndex   int64
	Headers      []Header
}

type DocumentRequest struct {
	DocumentId string
}

type Header struct {
	Title      string
	Author     string
	CreateTime *timestamppb.Timestamp
	UpdateTime *timestamppb.Timestamp
	ReplyCount int64
}

type Content struct {
	Reply []Reply
}

type Reply struct {
	Author    string
	UpCount   int64
	DownCount int64
	Text      string
	ReplyId   int64
}
type Document struct {
	Header  Header
	Content Content
}

type UpdateRequest struct {
	DocumentId string
	Document   Document
}

type UpdateResponse struct {
	Err    string
	Result string
}
