package search_service

import (
	"context"
	"github.com/ColanderProject/full-text-search-on-disk/util"
	"github.com/go-kit/kit/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	//“github.com/sxwxs/FullTextSearchOnDisk/search”
)

type service struct {
	logger log.Logger
}

type Service interface {
	Search(ctx context.Context, keyword string, startIndex int64) (string, int64, int64, string, string, int64, *timestamppb.Timestamp, *timestamppb.Timestamp)
	GetHeader(ctx context.Context, documentId string) (string, string, *timestamppb.Timestamp, *timestamppb.Timestamp, int64)
	GetContent(ctx context.Context, documentId string) (string, int64, int64, string, int64)
	UpdateDocuments(ctx context.Context, documentId string, document util.Document) (string, string)
}

// NewService func initializes a service
func NewService(logger log.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (s service) Search(ctx context.Context, keyword string, startIndex int64) (string, int64, int64, string, string, int64, *timestamppb.Timestamp, *timestamppb.Timestamp) {
	var currentIndex int64 = 1
	var err string = "This is a fake Search error message"
	var totalIndex int64 = 1
	var title string = "byrSearch"
	var author string = "author"
	var replyCount int64 = 1
	ct, createTimeErr := time.Parse(time.RFC3339, "2020-09-01T21:46:43Z")
	if createTimeErr != nil {
		panic(createTimeErr)
	}
	createTimeStamp := timestamppb.New(ct)

	ut, updateTimeErr := time.Parse(time.RFC3339, "2022-02-23T21:46:43Z")
	if createTimeErr != nil {
		panic(updateTimeErr)
	}
	updateTimeStamp := timestamppb.New(ut)

	return err, currentIndex, totalIndex, title, author, replyCount, createTimeStamp, updateTimeStamp
	/*
		return searchByte(keyword, startIndex)
	*/
}

func (s service) GetHeader(ctx context.Context, documentId string) (string, string, *timestamppb.Timestamp, *timestamppb.Timestamp, int64) {
	var title string = "Introduction to byrSearch"
	var author string = "Kouran Darkhand"
	ct, createTimeErr := time.Parse(time.RFC3339, "2020-09-01T21:46:43Z")
	if createTimeErr != nil {
		panic(createTimeErr)
	}
	createTimeStamp := timestamppb.New(ct)

	ut, updateTimeErr := time.Parse(time.RFC3339, "2022-02-23T21:46:43Z")
	if createTimeErr != nil {
		panic(updateTimeErr)
	}
	updateTimeStamp := timestamppb.New(ut)
	var replyCount int64 = 10

	return title, author, createTimeStamp, updateTimeStamp, replyCount
	/*
		head, err := DB_H.Get(id, nil) //此即为在grpc中定义的GetHeader API
	*/
}

func (s service) GetContent(ctx context.Context, documentId string) (string, int64, int64, string, int64) {
	var author string = "author"
	var upCount int64 = 10
	var downCount int64 = 1
	var text string = "This is a fake GetContent response string"
	var replyId int64 = 1

	return author, upCount, downCount, text, replyId
	/*
		content, err := DB_C.Get(id, nil) //此即为在grpc中定义的GetContent API
	*/
}

func (s service) UpdateDocuments(ctx context.Context, documentId string, document util.Document) (string, string) {
	var response string = "This is a fake UpdateCrawler response string"
	var err string = "This is a fake UpdateCrawler response string"
	return err, response
}
