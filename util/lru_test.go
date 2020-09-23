package util

import (
	"container/list"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func TestLru(t *testing.T) {

	l := NewLruCache(8, false)

	for i := 0; i < 10; i++ {
		l.Put(Key(strconv.Itoa(i)), i+100)
	}
	lruPrinter(t, l.(*lru))
	checkList(t, l.(*lru).List, []interface{}{
		node{k: Key(strconv.Itoa(9)), v: 109},
		node{k: Key(strconv.Itoa(8)), v: 108},
		node{k: Key(strconv.Itoa(7)), v: 107},
		node{k: Key(strconv.Itoa(6)), v: 106},
		node{k: Key(strconv.Itoa(5)), v: 105},
		node{k: Key(strconv.Itoa(4)), v: 104},
		node{k: Key(strconv.Itoa(3)), v: 103},
		node{k: Key(strconv.Itoa(2)), v: 102},
	})

	l.Get(Key(strconv.Itoa(6)))
	lruPrinter(t, l.(*lru))
	checkList(t, l.(*lru).List, []interface{}{
		node{k: Key(strconv.Itoa(6)), v: 106},
		node{k: Key(strconv.Itoa(9)), v: 109},
		node{k: Key(strconv.Itoa(8)), v: 108},
		node{k: Key(strconv.Itoa(7)), v: 107},
		node{k: Key(strconv.Itoa(5)), v: 105},
		node{k: Key(strconv.Itoa(4)), v: 104},
		node{k: Key(strconv.Itoa(3)), v: 103},
		node{k: Key(strconv.Itoa(2)), v: 102},
	})

	l.Get(Key(strconv.Itoa(999)))
	lruPrinter(t, l.(*lru))
	checkList(t, l.(*lru).List, []interface{}{
		node{k: Key(strconv.Itoa(6)), v: 106},
		node{k: Key(strconv.Itoa(9)), v: 109},
		node{k: Key(strconv.Itoa(8)), v: 108},
		node{k: Key(strconv.Itoa(7)), v: 107},
		node{k: Key(strconv.Itoa(5)), v: 105},
		node{k: Key(strconv.Itoa(4)), v: 104},
		node{k: Key(strconv.Itoa(3)), v: 103},
		node{k: Key(strconv.Itoa(2)), v: 102},
	})

	l.Put(Key(strconv.Itoa(11)), 111)
	lruPrinter(t, l.(*lru))
	checkList(t, l.(*lru).List, []interface{}{
		node{k: Key(strconv.Itoa(11)), v: 111},
		node{k: Key(strconv.Itoa(6)), v: 106},
		node{k: Key(strconv.Itoa(9)), v: 109},
		node{k: Key(strconv.Itoa(8)), v: 108},
		node{k: Key(strconv.Itoa(7)), v: 107},
		node{k: Key(strconv.Itoa(5)), v: 105},
		node{k: Key(strconv.Itoa(4)), v: 104},
		node{k: Key(strconv.Itoa(3)), v: 103},
	})

	l.Put(Key(strconv.Itoa(8)), 1008)
	lruPrinter(t, l.(*lru))
	checkList(t, l.(*lru).List, []interface{}{
		node{k: Key(strconv.Itoa(8)), v: 1008},
		node{k: Key(strconv.Itoa(11)), v: 111},
		node{k: Key(strconv.Itoa(6)), v: 106},
		node{k: Key(strconv.Itoa(9)), v: 109},
		node{k: Key(strconv.Itoa(7)), v: 107},
		node{k: Key(strconv.Itoa(5)), v: 105},
		node{k: Key(strconv.Itoa(4)), v: 104},
		node{k: Key(strconv.Itoa(3)), v: 103},
	})

	l.Rem(Key(strconv.Itoa(8)))
	l.Rem(Key(strconv.Itoa(7)))
	lruPrinter(t, l.(*lru))
	checkList(t, l.(*lru).List, []interface{}{
		node{k: Key(strconv.Itoa(11)), v: 111},
		node{k: Key(strconv.Itoa(6)), v: 106},
		node{k: Key(strconv.Itoa(9)), v: 109},
		node{k: Key(strconv.Itoa(5)), v: 105},
		node{k: Key(strconv.Itoa(4)), v: 104},
		node{k: Key(strconv.Itoa(3)), v: 103},
	})

	l.Put(Key(strconv.Itoa(12)), 112)
	lruPrinter(t, l.(*lru))
	checkList(t, l.(*lru).List, []interface{}{
		node{k: Key(strconv.Itoa(12)), v: 112},
		node{k: Key(strconv.Itoa(11)), v: 111},
		node{k: Key(strconv.Itoa(6)), v: 106},
		node{k: Key(strconv.Itoa(9)), v: 109},
		node{k: Key(strconv.Itoa(5)), v: 105},
		node{k: Key(strconv.Itoa(4)), v: 104},
		node{k: Key(strconv.Itoa(3)), v: 103},
	})
}

func TestLru2(t *testing.T) {
	l := NewLruCache(1e2, true)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for t := 1; t < 1e6; t++ {
				n := rand.Intn(1e6)
				switch n % 4 {
				case 0, 1:
					l.Put(Key(strconv.Itoa(n)), i)
				case 2:
					l.Get(Key(strconv.Itoa(n)))
				case 3:
					l.Rem(Key(strconv.Itoa(n)))
				}
			}
		}(i)
	}
	wg.Wait()
	lruPrinter(t, l.(*lru))
}

func lruPrinter(t *testing.T, l *lru) {
	t.Log("map--------------------")
	for k, v := range l.mp {
		t.Logf("k:%v, v:%v", k, v.Value)
	}

	t.Log("list-------------------")
	for e := l.Front(); e != nil; e = e.Next() {
		t.Logf("%+v", e)
	}
}

func checkList(t *testing.T, l *list.List, es []interface{}) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		le := e.Value.(node)
		if le != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, es[i])
		}
		i++
	}
}

func checkListLen(t *testing.T, l *list.List, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}
