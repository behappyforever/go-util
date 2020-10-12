package util

import (
	"reflect"
	"testing"
)

func TestHandleByPartition(t *testing.T) {
	expected := []int{0, 3, 3, 6, 6, 9, 9, 10}
	res := make([]int, 0)
	HandleByPartition(10, 3, func(l int, h int) {
		res = append(res, l, h)
	})

	t.Log(res)
	if len(expected) != len(res) {
		t.Fail()
	}
	for i := range expected {
		if expected[i] != res[i] {
			t.Fail()
		}
	}
}

func TestHandleByPartition1(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}
	res := make([][]int, 0)
	HandleByPartition(len(arr), 3, func(l int, h int) {
		res = append(res, arr[l:h])
	})

	t.Log(res)
	if len(expected) != len(res) {
		t.Fail()
	}

	for i := range res {
		if !reflect.DeepEqual(expected[i], res[i]) {
			t.Fail()
		}
	}
}

func TestHandleByPartition2(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	expected := [][]int{{1, 2, 3, 4, 5}}
	res := make([][]int, 0)
	HandleByPartition(len(arr), 10, func(l int, h int) {
		res = append(res, arr[l:h])
	})

	t.Log(res)
	if len(expected) != len(res) {
		t.Fail()
	}

	for i := range res {
		if !reflect.DeepEqual(expected[i], res[i]) {
			t.Fail()
		}
	}
}
