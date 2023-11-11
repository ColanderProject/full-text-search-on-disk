# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: search.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0csearch.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\x12\n\x03URL\x12\x0b\n\x03url\x18\x01 \x01(\t\"$\n\x12InvertedIndexValue\x12\x0e\n\x06\x64ocIds\x18\x01 \x03(\x03\"\x12\n\x10GetMaxURLRequest\"\x13\n\x05\x44ocID\x12\n\n\x02id\x18\x01 \x01(\x03\"@\n\x0f\x44ocumentRequest\x12\n\n\x02id\x18\x01 \x01(\x03\x12\x10\n\x08maxLevel\x18\x02 \x01(\x03\x12\x0f\n\x07maxSize\x18\x03 \x01(\x03\"\x10\n\x0eGetModeRequest\"\x12\n\x10\x44oCompactRequest\"\x14\n\x12UpdateIndexRequest\"\x1b\n\rHeaderRequest\x12\n\n\x02id\x18\x01 \x01(\x03\",\n\rInsertRequest\x12\x1b\n\x08\x64ocument\x18\x01 \x01(\x0b\x32\t.Document\"\x9b\x02\n\rSearchRequest\x12\x17\n\x0frequiredKeyword\x18\x01 \x03(\t\x12\x17\n\x0foptionalKeyword\x18\x02 \x03(\t\x12\x18\n\x10\x66orbiddenKeyword\x18\x03 \x03(\t\x12-\n\tstartTime\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12+\n\x07\x65ndTime\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x14\n\x0c\x64ocumentType\x18\x06 \x01(\t\x12\x17\n\x0frequireAbstract\x18\x07 \x01(\x08\x12#\n\x1bisContainsAuthorConstraints\x18\x08 \x01(\x08\x12\x0e\n\x06\x61uthor\x18\t \x03(\t\"(\n\x0bServerState\x12\x19\n\x04mode\x18\x01 \x01(\x0e\x32\x0b.ServerMode\"8\n\rUpdateRequest\x12\x1b\n\x08\x64ocument\x18\x01 \x01(\x0b\x32\t.Document\x12\n\n\x02id\x18\x02 \x01(\x03\"F\n\x0e\x44\x65leteResponse\x12\n\n\x02id\x18\x01 \x01(\x03\x12\x19\n\x05state\x18\x02 \x01(\x0e\x32\n.StateCode\x12\r\n\x05\x65rror\x18\x03 \x01(\t\"\xb8\x01\n\tDocHeader\x12\x0b\n\x03url\x18\x01 \x01(\t\x12\r\n\x05title\x18\x02 \x01(\t\x12\x0f\n\x07\x61uthors\x18\x03 \x03(\t\x12.\n\ncreateTime\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12.\n\nupdateTime\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x0c\n\x04hash\x18\x06 \x01(\x03\x12\x10\n\x08\x61\x62stract\x18\x07 \x01(\t\"\'\n\rDocHeaderList\x12\x16\n\x0e\x63hildrenDocIds\x18\x01 \x03(\x03\"4\n\x0e\x44ocumentIdList\x12\"\n\x0e\x63hildrenDocIds\x18\x01 \x03(\x0b\x32\n.DocHeader\"\xfd\x01\n\x08\x44ocument\x12\x0b\n\x03url\x18\x01 \x01(\t\x12\r\n\x05title\x18\x02 \x01(\t\x12\x0f\n\x07\x61uthors\x18\x03 \x03(\t\x12.\n\ncreateTime\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12.\n\nupdateTime\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x0c\n\x04hash\x18\x06 \x01(\x03\x12\x0c\n\x04\x64\x61ta\x18\x07 \x01(\t\x12\x0c\n\x04size\x18\x08 \x01(\x03\x12\n\n\x02id\x18\t \x01(\x03\x12\x18\n\x10parentDocumentID\x18\n \x01(\x03\x12\x14\n\x0c\x64ocumentType\x18\x0b \x01(\t\"F\n\x0eInsertResponse\x12\n\n\x02id\x18\x01 \x01(\x03\x12\x19\n\x05state\x18\x02 \x01(\x0e\x32\n.StateCode\x12\r\n\x05\x65rror\x18\x03 \x01(\t\"&\n\x07\x44ocList\x12\x1b\n\x08response\x18\x01 \x03(\x0b\x32\t.Document\"]\n\x0fSetModeResponse\x12\x19\n\x05state\x18\x01 \x01(\x0e\x32\n.StateCode\x12 \n\x0b\x63urrentMode\x18\x02 \x01(\x0e\x32\x0b.ServerMode\x12\r\n\x05\x65rror\x18\x03 \x01(\t\"F\n\x0eUpdateResponse\x12\n\n\x02id\x18\x01 \x01(\x03\x12\x19\n\x05state\x18\x02 \x01(\x0e\x32\n.StateCode\x12\r\n\x05\x65rror\x18\x03 \x01(\t\";\n\x0fGeneralResponse\x12\x19\n\x05state\x18\x01 \x01(\x0e\x32\n.StateCode\x12\r\n\x05\x65rror\x18\x02 \x01(\t*9\n\nServerMode\x12\r\n\tREAD_ONLY\x10\x00\x12\n\n\x06UPDATE\x10\x01\x12\x10\n\x0cINDEX_UPDATE\x10\x02*6\n\tStateCode\x12\x0b\n\x07SUCCESS\x10\x00\x12\x0b\n\x07\x46\x41ILURE\x10\x01\x12\x0f\n\x0bREPEATED_ID\x10\x02\x32\xd3\x05\n\rSearchService\x12\x18\n\x06URL2ID\x12\x04.URL\x1a\x06.DocID\"\x00\x12(\n\tGetMaxURL\x12\x11.GetMaxURLRequest\x1a\x06.DocID\"\x00\x12(\n\x11GetDocumentHeader\x12\x06.DocID\x1a\t.Document\"\x00\x12+\n\x16GetDocumentHeaderByURL\x12\x04.URL\x1a\t.Document\"\x00\x12\x30\n\x10GetChildDocument\x12\x10.DocumentRequest\x1a\x08.DocList\"\x00\x12*\n\x07GetMode\x12\x0f.GetModeRequest\x1a\x0c.ServerState\"\x00\x12+\n\x07SetMode\x12\x0c.ServerState\x1a\x10.SetModeResponse\"\x00\x12+\n\x0e\x44\x65leteDocument\x12\x06.DocID\x1a\x0f.DeleteResponse\"\x00\x12\x33\n\x0eInsertDocument\x12\x0e.InsertRequest\x1a\x0f.InsertResponse\"\x00\x12\x33\n\x0eUpdateDocument\x12\x0e.UpdateRequest\x1a\x0f.UpdateResponse\"\x00\x12\x33\n\x0eSearchForDocId\x12\x0e.SearchRequest\x1a\x0f.DocumentIdList\"\x00\x12\x33\n\x0fSearchForHeader\x12\x0e.SearchRequest\x1a\x0e.DocHeaderList\"\x00\x12/\n\x11SearchForDocument\x12\x0e.SearchRequest\x1a\x08.DocList\"\x00\x12\x32\n\tDoCompact\x12\x11.DoCompactRequest\x1a\x10.GeneralResponse\"\x00\x12\x36\n\x0bUpdateIndex\x12\x13.UpdateIndexRequest\x1a\x10.GeneralResponse\"\x00\x42*Z(github.com/qdmalekith/byrSearch-micro/pbb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'search_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z(github.com/qdmalekith/byrSearch-micro/pb'
  _globals['_SERVERMODE']._serialized_start=1685
  _globals['_SERVERMODE']._serialized_end=1742
  _globals['_STATECODE']._serialized_start=1744
  _globals['_STATECODE']._serialized_end=1798
  _globals['_URL']._serialized_start=49
  _globals['_URL']._serialized_end=67
  _globals['_INVERTEDINDEXVALUE']._serialized_start=69
  _globals['_INVERTEDINDEXVALUE']._serialized_end=105
  _globals['_GETMAXURLREQUEST']._serialized_start=107
  _globals['_GETMAXURLREQUEST']._serialized_end=125
  _globals['_DOCID']._serialized_start=127
  _globals['_DOCID']._serialized_end=146
  _globals['_DOCUMENTREQUEST']._serialized_start=148
  _globals['_DOCUMENTREQUEST']._serialized_end=212
  _globals['_GETMODEREQUEST']._serialized_start=214
  _globals['_GETMODEREQUEST']._serialized_end=230
  _globals['_DOCOMPACTREQUEST']._serialized_start=232
  _globals['_DOCOMPACTREQUEST']._serialized_end=250
  _globals['_UPDATEINDEXREQUEST']._serialized_start=252
  _globals['_UPDATEINDEXREQUEST']._serialized_end=272
  _globals['_HEADERREQUEST']._serialized_start=274
  _globals['_HEADERREQUEST']._serialized_end=301
  _globals['_INSERTREQUEST']._serialized_start=303
  _globals['_INSERTREQUEST']._serialized_end=347
  _globals['_SEARCHREQUEST']._serialized_start=350
  _globals['_SEARCHREQUEST']._serialized_end=633
  _globals['_SERVERSTATE']._serialized_start=635
  _globals['_SERVERSTATE']._serialized_end=675
  _globals['_UPDATEREQUEST']._serialized_start=677
  _globals['_UPDATEREQUEST']._serialized_end=733
  _globals['_DELETERESPONSE']._serialized_start=735
  _globals['_DELETERESPONSE']._serialized_end=805
  _globals['_DOCHEADER']._serialized_start=808
  _globals['_DOCHEADER']._serialized_end=992
  _globals['_DOCHEADERLIST']._serialized_start=994
  _globals['_DOCHEADERLIST']._serialized_end=1033
  _globals['_DOCUMENTIDLIST']._serialized_start=1035
  _globals['_DOCUMENTIDLIST']._serialized_end=1087
  _globals['_DOCUMENT']._serialized_start=1090
  _globals['_DOCUMENT']._serialized_end=1343
  _globals['_INSERTRESPONSE']._serialized_start=1345
  _globals['_INSERTRESPONSE']._serialized_end=1415
  _globals['_DOCLIST']._serialized_start=1417
  _globals['_DOCLIST']._serialized_end=1455
  _globals['_SETMODERESPONSE']._serialized_start=1457
  _globals['_SETMODERESPONSE']._serialized_end=1550
  _globals['_UPDATERESPONSE']._serialized_start=1552
  _globals['_UPDATERESPONSE']._serialized_end=1622
  _globals['_GENERALRESPONSE']._serialized_start=1624
  _globals['_GENERALRESPONSE']._serialized_end=1683
  _globals['_SEARCHSERVICE']._serialized_start=1801
  _globals['_SEARCHSERVICE']._serialized_end=2524
# @@protoc_insertion_point(module_scope)