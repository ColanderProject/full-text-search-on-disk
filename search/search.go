package search_service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
	"unsafe"
	_ "unsafe"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var DB_KEY_TO_ID *leveldb.DB // KEY to ID
var DB_H *leveldb.DB         // Document head
var DB_C *leveldb.DB         // Document content
var DB_H_INDEX *leveldb.DB   // Inverted index for content
var DB_C_INDEX *leveldb.DB   // Inverted index for content

var CONNECTION_KEY []byte // Use to validate a client when get a connection request
var BASE_DB_PATH string   // Path to db base directory

var exitKye = []byte("qC_welj!OIJ-Wsdfar") // Use to graceful shutdown the server

func init() {

}

func getContentIDs(key []byte, set *Set) {
	ids_b, err := DB_C_INDEX.Get(key, nil)
	if err != nil {
		log.Println(err, key)
		return
	}
	ids := bytes.Split(ids_b, []byte{9})
	for _, i := range ids {
		set.Add(*(*string)(unsafe.Pointer(&i)))
	}
}

// Perform a compaction
func compaction(db *leveldb.DB) {
	var islice util.Range
	err := db.CompactRange(islice)
	if err != nil {
		log.Println(err)
	}
}

func getHeadIDs(key []byte, set *Set) {
	ids_b, err := DB_H_INDEX.Get(key, nil)
	if err != nil {
		log.Println(err, key)
		return
	}
	ids := bytes.Split(ids_b, []byte{9})
	for _, i := range ids {
		set.Add(*(*string)(unsafe.Pointer(&i)))
	}
}

func searchByte(keyByte []byte, start_id int) []byte {
	if start_id < 0 {
		return nil
	}
	bt := time.Now().UnixNano()
	var set0 *Set
	var hset0 *Set
	set0 = nil
	hset0 = nil
	keys := bytes.Split(keyByte, []byte{32})
	//idsx := [][]byte{}
	// Get content document ids
	vis_map := map[int32]Empty{}
	for _, k_bytes := range keys {
		if len(k_bytes) <= 0 {
			continue
		}

		k := string(k_bytes)
		for _, ch := range k {
			_, ok := vis_map[ch]
			if ok {
				continue
			}
			vis_map[ch] = empty
			if set0 == nil {
				set0 = SetFactory()
				getContentIDs([]byte(string(ch)), set0)
			} else {
				set1 := SetFactory()
				getContentIDs([]byte(string(ch)), set1)
				set0.intersection(set1)
			}
			if set0.Len() == 0 {
				break
			}
		}
		if set0.Len() == 0 {
			break
		}
	}
	// Get title document ids
	vis_map = map[int32]Empty{}
	for _, k_bytes := range keys {
		k := string(k_bytes)
		for _, ch := range k {
			_, ok := vis_map[ch]
			if ok {
				continue
			}
			vis_map[ch] = empty
			if hset0 == nil {
				hset0 = SetFactory()
				getHeadIDs([]byte(string(ch)), hset0)
			} else {
				hset1 := SetFactory()
				getHeadIDs([]byte(string(ch)), hset1)
				hset0.intersection(hset1)
			}
			if hset0.Len() == 0 {
				break
			}
		}
	}

	// No document recalled
	if set0 == nil && hset0 == nil {
		return nil
	}

	var buffer bytes.Buffer
	target_ids := make([][2]int, 0)
	both := 0
	if set0 != nil {
		if hset0 != nil {
			for k := range set0.m {
				i, err := strconv.Atoi(k)
				if err != nil {
					log.Println(err)
				} else {
					_, ok := hset0.m[k]
					if ok {
						target_ids = append(target_ids, [2]int{i, 2})
						both++
					} else {
						target_ids = append(target_ids, [2]int{i, 0})
					}
				}
			}
		} else {
			for k := range set0.m {
				i, err := strconv.Atoi(k)
				if err != nil {
					log.Println(err)
				} else {
					target_ids = append(target_ids, [2]int{i, 0})
				}
			}
		}
	}

	if hset0 != nil {
		if set0 != nil {
			for k := range hset0.m {

				i, err := strconv.Atoi(k)
				if err == nil {
					_, ok := set0.m[k]
					if !ok {
						// head only
						target_ids = append(target_ids, [2]int{i, 1})
					}
				}
			}
		} else {
			for k := range hset0.m {
				i, err := strconv.Atoi(k)
				if err != nil {
					log.Println(err)
				} else {
					target_ids = append(target_ids, [2]int{i, 1})
				}
			}
		}
	}
	sort.Slice(target_ids, func(i, j int) bool {
		if target_ids[i][1] == target_ids[j][1] { // related score
			return target_ids[i][0] < target_ids[j][0] // document id
		}
		return target_ids[i][1] > target_ids[j][1]
	})
	i := start_id
	cnt := 0
	target_ids_len := len(target_ids)
	max_id := target_ids_len
	if max_id-i > 5000 {
		max_id = i + 5000
	}
	mt := time.Now().UnixNano()
	for i < max_id {
		var head, content []byte
		head = nil
		content = nil
		flag := true
		key := []byte(strconv.Itoa(target_ids[i][0]))
		if i%500 == 499 && time.Now().UnixNano()-mt > 1000000000 {
			break
		}
		if target_ids[i][1] == 2 || target_ids[i][1] == 1 {
			var err error
			head, err = DB_H.Get(key, nil)
			if err != nil {
				log.Println("Warning!! Document not found", target_ids[i])
				flag = false
			} else {
				for _, k := range keys {
					if len(k) < 2 {
						continue
					}
					if !bytes.Contains(head, k) {
						flag = false
						break
					}
				}
			}
			if !flag {
				i++
				continue
			} else {
				buffer.Write(key)
				buffer.Write([]byte{13})
				buffer.Write(head)
				buffer.Write([]byte{13})
				var err error
				content, err = DB_C.Get(key, nil)
				if err != nil {
					log.Println("Document not found", target_ids[i])
					content = []byte("")
				}
				if len(content) > 1000 {
					i := bytes.IndexByte(content, 10)
					if i > 0 && i < 1000 {
						buffer.Write(content[:i])
					}
				} else {
					buffer.Write(content)
				}
				buffer.Write([]byte{13})
				cnt++
				if cnt == 10 {
					buffer.Write([]byte(strconv.Itoa(target_ids_len)))
					buffer.Write([]byte{13})
					buffer.Write([]byte(strconv.Itoa(i + 1)))
					et := time.Now().UnixNano()
					log.Println(string(keyByte), hset0.Len(), set0.Len(), both, et-bt, mt-bt, start_id, i)
					return buffer.Bytes()
				}
				i++
				continue
			}
		}

		if target_ids[i][1] == 2 || target_ids[i][1] == 0 {
			var err error
			content, err = DB_C.Get(key, nil)
			if err != nil {
				log.Println("Document not found", target_ids[i])
				flag = false
			} else {
				for _, k := range keys {
					if len(k) < 2 {
						continue
					}
					if !bytes.Contains(content, k) {
						flag = false
						break
					}
				}
			}
		}
		if flag {
			buffer.Write(key)
			buffer.Write([]byte{13})
			head, _ = DB_H.Get(key, nil)
			buffer.Write(head)
			buffer.Write([]byte{13})
			if len(content) > 1000 {
				i := bytes.IndexByte(content, 10)
				if i > 0 && i < 1000 {
					buffer.Write(content[:i])
				}
			} else {
				buffer.Write(content)
			}
			buffer.Write([]byte{13})
			cnt++
			if cnt == 10 {
				buffer.Write([]byte(strconv.Itoa(target_ids_len)))
				buffer.Write([]byte{13})
				buffer.Write([]byte(strconv.Itoa(i + 1)))
				et := time.Now().UnixNano()
				log.Println(string(keyByte), hset0.Len(), set0.Len(), both, et-bt, mt-bt, start_id, i)
				return buffer.Bytes()
			}
		}
		i++
	}
	buffer.Write([]byte(strconv.Itoa(target_ids_len)))
	buffer.Write([]byte{13})
	if i >= len(target_ids) { // all document checked
		buffer.Write([]byte(strconv.Itoa(-1)))
	} else {
		buffer.Write([]byte(strconv.Itoa(i)))
	}
	et := time.Now().UnixNano()
	log.Println(string(keyByte), hset0.Len(), set0.Len(), both, et-bt, mt-bt, start_id, i)
	return buffer.Bytes()
}

func getBytesCharSet(byteString []byte, set0 *Set) {
	k := string(byteString)
	for _, chi := range k {
		ch := string(chi)
		if ch != " " && ch != "\t" && ch != "\n" && ch != "\r" && ch != "'" && ch != "\"" && ch != "\xa0" {
			set0.Add(ch)
		}
	}
}

func getContentCharSet(content []byte) *Set {
	set0 := SetFactory()
	for _, l := range bytes.Split(content, []byte{10}) {
		x := bytes.Split(l, []byte{9})
		i := 1
		for i < len(x) {
			getBytesCharSet(x[i], set0)
			i++
		}
	}
	return set0
}

func realReadLine(buf_reader *bufio.Reader) []byte {
	v, isPrefix, _ := buf_reader.ReadLine()
	//log.Println("Get", v, isPrefix)
	if !isPrefix {
		result := make([]byte, len(v))
		copy(result, v)
		return result
	}
	var buffer bytes.Buffer
	buffer.Write(v)
	for isPrefix {
		v, isPrefix, _ = buf_reader.ReadLine()
		buffer.Write(v)
	}
	return buffer.Bytes()
}

func validateConnection(c net.Conn) (int, *bufio.Reader) {
	c.SetReadDeadline(time.Now().Add(1 * time.Second)) // Set connection timeout to prevent connection request backlog
	bufReader := bufio.NewReader(c)
	key, _, _ := bufReader.ReadLine()
	if !bytes.Equal(key, CONNECTION_KEY) {
		if bytes.Equal(key, exitKye) {
			log.Println(c.RemoteAddr(), "Received exit command!")
			return 2, nil
		}
		if len(key) > 1000 {
			key = key[:1000]
		}
		log.Println(c.RemoteAddr(), "key mismatch", key)
		return 1, nil
	}
	log.Println("Get connect", c.RemoteAddr())
	c.SetReadDeadline(time.Time{})
	c.Write([]byte("ok\n"))
	return 0, bufReader
}

func handleConn(c net.Conn, buf_reader *bufio.Reader, debug bool) {
	defer c.Close()
	for {
		line, _, eof := buf_reader.ReadLine()
		if eof == io.EOF || len(line) == 0 {
			log.Println("finished", c.RemoteAddr())
			break
		}
		if line[0] == 115 { // 's' : search
			keys := realReadLine(buf_reader)
			i_byte, _, _ := buf_reader.ReadLine()
			i, err := strconv.Atoi(string(i_byte))
			if err != nil {
				log.Println("error start index", i_byte)
				return
			}
			//fmt.Println(string(keys), string(i_byte))
			r := searchByte(keys, i)
			c.Write([]byte(strconv.Itoa(len(r))))
			c.Write([]byte{10})
			c.Write(r)
		} else if line[0] == 103 { // 'g': get head
			id, _, _ := buf_reader.ReadLine()
			if len(id) == 0 {
				break
			}
			//println(string(id))
			head, err := DB_H.Get(id, nil)
			if err != nil {
				c.Write([]byte(strconv.Itoa(0)))
				c.Write([]byte{10})
			} else {
				c.Write([]byte(strconv.Itoa(len(head))))
				c.Write([]byte{10})
				c.Write(head)
			}
		} else if line[0] == 71 { // 'G': get content
			id, _, _ := buf_reader.ReadLine()
			if len(id) == 0 {
				break
			}
			//println(string(id))
			content, err := DB_C.Get(id, nil)
			if err != nil {
				c.Write([]byte(strconv.Itoa(0)))
				c.Write([]byte{10})
			} else {
				c.Write([]byte(strconv.Itoa(len(content))))
				c.Write([]byte{10})
				c.Write(content)
			}
		} else if line[0] == 117 { // 'u': update
			log.Println("Update mode")
			indexChangeH := map[string]IndexChange{}
			indexChangeC := map[string]IndexChange{}
			delHeadBytes := 0
			delContentBytes := 0
			for {
				buf, _, _ := buf_reader.ReadLine()
				if len(buf) == 0 {
					break
				}
				id := make([]byte, len(buf))
				copy(id, buf)
				idi, err := strconv.Atoi(string(id))
				if err != nil {
					log.Println("Error! bad id: ", string(id))
					return
				}
				if debug {
					log.Println("Get Id", id)
				}
				buf, _, _ = buf_reader.ReadLine()
				if buf[0] == 100 { // 'd': delete
					if debug {
						log.Println("To delete", id)
					}
					old_head, err := DB_H.Get(id, nil)
					if err != nil {
						fmt.Println("delete record head not found", id)
						continue
					}
					headCharset := SetFactory()
					getBytesCharSet(old_head, headCharset)
					for ch := range headCharset.m {
						//if ch != " " && ch != "\t" && ch != "\n" && ch != "\r" && ch != "'" && ch != "\"" && ch != "\xa0" {
						ic, ok := indexChangeH[ch]
						if !ok {
							ic = IndexChange{[]int{}, map[int]Empty{}}
						}
						ic.Remove[idi] = empty
						indexChangeH[ch] = ic
						//}
					}
					delHeadBytes += len(old_head)
					old_content, err := DB_C.Get(id, nil)
					set_cont_old := getContentCharSet(old_content)
					for ch := range set_cont_old.m {
						//if ch != " " && ch != "\t" && ch != "\n" && ch != "\r" && ch != "'" && ch != "\"" && ch != "\xa0" {
						ic, ok := indexChangeC[ch]
						if !ok {
							ic = IndexChange{[]int{}, map[int]Empty{}}
						}
						ic.Remove[idi] = empty
						indexChangeC[ch] = ic
						//}
					}
					err = DB_H.Delete(id, nil)
					if err != nil {
						log.Println("Error! Delete head, id", err, id)
					}
					delContentBytes += len(old_content)
					err = DB_C.Delete(id, nil)
					if err != nil {
						log.Println("Error! Delete content, id", err, id)
					}
					fd, err := os.OpenFile(BASE_DB_PATH+"byr_search_update.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
					if err != nil {
						log.Println("Error! write log, id", err, id)
					}
					fd.Write([]byte("delete id\t"))
					fd.Write(id)
					fd.Write([]byte("\t"))

					fd.Write([]byte{10})
					fd.Close()
				} else if buf[0] == 117 { // 'u': update
					if debug {
						log.Println("To update", id)
					}
					head := realReadLine(buf_reader)
					if debug {
						log.Println("Head is", head)
					}
					content := realReadLine(buf_reader)
					i := len(content) - 1
					for i >= 0 {
						if content[i] == 7 {
							content[i] = 10
						}
						i--
					}
					if debug {
						log.Println("Content is", content)
					}
					old_head, err := DB_H.Get(id, nil)
					if err != nil {
						//fmt.Println("New record", idi, id)
						// new record
						//DB_H.Put(id, head, nil)
						//DB_C.Put(id, content, nil)
						set_head := SetFactory()
						getBytesCharSet(bytes.Split(head, []byte{9})[0], set_head)
						set_cont := getContentCharSet(content)
						for ch := range set_head.m {
							ic, ok := indexChangeH[ch]
							if !ok {
								ic = IndexChange{[]int{}, map[int]Empty{}}
							}
							ic.Add = append(ic.Add, idi)
							indexChangeH[ch] = ic
						}
						for ch := range set_cont.m {
							ic, ok := indexChangeC[ch]
							if !ok {
								ic = IndexChange{[]int{}, map[int]Empty{}}
							}
							ic.Add = append(ic.Add, idi)
							indexChangeC[ch] = ic
						}
						fd, err := os.OpenFile(BASE_DB_PATH+"byr_search_update.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
						if err != nil {
							log.Println("Error! write log, id", err, id)
						}
						fd.Write([]byte("add id "))
						fd.Write(id)
						fd.Write([]byte{10})
						fd.Close()
					} else {
						old_content, err := DB_C.Get(id, nil)
						if err != nil {
							log.Println("Error!!! Content does not exist")
							old_content = []byte("")
						}
						set_head_old := SetFactory()
						getBytesCharSet(bytes.Split(old_head, []byte{9})[0], set_head_old)
						set_cont_old := getContentCharSet(old_content)

						set_head := SetFactory()
						getBytesCharSet(bytes.Split(head, []byte{9})[0], set_head)
						set_cont := getContentCharSet(content)
						for ch := range set_head.m {
							_, ok := set_head_old.m[ch]
							if !ok {
								ic, ok := indexChangeH[ch]
								if !ok {
									ic = IndexChange{[]int{}, map[int]Empty{}}
								}
								ic.Add = append(ic.Add, idi)
								indexChangeH[ch] = ic
							}
						}
						for ch := range set_head_old.m {
							_, ok := set_head.m[ch]
							if !ok {
								ic, ok := indexChangeH[ch]
								if !ok {
									ic = IndexChange{[]int{}, map[int]Empty{}}
								}
								ic.Remove[idi] = empty
								indexChangeH[ch] = ic
							}
						}

						for ch := range set_cont.m {
							_, ok := set_cont_old.m[ch]
							if !ok {
								ic, ok := indexChangeC[ch]
								if !ok {
									ic = IndexChange{[]int{}, map[int]Empty{}}
								}
								ic.Add = append(ic.Add, idi)
								indexChangeC[ch] = ic
							}
						}
						for ch := range set_cont_old.m {
							_, ok := set_cont.m[ch]
							if !ok {
								ic, ok := indexChangeC[ch]
								if !ok {
									ic = IndexChange{[]int{}, map[int]Empty{}}
								}
								ic.Remove[idi] = empty
								indexChangeC[ch] = ic
							}
						}
						fd, err := os.OpenFile(BASE_DB_PATH+"byr_search_update.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
						if err != nil {
							log.Println("Error! write log, id", err, id)
						}
						fd.Write([]byte("update id "))
						fd.Write(id)
						fd.Write([]byte{10})
						fd.Close()
					}
					err = DB_H.Put(id, head, nil)
					if err != nil {
						log.Fatal("Error! Write Head data", id)
					}
					err = DB_C.Put(id, content, nil)
					if err != nil {
						log.Fatal("Error! Write content data", id)
					}
				} else {
					log.Fatal("Unknown command!!!")
				}
			}
			log.Println("Head change", len(indexChangeH), delHeadBytes/1024, "KB Removed")
			for ch := range indexChangeH {
				old_ri, err := DB_H_INDEX.Get([]byte(ch), nil)
				outbuf := bytes.Buffer{}
				if err != nil {
					log.Println("New Char in head", ch)
					if len(indexChangeH[ch].Add) == 0 {
						log.Println("No change skipped!!")
						continue
					}
				} else {
					for _, x := range bytes.Split(old_ri, []byte{9}) {
						xi, err := strconv.Atoi(string(x))
						if err != nil {
							log.Println("Error rindex in head", ch, x, "Removed!")
						} else {
							_, ok := indexChangeH[ch].Remove[xi]
							if !ok {
								if outbuf.Len() != 0 {
									outbuf.Write([]byte{9})
								}
								outbuf.Write(x)
							}
						}
					}
				}
				for _, xi := range indexChangeH[ch].Add {
					if outbuf.Len() != 0 {
						outbuf.Write([]byte{9})
					}
					outbuf.Write([]byte(strconv.Itoa(xi)))
				}
				err = DB_H_INDEX.Put([]byte(ch), outbuf.Bytes(), nil)
				if err != nil {
					log.Fatal("Write Error!!!! head ri", ch, err)
				}
			}
			log.Println("Content change", len(indexChangeC), delContentBytes/1024, "KB Removed")
			ii := 0
			for ch := range indexChangeC {
				old_ri, err := DB_C_INDEX.Get([]byte(ch), nil)
				outbuf := bytes.Buffer{}
				if err != nil {
					log.Println("New Char in content", ch)
					if len(indexChangeC[ch].Add) == 0 {
						log.Println("No change skipped!!")
						continue
					}
				} else {
					for _, x := range bytes.Split(old_ri, []byte{9}) {
						xi, err := strconv.Atoi(string(x))
						if err != nil {
							log.Println("Error rindex in content", ch, x, "Removed!")
						} else {
							_, ok := indexChangeC[ch].Remove[xi]
							if !ok {
								if outbuf.Len() != 0 {
									outbuf.Write([]byte{9})
								}
								outbuf.Write(x)
							}
						}
					}
				}
				for _, xi := range indexChangeC[ch].Add {
					if outbuf.Len() != 0 {
						outbuf.Write([]byte{9})
					}
					outbuf.Write([]byte(strconv.Itoa(xi)))
				}
				err = DB_C_INDEX.Put([]byte(ch), outbuf.Bytes(), nil)
				if err != nil {
					log.Fatal("Write Error!!!!", ch, err)
				}
				ii += 1
				//fmt.Print("\r", ii, len(indexChangeC))
			}
			// DB_C.CompactRange(util.Range{})
			log.Println("Finished!!!")
			c.Write([]byte("ok\n"))
		} else if line[0] == 99 { // c CompactRange
			log.Println("start compact range DB_H")
			compaction(DB_H)
			log.Println("finish compact range DB_H")

			log.Println("start compact range DB_C")
			compaction(DB_C)
			log.Println("finish compact range DB_C")

			log.Println("start compact range DB_C_INDEX")
			compaction(DB_C_INDEX)
			log.Println("finish compact range DB_C_INDEX")

			log.Println("start compact range DB_H_INDEX")
			compaction(DB_H_INDEX)
			log.Println("finish compact range DB_H_INDEX")
			c.Write([]byte("ok\n"))
		} else if line[0] == 114 { // r: remote connect as client
			buf, _, _ := buf_reader.ReadLine()
			ip := make([]byte, len(buf))
			copy(ip, buf)
			buf, _, _ = buf_reader.ReadLine()
			port := make([]byte, len(buf))
			copy(port, buf)
			buf, _, _ = buf_reader.ReadLine()
			key := make([]byte, len(buf))
			copy(key, buf)

			ip_port := string(ip) + ":" + string(port)
			log.Println("Connect to: ", ip_port, string(key))
			conn, _ := net.Dial("tcp", ip_port)
			conn.Write(key)
			conn.Write([]byte{10})
			bufReader := bufio.NewReader(conn)
			message, _ := bufReader.ReadString('\n')
			if message[0] == 111 && message[1] == 107 { // ok
				go handleConn(conn, bufReader, false)
				c.Write([]byte("ok\n"))
			}
			c.Write([]byte("error\n"))
		}
	}
}

func MainServer(dbbase string, ip string, port string, key string) {
	var err error
	BASE_DB_PATH = dbbase
	CONNECTION_KEY = []byte(key)

	DB_KEY_TO_ID, err = leveldb.OpenFile(dbbase+"key2id/", nil)
	if err != nil {
		log.Fatal(err)
	}

	headOption := opt.Options{
		BlockCacheCapacity: 50 * 1048576, // 50 MB
	}
	DB_H, err = leveldb.OpenFile(dbbase+"head/", &headOption)
	if err != nil {
		log.Fatal(err)
	}

	contentOption := opt.Options{
		BlockCacheCapacity: 100 * 1048576, // 100 MB
	}
	DB_C, err = leveldb.OpenFile(dbbase+"content/", &contentOption)
	if err != nil {
		log.Fatal(err)
	}

	contentIndexOption := opt.Options{
		BlockCacheCapacity: 100 * 1048576, // 100 MB
	}
	DB_C_INDEX, err = leveldb.OpenFile(dbbase+"base_rindex/", &contentIndexOption)
	if err != nil {
		log.Fatal(err)
	}

	headIndexOption := opt.Options{
		BlockCacheCapacity: 50 * 1048576, // 50 MB
	}
	DB_H_INDEX, err = leveldb.OpenFile(dbbase+"base_rhindex/", &headIndexOption)
	if err != nil {
		log.Fatal(err)
	}

	defer DB_KEY_TO_ID.Close()
	defer DB_C.Close()
	defer DB_H.Close()
	defer DB_C_INDEX.Close()
	defer DB_H_INDEX.Close()

	l, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Start Listen on", ip+":"+port)
	defer l.Close()
	if err != nil {
		log.Fatal("listen error:", err)
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		state, bufReader := validateConnection(c)
		if state == 0 {
			go handleConn(c, bufReader, false)
		} else {
			c.Close()
			if state == 2 {
				break
			}
		}
	}
	return
}
