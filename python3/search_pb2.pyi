from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ServerMode(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    READ_ONLY: _ClassVar[ServerMode]
    UPDATE: _ClassVar[ServerMode]
    INDEX_UPDATE: _ClassVar[ServerMode]

class StateCode(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    SUCCESS: _ClassVar[StateCode]
    FAILURE: _ClassVar[StateCode]
    REPEATED_ID: _ClassVar[StateCode]
READ_ONLY: ServerMode
UPDATE: ServerMode
INDEX_UPDATE: ServerMode
SUCCESS: StateCode
FAILURE: StateCode
REPEATED_ID: StateCode

class URL(_message.Message):
    __slots__ = ["url"]
    URL_FIELD_NUMBER: _ClassVar[int]
    url: str
    def __init__(self, url: _Optional[str] = ...) -> None: ...

class InvertedIndexValue(_message.Message):
    __slots__ = ["docIds"]
    DOCIDS_FIELD_NUMBER: _ClassVar[int]
    docIds: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, docIds: _Optional[_Iterable[int]] = ...) -> None: ...

class GetMaxURLRequest(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class DocID(_message.Message):
    __slots__ = ["id"]
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class DocumentRequest(_message.Message):
    __slots__ = ["id", "maxLevel", "maxSize"]
    ID_FIELD_NUMBER: _ClassVar[int]
    MAXLEVEL_FIELD_NUMBER: _ClassVar[int]
    MAXSIZE_FIELD_NUMBER: _ClassVar[int]
    id: int
    maxLevel: int
    maxSize: int
    def __init__(self, id: _Optional[int] = ..., maxLevel: _Optional[int] = ..., maxSize: _Optional[int] = ...) -> None: ...

class GetModeRequest(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class DoCompactRequest(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class UpdateIndexRequest(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class HeaderRequest(_message.Message):
    __slots__ = ["id"]
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class InsertRequest(_message.Message):
    __slots__ = ["document"]
    DOCUMENT_FIELD_NUMBER: _ClassVar[int]
    document: Document
    def __init__(self, document: _Optional[_Union[Document, _Mapping]] = ...) -> None: ...

class SearchRequest(_message.Message):
    __slots__ = ["requiredKeyword", "optionalKeyword", "forbiddenKeyword", "startTime", "endTime", "documentType", "requireAbstract", "isContainsAuthorConstraints", "author"]
    REQUIREDKEYWORD_FIELD_NUMBER: _ClassVar[int]
    OPTIONALKEYWORD_FIELD_NUMBER: _ClassVar[int]
    FORBIDDENKEYWORD_FIELD_NUMBER: _ClassVar[int]
    STARTTIME_FIELD_NUMBER: _ClassVar[int]
    ENDTIME_FIELD_NUMBER: _ClassVar[int]
    DOCUMENTTYPE_FIELD_NUMBER: _ClassVar[int]
    REQUIREABSTRACT_FIELD_NUMBER: _ClassVar[int]
    ISCONTAINSAUTHORCONSTRAINTS_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    requiredKeyword: _containers.RepeatedScalarFieldContainer[str]
    optionalKeyword: _containers.RepeatedScalarFieldContainer[str]
    forbiddenKeyword: _containers.RepeatedScalarFieldContainer[str]
    startTime: _timestamp_pb2.Timestamp
    endTime: _timestamp_pb2.Timestamp
    documentType: str
    requireAbstract: bool
    isContainsAuthorConstraints: bool
    author: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, requiredKeyword: _Optional[_Iterable[str]] = ..., optionalKeyword: _Optional[_Iterable[str]] = ..., forbiddenKeyword: _Optional[_Iterable[str]] = ..., startTime: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., endTime: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., documentType: _Optional[str] = ..., requireAbstract: bool = ..., isContainsAuthorConstraints: bool = ..., author: _Optional[_Iterable[str]] = ...) -> None: ...

class ServerState(_message.Message):
    __slots__ = ["mode"]
    MODE_FIELD_NUMBER: _ClassVar[int]
    mode: ServerMode
    def __init__(self, mode: _Optional[_Union[ServerMode, str]] = ...) -> None: ...

class UpdateRequest(_message.Message):
    __slots__ = ["document", "id"]
    DOCUMENT_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    document: Document
    id: int
    def __init__(self, document: _Optional[_Union[Document, _Mapping]] = ..., id: _Optional[int] = ...) -> None: ...

class DeleteResponse(_message.Message):
    __slots__ = ["id", "state", "error"]
    ID_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    id: int
    state: StateCode
    error: str
    def __init__(self, id: _Optional[int] = ..., state: _Optional[_Union[StateCode, str]] = ..., error: _Optional[str] = ...) -> None: ...

class DocHeader(_message.Message):
    __slots__ = ["url", "title", "authors", "createTime", "updateTime", "hash", "abstract"]
    URL_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    AUTHORS_FIELD_NUMBER: _ClassVar[int]
    CREATETIME_FIELD_NUMBER: _ClassVar[int]
    UPDATETIME_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    ABSTRACT_FIELD_NUMBER: _ClassVar[int]
    url: str
    title: str
    authors: _containers.RepeatedScalarFieldContainer[str]
    createTime: _timestamp_pb2.Timestamp
    updateTime: _timestamp_pb2.Timestamp
    hash: int
    abstract: str
    def __init__(self, url: _Optional[str] = ..., title: _Optional[str] = ..., authors: _Optional[_Iterable[str]] = ..., createTime: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., updateTime: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., hash: _Optional[int] = ..., abstract: _Optional[str] = ...) -> None: ...

class DocHeaderList(_message.Message):
    __slots__ = ["childrenDocIds"]
    CHILDRENDOCIDS_FIELD_NUMBER: _ClassVar[int]
    childrenDocIds: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, childrenDocIds: _Optional[_Iterable[int]] = ...) -> None: ...

class DocumentIdList(_message.Message):
    __slots__ = ["childrenDocIds"]
    CHILDRENDOCIDS_FIELD_NUMBER: _ClassVar[int]
    childrenDocIds: _containers.RepeatedCompositeFieldContainer[DocHeader]
    def __init__(self, childrenDocIds: _Optional[_Iterable[_Union[DocHeader, _Mapping]]] = ...) -> None: ...

class Document(_message.Message):
    __slots__ = ["url", "title", "authors", "createTime", "updateTime", "hash", "data", "size", "id", "parentDocumentID", "documentType"]
    URL_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    AUTHORS_FIELD_NUMBER: _ClassVar[int]
    CREATETIME_FIELD_NUMBER: _ClassVar[int]
    UPDATETIME_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    SIZE_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    PARENTDOCUMENTID_FIELD_NUMBER: _ClassVar[int]
    DOCUMENTTYPE_FIELD_NUMBER: _ClassVar[int]
    url: str
    title: str
    authors: _containers.RepeatedScalarFieldContainer[str]
    createTime: _timestamp_pb2.Timestamp
    updateTime: _timestamp_pb2.Timestamp
    hash: int
    data: str
    size: int
    id: int
    parentDocumentID: int
    documentType: str
    def __init__(self, url: _Optional[str] = ..., title: _Optional[str] = ..., authors: _Optional[_Iterable[str]] = ..., createTime: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., updateTime: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., hash: _Optional[int] = ..., data: _Optional[str] = ..., size: _Optional[int] = ..., id: _Optional[int] = ..., parentDocumentID: _Optional[int] = ..., documentType: _Optional[str] = ...) -> None: ...

class InsertResponse(_message.Message):
    __slots__ = ["id", "state", "error"]
    ID_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    id: int
    state: StateCode
    error: str
    def __init__(self, id: _Optional[int] = ..., state: _Optional[_Union[StateCode, str]] = ..., error: _Optional[str] = ...) -> None: ...

class DocList(_message.Message):
    __slots__ = ["response"]
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    response: _containers.RepeatedCompositeFieldContainer[Document]
    def __init__(self, response: _Optional[_Iterable[_Union[Document, _Mapping]]] = ...) -> None: ...

class SetModeResponse(_message.Message):
    __slots__ = ["state", "currentMode", "error"]
    STATE_FIELD_NUMBER: _ClassVar[int]
    CURRENTMODE_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    state: StateCode
    currentMode: ServerMode
    error: str
    def __init__(self, state: _Optional[_Union[StateCode, str]] = ..., currentMode: _Optional[_Union[ServerMode, str]] = ..., error: _Optional[str] = ...) -> None: ...

class UpdateResponse(_message.Message):
    __slots__ = ["id", "state", "error"]
    ID_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    id: int
    state: StateCode
    error: str
    def __init__(self, id: _Optional[int] = ..., state: _Optional[_Union[StateCode, str]] = ..., error: _Optional[str] = ...) -> None: ...

class GeneralResponse(_message.Message):
    __slots__ = ["state", "error"]
    STATE_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    state: StateCode
    error: str
    def __init__(self, state: _Optional[_Union[StateCode, str]] = ..., error: _Optional[str] = ...) -> None: ...
