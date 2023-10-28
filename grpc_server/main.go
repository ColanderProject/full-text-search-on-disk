package main

import (
	context "context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/cockroachdb/pebble"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var db_head *pebble.DB
var db_content *pebble.DB
var db_index *pebble.DB
var db_url_to_id *pebble.DB
var nextDocId int64

var (
	port    = flag.Int("port", 50051, "The server port")
	db_path = flag.String("db", "db", "Database root path")
)

type server struct {
	SearchServiceServer
	mode ServerMode
}

func (s *server) GetMaxURL(_ context.Context, _ *GetMaxURLRequest) (*DocID, error) {
	var doc_id DocID
	doc_id.Id = nextDocId - 1
	return &doc_id, nil
}

func (s *server) URL2ID(_ context.Context, url *URL) (*DocID, error) {
	raw_url, err := proto.Marshal(url)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	raw_doc_id, _, err := db_url_to_id.Get(raw_url)
	var doc_id DocID
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	proto.Unmarshal(raw_doc_id, &doc_id)
	return &doc_id, nil
}

func (s *server) GetChildDocument(context.Context, *DocumentRequest) (*DocumentList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChildDocument not implemented")
}

func (s *server) GetDocumentHeader(_ context.Context, doc_id *DocID) (*Document, error) {
	raw_doc_id, err := proto.Marshal(doc_id)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	raw_head, _, err := db_head.Get(raw_doc_id)
	if err == nil {
		log.Println(err)
		return nil, nil
	}
	var doc_head Document
	proto.Unmarshal(raw_head, &doc_head)
	return &doc_head, nil
}

func (s *server) GetMode(context.Context, *GetModeRequest) (*ServerState, error) {
	var state ServerState
	state.Mode = s.mode
	return &state, nil
}

func (s *server) SetMode(_ context.Context, state *ServerState) (*SetModeResponse, error) {
	s.mode = state.Mode
	var response SetModeResponse
	response.State = StateCode_SUCCESS
	response.CurrentMode = state.Mode
	return &response, nil
}

func (s *server) DeleteDocument(context.Context, *DocID) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDocument not implemented")
}

func (s *server) InsertDocument(_ context.Context, req *InsertRequest) (*InsertResponse, error) {
	var docId DocID
	var response InsertResponse
	response.Id = -1
	req.Document.Id = atomic.AddInt64(&nextDocId, 1)
	req.Document.Id--
	docId.Id = req.Document.Id
	raw_id, err := proto.Marshal(&docId)
	if err != nil {
		log.Println(err)
		return &response, nil
	}
	log.Printf("Current id = %d, url = %s", nextDocId, req.Document.Url)
	raw_doc, err := proto.Marshal(req.Document)
	response.Id = docId.Id
	if err != nil {
		log.Println(err)
		return &response, nil
	}
	db_head.Set(raw_id, raw_doc, pebble.Sync)
	return &response, nil
}

func (s *server) UpdateDocument(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDocument not implemented")
}

func (s *server) Search(context.Context, *SearchRequest) (*DocumentList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}

func main() {
	flag.Parse()
	var err error
	db_head, err = pebble.Open(*db_path+"/doc_head", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	db_content, err = pebble.Open(*db_path+"/doc_content", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	db_index, err = pebble.Open(*db_path+"/doc_inverted_index", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	db_url_to_id, err = pebble.Open(*db_path+"/url_to_id", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterSearchServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	raw_doc_id, _, err := db_url_to_id.Get([]byte("MAX_DOC_ID_KEY"))

	if err != nil {
		log.Println(err)
		nextDocId = 0
	} else {
		var doc_id DocID
		proto.Unmarshal(raw_doc_id, &doc_id)
		nextDocId = doc_id.Id
	}
}
