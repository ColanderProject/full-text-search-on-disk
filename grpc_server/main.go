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
var db_children_id *pebble.DB
var db_prepare_index *pebble.DB
var db_index *pebble.DB
var db_url_to_id *pebble.DB
var nextDocId int64

func getBytesCharSet(str string, set0 *Set) {
	for _, chi := range str {
		ch := string(chi)
		if ch != " " && ch != "\t" && ch != "\n" && ch != "\r" && ch != "'" && ch != "\"" && ch != "\xa0" {
			set0.Add(ch)
		}
	}
}

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

func (s *server) GetChildDoc(context.Context, *DocumentRequest) (*DocList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChildDocument not implemented")
}

func (s *server) GetChildDocId(context.Context, *DocumentRequest) (*DocIdList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChildDocument not implemented")
}

func (s *server) GetDocumentHeader(_ context.Context, doc_id *DocID) (*DocHeader, error) {
	raw_doc_id, err := proto.Marshal(doc_id)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	raw_head, _, err := db_head.Get(raw_doc_id)
	if err != nil {
		log.Println("To get docid=", doc_id.Id, err)
		return nil, nil
	}
	var doc_head DocHeader
	proto.Unmarshal(raw_head, &doc_head)
	return &doc_head, nil
}

func (s *server) GetDocumentHeaderByURL(_ context.Context, url *URL) (*DocHeader, error) {
	raw_url, err := proto.Marshal(url)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	raw_doc_id, _, err := db_url_to_id.Get(raw_url)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	raw_head, _, err := db_head.Get(raw_doc_id)
	if err == nil {
		log.Println(err)
		return nil, nil
	}
	var doc_head DocHeader
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
	var url URL
	var response InsertResponse
	response.Id = -1
	url.Url = req.Document.Url
	raw_url, err := proto.Marshal(&url)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	raw_doc_id, _, err := db_url_to_id.Get(raw_url)

	if err != pebble.ErrNotFound {
		if err != nil {
			log.Println(err)
			response.Error = err.Error()
			response.State = 1
			return &response, nil
		}
		response.State = 2
		var doc_id DocID
		response.Error = "already existed"
		proto.Unmarshal(raw_doc_id, &doc_id)
		response.Id = doc_id.Id
		return &response, nil

	}
	req.Document.Id = atomic.AddInt64(&nextDocId, 1)
	var new_max_doc_id, current_doc_id DocID
	new_max_doc_id.Id = req.Document.Id
	req.Document.Id--
	current_doc_id.Id = req.Document.Id
	raw_doc_id, err = proto.Marshal(&new_max_doc_id)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	err = db_url_to_id.Set([]byte("MAX_DOC_ID_KEY"), raw_doc_id, pebble.Sync)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	raw_doc_id, err = proto.Marshal(&current_doc_id)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	err = db_url_to_id.Set(raw_url, raw_doc_id, pebble.Sync)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	docId.Id = req.Document.Id
	log.Println("To write id", docId.Id)
	raw_id, err := proto.Marshal(&docId)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	log.Printf("Current id = %d, url = %s", nextDocId, req.Document.Url)
	var header DocHeader
	header.Url = req.Document.Url
	header.Title = req.Document.Title
	header.Authors = req.Document.Authors
	header.CreateTime = req.Document.CreateTime
	header.UpdateTime = req.Document.UpdateTime
	header.Hash = req.Document.Hash
	raw_header, err := proto.Marshal(&header)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	raw_doc, err := proto.Marshal(req.Document)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	response.Id = docId.Id
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, nil
	}
	// If current doc has parent, update parent doc's childrend list
	if req.Document.ParentDocumentID != -1 {
		var parentId DocID
		parentId.Id = req.Document.ParentDocumentID
		rawParentId, _ := proto.Marshal(&parentId)
		rawParentDocChildList, _, err := db_children_id.Get(rawParentId)
		var docIdList DocIdList
		if err == pebble.ErrNotFound {

		} else if err != nil {

			log.Println(err)
			response.Error = err.Error()
			return &response, nil
		} else {
			proto.Unmarshal(rawParentDocChildList, &docIdList)
		}

		docIdList.DocIdList = append(docIdList.DocIdList, req.Document.Id)
	}
	err = db_head.Set(raw_id, raw_header, pebble.Sync)
	if err != nil {
		log.Println("Insert header error", err)
		response.Error = err.Error()
		return &response, nil
	}
	err = db_content.Set(raw_id, raw_doc, pebble.Sync)
	if err != nil {
		log.Println("Insert content error", err)
		response.Error = err.Error()
		return &response, nil
	}
	err = db_prepare_index.Set(raw_id, nil, pebble.Sync)
	if err != nil {
		log.Println("Insert prepare_index error", err)
		response.Error = err.Error()
		return &response, nil
	}
	return &response, nil
}

func (s *server) UpdateDocument(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDocument not implemented")
}

func (s *server) SearchForDocId(context.Context, *SearchRequest) (*DocIdList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}

func (s *server) SearchForHeader(context.Context, *SearchRequest) (*DocHeaderList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}

func (s *server) SearchForDocument(context.Context, *SearchRequest) (*DocList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}

func (s *server) DoCompact(ctx context.Context, in *DoCompactRequest) (*GeneralResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}

func (s *server) UpdateIndex(ctx context.Context, in *UpdateIndexRequest) (*GeneralResponse, error) {
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
	db_prepare_index, err = pebble.Open(*db_path+"/prepare_index", &pebble.Options{})
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
	db_children_id, err = pebble.Open(*db_path+"/children_id", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterSearchServiceServer(s, &server{})
	raw_doc_id, _, err := db_url_to_id.Get([]byte("MAX_DOC_ID_KEY"))

	if err != nil {
		log.Println(err)
		nextDocId = 0
		log.Println("nextDocId is 0")
	} else {
		var doc_id DocID
		proto.Unmarshal(raw_doc_id, &doc_id)
		nextDocId = doc_id.Id
		log.Println("nextDocId is ", nextDocId)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
