
# FullTextSearchOnDisk

## Build and Run

Get code

```
git clone https://github.com/sxwxs/FullTextSearchOnDisk.git
```

Build

```
cd FullTextSearchOnDisk
cd cmd
go build
```

Run

```
cmd db_path ip port key
# cmd E:\sxw\code\fileTrans\byr_search_leveldb_2022_03_28\ 127.0.0.1 8000 1234
```

## Overview

This project aims to build a light weight full text search engine which can work with limited memory.

When the main server started. It will:

1. Open four KV-store database on disk(goleveldb for now)
	- Head DB
	- Content DB
	- Head Index DB
	- Content Index DB
2. Listen on a TCP port and provide some commands:
   Client must send the connection key when the TCP connection established for authentication.
	- Graceful exit: send exit key when conection
	- Search for a key workd: command `s`
	- Get document head by document ID: command `g`
	- Get document content by document ID: command `G`
	- Update or insert a document
	   When update or insert or delete documents. First send `u` to switch to Update mode.
	   In Update mode, first send an ID, then use `d` to delete, `u` to update, send empty line to exit update mode. After exiting, the index will begin to update, and when the index update is complete, a compaction will be performed.

## How does search and index work

1. Document Retrieval

If we have two documents: 
> Document id: 0
> Title: "欢迎你"
> Content: "欢迎光临我们的网站"
> 
> Document id: 1
> Title: "欢迎光临"
> Content: "新店开业"

We will create character level inverted index， for head index:

> 欢： [0, 1]
> 迎： [0, 1]
> 你： [0 ]
> 光： [1 ]
> 临： [1 ]

For content index:

> 欢： [0 ]
> 迎： [0 ]
> 光： [0 ]
> 临： [0 ]
> 我： [0 ]
> 们： [0 ]
> 的： [0 ]
> 网： [0 ]
> 站： [0 ]
> 新： [1 ]
> 店： [1 ]
> 开： [1 ]
> 业： [1 ]

When get a key word, we fist go to index to get possible document ids.

For example the key word is "新开". We cannot get mathc in head index, but can get `新： [1 ]` and `开： [1 ]`. Will take intersection on them, then read document from content db (because this document is only matched by content index) to see if the key work "新开" in the content. The key word not in "新店开业", so we will not recall this document.

2. Ranking

Currently we have two field for ranking:
1. Matching score: 
   - Score 2: Key word in both title and content
   - Score 1: Key word only in title
   - Score 0: Key word only in content
2. Document ID

## Next step

1. Support arbitrary byte as key
   Create a key2id db, for head, content, head index, content index table we still use int as id, and we store key2id map in the new db. 
2. Provide HTTP and GRPC interface instead TCP interface.
    We need interfaces as follow:
    - Search(keyWord, maxReturnDocumentCount, ContinueID, maxSearchTime)
    - Exist(documentID): check if a ducument id exist
    - Exist(key): check if a key exist
    - Upsert(key, head, content): will log the operation in cache
    - Delete(key): will log the operation in cache
    - Commit(): All Upsert and delete will not be write before a commit call
    - GetCurrentCache(): Get current cache content
3. Support index aggration
4. Use protobuf to save id list for inverted index
