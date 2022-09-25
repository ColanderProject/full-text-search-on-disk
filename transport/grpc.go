package transport

import (
	"context"
	"github.com/ColanderProject/full-text-search-on-disk/endpoint"
	"github.com/ColanderProject/full-text-search-on-disk/util"
	"github.com/ColanderProject/search-protobuf/pb"
	"github.com/go-kit/kit/log"
	gt "github.com/go-kit/kit/transport/grpc"
)

/*
* This struct type implements pb.SearchServiceServer interface
 */
type gRPCServer struct {
	searchHandler          gt.Handler
	getHeaderHandler       gt.Handler
	getContentHandler      gt.Handler
	updateDocumentsHandler gt.Handler
	pb.UnimplementedSearchServiceServer
}

func NewGRPCServer(endpoints endpoint.Endpoints, logger log.Logger) pb.SearchServiceServer {
	return &gRPCServer{
		searchHandler: gt.NewServer(
			endpoints.SearchOperation,
			decodeSearchRequest,
			encodeSearchResponse,
		),

		getHeaderHandler: gt.NewServer(
			endpoints.GetHeaderOperation,
			decodeDocumentRequest,
			encodeHeader,
		),

		getContentHandler: gt.NewServer(
			endpoints.GetContentOperation,
			decodeDocumentRequest,
			encodeContent,
		),

		updateDocumentsHandler: gt.NewServer(
			endpoints.UpdateDocumentsOperation,
			decodeUpdateRequest,
			encodeUpdateResponse,
		),
	}
}

func decodeSearchRequest(_ context.Context, request interface{}) (interface{}, error) {
	//req := request.(*pb.MathRequest) return a pointer type
	var req *pb.SearchRequest = request.(*pb.SearchRequest)
	return util.SearchRequest{Keyword: req.Keyword, StartIndex: req.StartIndex}, nil
}

func encodeSearchResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(util.SearchResponse)
	var headers []util.Header = resp.Headers
	var resHeaders []*pb.Header

	for i := 0; i < len(headers); i++ {
		var header util.Header = headers[i]
		var resHeader *pb.Header = &pb.Header{Title: header.Title, Author: header.Author, CreateTime: header.CreateTime, UpdateTime: header.UpdateTime, ReplyCount: header.ReplyCount}
		resHeaders = append(resHeaders, resHeader)
	}

	return &pb.SearchResponse{Err: resp.Err, CurrentIndex: resp.CurrentIndex, TotalIndex: resp.TotalIndex, Headers: resHeaders}, nil
}

func decodeDocumentRequest(_ context.Context, request interface{}) (interface{}, error) {
	var req *pb.DocumentRequest = request.(*pb.DocumentRequest)
	return util.DocumentRequest{DocumentId: req.DocumentId}, nil
}

func encodeHeader(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(util.Header)
	return &pb.Header{Title: resp.Title, Author: resp.Author, CreateTime: resp.CreateTime, UpdateTime: resp.UpdateTime, ReplyCount: resp.ReplyCount}, nil
}

func encodeContent(_ context.Context, response interface{}) (interface{}, error) {
	var resp util.Content = response.(util.Content)
	var replies []util.Reply = resp.Reply
	var resReplies []*pb.Reply

	for i := 0; i < len(replies); i++ {
		var reply util.Reply = replies[i]
		var resReply *pb.Reply = &pb.Reply{Author: reply.Author, UpCount: reply.UpCount, DownCount: reply.DownCount, Text: reply.Text, ReplyId: reply.ReplyId}
		resReplies = append(resReplies, resReply)
	}

	return &pb.Content{Reply: resReplies}, nil
}

func decodeUpdateRequest(_ context.Context, request interface{}) (interface{}, error) {
	var req *pb.UpdateRequest = request.(*pb.UpdateRequest)
	var document *pb.Document = req.Document
	var header *pb.Header = document.Header
	var content *pb.Content = document.Content
	replies := content.Reply
	var resReplies []util.Reply

	for i := 0; i < len(replies); i++ {
		var reply *pb.Reply = replies[i]
		var resReply util.Reply = util.Reply{Author: reply.Author, UpCount: reply.UpCount, DownCount: reply.DownCount, Text: reply.Text, ReplyId: reply.ReplyId}
		resReplies = append(resReplies, resReply)
	}

	var resDocument util.Document = util.Document{
		Header:  util.Header{Title: header.Title, Author: header.Author, CreateTime: header.CreateTime, UpdateTime: header.UpdateTime, ReplyCount: header.ReplyCount},
		Content: util.Content{Reply: resReplies},
	}

	return util.UpdateRequest{DocumentId: req.DocumentId, Document: resDocument}, nil
}

func encodeUpdateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(util.UpdateResponse)
	return &pb.UpdateResponse{Err: resp.Err, Result: resp.Result}, nil
}

//Search
/*
* These are grpc server API
* All implementations must embed UnimplementedSearchServiceServer for forward compatibility
 */
func (s *gRPCServer) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	_, resp, err := s.searchHandler.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err

	}

	return resp.(*pb.SearchResponse), nil
}

func (s *gRPCServer) GetHeader(ctx context.Context, request *pb.DocumentRequest) (*pb.Header, error) {
	_, resp, err := s.getHeaderHandler.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return resp.(*pb.Header), nil
}

func (s *gRPCServer) GetContent(ctx context.Context, request *pb.DocumentRequest) (*pb.Content, error) {
	_, resp, err := s.getContentHandler.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return resp.(*pb.Content), nil
}

func (s *gRPCServer) UpdateDocuments(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	_, resp, err := s.updateDocumentsHandler.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return resp.(*pb.UpdateResponse), nil
}
