import search_pb2
import search_pb2_grpc
import grpc
from google.protobuf import timestamp_pb2
import time
['DESCRIPTOR', 'DeleteResponse', 'DoCompactRequest', 'DocHeader', 'DocHeaderList', 'DocID', 'DocList', 'Document', 'DocumentIdList', 'DocumentRequest', 'FAILURE', 'GeneralResponse', 'GetMaxURLRequest', 'GetModeRequest', 'HeaderRequest', 'INDEX_UPDATE', 'InsertRequest', 'InsertResponse', 'InvertedIndexValue', 'READ_ONLY', 'REPEATED_ID', 'SUCCESS', 'SearchRequest', 'ServerMode', 'ServerState', 'SetModeResponse', 'StateCode', 'UPDATE', 'URL', 'UpdateIndexRequest', 'UpdateRequest', 'UpdateResponse', '_DELETERESPONSE', '_DOCHEADER', '_DOCHEADERLIST', '_DOCID', '_DOCLIST', '_DOCOMPACTREQUEST', '_DOCUMENT', '_DOCUMENTIDLIST', '_DOCUMENTREQUEST', '_GENERALRESPONSE', '_GETMAXURLREQUEST', '_GETMODEREQUEST', '_HEADERREQUEST', '_INSERTREQUEST', '_INSERTRESPONSE', '_INVERTEDINDEXVALUE', '_SEARCHREQUEST', '_SEARCHSERVICE', '_SERVERMODE', '_SERVERSTATE', '_SETMODERESPONSE', '_STATECODE', '_UPDATEINDEXREQUEST', '_UPDATEREQUEST', '_UPDATERESPONSE', '_URL', '__builtins__', '__cached__', '__doc__', '__file__', '__loader__', '__name__', '__package__', '__spec__', '_builder', '_descriptor', '_descriptor_pool', '_globals', '_sym_db', '_symbol_database', 'google_dot_protobuf_dot_timestamp__pb2']
['SearchService', 'SearchServiceServicer', 'SearchServiceStub', '__builtins__', '__cached__', '__doc__', '__file__', '__loader__', '__name__', '__package__', '__spec__', 'add_SearchServiceServicer_to_server', 'grpc', 'search__pb2']

['DeleteDocument', 'DoCompact', 'GetChildDocument', 'GetDocumentHeader', 'GetDocumentHeaderByURL', 'GetMaxURL', 'GetMode', 'InsertDocument', 'SearchForDocId', 'SearchForDocument', 'SearchForHeader', 'SetMode', 'URL2ID', 'UpdateDocument', 'UpdateIndex', '__class__', '__delattr__', '__dict__', '__dir__', '__doc__', '__eq__', '__format__', '__ge__', '__getattribute__', '__gt__', '__hash__', '__init__', '__init_subclass__', '__le__', '__lt__', '__module__', '__ne__', '__new__', '__reduce__', '__reduce_ex__', '__repr__', '__setattr__', '__sizeof__', '__str__', '__subclasshook__', '__weakref__']

doc1 = search_pb2.Document(
    url="https://bbs.xxx.com/threads/1/",
    title='Welcom to bbs!',
    authors=['admin'],
    createTime=timestamp_pb2.Timestamp(
                    seconds=int(time.time()), 
                ),
    data='Hello everyone, welcom to this bbs!',
    parentDocumentID=-1 # means this is the root document
)
doc2 = search_pb2.Document(
    url="https://bbs.xxx.com/threads/1/#1",
    title='Re: Welcom to bbs!',
    authors=['admin'],
    createTime=timestamp_pb2.Timestamp(
                    seconds=int(time.time()), 
                ),
    data='Greeting!',
    parentDocumentID=-1 # means this is the root document
)
doc3 = search_pb2.Document(
    url="https://bbs.xxx.com/threads/1/#2",
    title='Re: Welcom to bbs!',
    authors=['admin'],
    createTime=timestamp_pb2.Timestamp(
                    seconds=int(time.time()), 
                ),
    data='Nice to meet you!',
    parentDocumentID=0
)


channel = grpc.insecure_channel('localhost:50051')
# test insert
stub = search_pb2_grpc.SearchServiceStub(channel)
ins = search_pb2.InsertRequest(document=doc1)
response = stub.InsertDocument(ins)
print("client received: " + response.error)
print("Doc id: " , response.id)
doc2.parentDocumentID = response.id
doc3.parentDocumentID = response.id
ins2 = search_pb2.InsertRequest(document=doc2)
ins3 = search_pb2.InsertRequest(document=doc3)
response = stub.InsertDocument(ins2)
print("client received: " + response.error)
print("Doc id: ", response.id)
response = stub.InsertDocument(ins3)
print("client received: " + response.error)
print("Doc id: ", response.id)
rr0 = stub.GetDocumentHeader(search_pb2.DocID(id=0))
assert rr0.title == 'Welcom to bbs!'
assert rr0.url == 'https://bbs.xxx.com/threads/1/'
rr1 = stub.GetDocumentHeader(search_pb2.DocID(id=1))
assert rr1.title == 'Re: Welcom to bbs!'
assert rr1.url == 'https://bbs.xxx.com/threads/1/#1'
rr2 = stub.GetDocumentHeader(search_pb2.DocID(id=2))
assert rr2.title == 'Re: Welcom to bbs!'
assert rr2.url == 'https://bbs.xxx.com/threads/1/#2'
# stub.GetChildDocId(search_pb2.DocID(id=0))
