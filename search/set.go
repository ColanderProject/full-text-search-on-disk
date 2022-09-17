package search_service

import (
	"sort"
	"strconv"
)

// https://blog.csdn.net/haodawang/article/details/80006059
type Empty struct { }

var empty Empty

type Set struct {
	m map[string]Empty
}

func SetFactory() *Set{
	return &Set{
		m:map[string]Empty{},
	}
}

//添加元素
func (s *Set) Add(val string) {
	s.m[val] = empty
}

//删除元素
func (s *Set) Remove(val string) {
	delete(s.m, val)
}

//获取长度
func (s *Set) Len() int {
	return len(s.m)
}

//清空set
func (s *Set) Clear() {
	s.m = make(map[string]Empty)
}

//排序输出
func (s *Set) GetItems() []string {
	vals := make([]string, 0, s.Len())

	for v := range s.m {
		vals = append(vals, v)
	}

	//排序
	sort.Slice(vals, func(i, j int) bool {
		numA, _ := strconv.Atoi(vals[i])
		numB, _ := strconv.Atoi(vals[j])
		return numA < numB
	})

	return vals
}

func (s *Set)intersection(s2 *Set) {
	for k, _ := range s.m {
		_, ok := s2.m[k]
		if !ok {
			s.Remove(k)
		}
	}
}