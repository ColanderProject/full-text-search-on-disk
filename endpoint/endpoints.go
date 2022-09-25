package endpoint

import (
	"context"
	"github.com/ColanderProject/full-text-search-on-disk/indexService"
	"github.com/ColanderProject/full-text-search-on-disk/util"
	"github.com/go-kit/kit/endpoint"
)

/*
Endpoints type struct defined in this file.
An endpoint represents a single RPC method.
Each Service method converts into an Endpoint to make RPC style communication between servers and clients.
*/
type Endpoints struct {
	SearchOperation          endpoint.Endpoint
	GetHeaderOperation       endpoint.Endpoint
	GetContentOperation      endpoint.Endpoint
	UpdateDocumentsOperation endpoint.Endpoint
}

/*
MakeEndpoints func initializes the Endpoint instances.
Endpoint is the fundamental building block of servers and clients, it represents a single RPC method.
The input is Service interface defined in Service tier.
*/
func MakeEndpoints(s indexService.Service) Endpoints {
	return Endpoints{
		SearchOperation:          makeSearchKeywordEndpoint(s),
		GetHeaderOperation:       makeGetHeaderEndpoint(s),
		GetContentOperation:      makeGetContentEndpoint(s),
		UpdateDocumentsOperation: makeUpdateDocumentsEndpoint(s),
	}
}

func makeSearchKeywordEndpoint(s indexService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(util.SearchRequest)
		errMsg, currentIndex, totalIndex, title, author, replyCount, createTimeStamp, updateTimeStamp := s.Search(ctx, req.Keyword, req.StartIndex)
		headers := []util.Header{
			{title, author, createTimeStamp, updateTimeStamp, replyCount},
		}
		return util.SearchResponse{Err: errMsg, CurrentIndex: currentIndex, TotalIndex: totalIndex, Headers: headers}, nil
	}
}

func makeGetHeaderEndpoint(s indexService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(util.DocumentRequest)
		title, author, createTimeStamp, updateTimeStamp, replyCount := s.GetHeader(ctx, req.DocumentId)
		return util.Header{Title: title, Author: author, CreateTime: createTimeStamp, UpdateTime: updateTimeStamp, ReplyCount: replyCount}, nil
	}
}

func makeGetContentEndpoint(s indexService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(util.DocumentRequest)
		author, upCount, downCount, text, replyId := s.GetContent(ctx, req.DocumentId)
		return util.Content{
			Reply: []util.Reply{
				{author, upCount, downCount, text, replyId},
			},
		}, nil
	}
}

func makeUpdateDocumentsEndpoint(s indexService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(util.UpdateRequest)
		errMsg, statusRes := s.UpdateDocuments(ctx, req.DocumentId, req.Document)
		return util.UpdateResponse{Err: errMsg, Result: statusRes}, nil
	}
}

/*
之所以要使用return func这样的定义方式，是因为在go-kit中定义的Endpoint就是一个这样的形式，它定义的func没有具体的实现，
只有参数和返回值的声明，有点像是interface中的一个method，在实际函数中我们实现了这个func，给予了它具体的函数内容。
*/
