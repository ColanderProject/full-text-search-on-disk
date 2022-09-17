package search_service

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestX1(t *testing.T) {
	br := bufio.NewReader(os.Stdin)
	x, _, _ := br.ReadLine()
	print(x)
	/*
		b := []byte("Hel\to")
		print(string(b))
		i := len(b)-1
		for i >= 0 {
			if b[i] == 9 {
				b[i] = 10
			}
			i--
		}

		print(string(b))
	*/
}

func TestX(t *testing.T) {
	bin_buf := bytes.Buffer{}
	var dd int
	dd = 1024
	binary.Write(&bin_buf, binary.BigEndian, int32(dd))
	bb := bin_buf.Bytes()
	fmt.Println(bb)

	//x := []byte{0,0,4,0,0,0,4,0}
	//bin_buf := bytes.NewBuffer(x)
	//y := (*int32)(unsafe.Pointer(&(x[4])))
	//fmt.Println(*y)
	//var zzz int32
	//binary.Read(bin_buf, binary.BigEndian, &zzz)
	//fmt.Println(zzz)
	//x := "你好"
	//for i, ch := range x {
	//	fmt.Println(i, string(ch))
	//}
}

func TestName2(t *testing.T) {
	//DB_C_INDEX, err := leveldb.OpenFile("C:\\Users\\v-xiaowensun\\data\\byr_search_leveldb\\"+"base_rindex/", nil)

	//x := "你好"
	//b := (*[]byte)(unsafe.Pointer(&x))
	//l := 0
	//for i, c := range(x) {
	//	fmt.Println(i, c, x[i], (*b)[l:i])
	//	l = i
	//}
	//fmt.Println((*b)[l:])
}

func TestName(t *testing.T) {
	//search()
	//r, i := Main("门头沟", "/home/sxw/jupyter/filehub_client/byr_search_leveldb/")
	r, i := Main("控制科学与工程", "C:\\Users\\v-xiaowensun\\data\\byr_search_leveldb\\", true)
	fmt.Println(i)
	f, err := os.Create("result.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write(r)
	//x := "你好"
	//b := (*[]byte)(unsafe.Pointer(&x))
	//l := 0
	//for i, c := range(x) {
	//	fmt.Println(i, c, x[i], (*b)[l:i])
	//	l = i
	//}
	//fmt.Println((*b)[l:])
}

func TestMainServer(t *testing.T) {
	MainServer("/home/sxw/jupyter/byr_search/byr_search_leveldb/", "127.0.0.1", "12345", "123")
	//MainServer("/home/sxw/jupyter/filehub_client/byr_search_leveldb/", "127.0.0.1", "12345", "123")
}
