package tools

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

func ImplodeInt(sep string, items ...int) string {
	if len(items) < 1 {
		return ""
	}
	tmp := make([]string, len(items))
	for k, v := range items {
		tmp[k] = strconv.FormatInt(int64(v), 10)
	}
	return strings.Join(tmp, sep)
}

func ImplodeUint(sep string, items ...uint) string {
	if len(items) < 1 {
		return ""
	}
	tmp := make([]string, len(items))
	for k, v := range items {
		tmp[k] = strconv.FormatInt(int64(v), 10)
	}
	return strings.Join(tmp, sep)
}

func ExplodeInt(sep, s string) []int {
	if len(s) < 1 {
		return []int{}
	}
	data := strings.Split(s, sep)
	tmp := make([]int, 0, len(data))
	for _, v := range data {
		item, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			continue
		}
		tmp = append(tmp, int(item))
	}
	return tmp
}

func ExplodeUint(sep, s string) []uint {
	if len(s) < 1 {
		return []uint{}
	}
	data := strings.Split(s, sep)
	tmp := make([]uint, 0, len(data))
	for _, v := range data {
		item, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			continue
		}
		tmp = append(tmp, uint(item))
	}
	return tmp
}

func IntArrayDel(a []int, key int) []int {
	if key < 0 {
		key = len(a) + key
	}
	if key <= 0 {
		return a
	}
	return append(a[:key], a[key+1:]...)
}

func UintArrayDel(a []uint, key int) []uint {
	if key < 0 {
		key = len(a) + key
	}
	if key <= 0 {
		return a
	}
	return append(a[:key], a[key+1:]...)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ArrayDiffInt(array1 []int, othersParams ...[]int) []int {
	if len(array1) == 0 {
		return []int{}
	}
	if len(array1) > 0 && len(othersParams) == 0 {
		return array1
	}
	var tmp = make(map[int]int, len(array1))
	for _, v := range array1 {
		tmp[v] = 1
	}
	for _, param := range othersParams {
		for _, arg := range param {
			if tmp[arg] != 0 {
				tmp[arg]++
			}
		}
	}
	var res = make([]int, 0, len(tmp))
	for k, v := range tmp {
		if v == 1 {
			res = append(res, k)
		}
	}
	return res
}

func ArrayDiffUint(array1 []uint, othersParams ...[]uint) []uint {
	if len(array1) == 0 {
		return []uint{}
	}
	if len(array1) > 0 && len(othersParams) == 0 {
		return array1
	}
	var tmp = make(map[uint]int, len(array1))
	for _, v := range array1 {
		tmp[v] = 1
	}
	for _, param := range othersParams {
		for _, arg := range param {
			if tmp[arg] != 0 {
				tmp[arg]++
			}
		}
	}
	var res = make([]uint, 0, len(tmp))
	for k, v := range tmp {
		if v == 1 {
			res = append(res, k)
		}
	}
	return res
}
