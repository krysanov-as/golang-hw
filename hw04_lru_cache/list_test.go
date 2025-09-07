package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func Test_list_PushFront(t *testing.T) {
	t.Run("push to empty list", func(t *testing.T) {
		var l list
		got := l.PushFront(100)

		if got.Value != 100 {
			t.Errorf("Value = %v, want %v", got.Value, 100)
		}

		if got.Next != nil {
			t.Errorf("Next = %v, want nil", got.Next)
		}
		if got.Prev != nil {
			t.Errorf("Prev = %v, want nil", got.Prev)
		}

		if l.head != got {
			t.Errorf("head = %v, want %v", l.head, got)
		}

		if l.tail != got {
			t.Errorf("tail = %v, want %v", l.tail, got)
		}
	})

	t.Run("push to non-empty list", func(t *testing.T) {
		var l list
		first := l.PushFront(50)
		got := l.PushFront(100)

		if got.Value != 100 {
			t.Errorf("Value = %v, want %v", got.Value, 100)
		}
		if first.Prev != got {
			t.Errorf("first.Prev = %v, want %v", first.Prev, got)
		}
		if got.Next != first {
			t.Errorf("got.Next = %v, want %v", got.Next, first)
		}

		if l.head != got {
			t.Errorf("head = %v, want %v", l.head, got)
		}
	})
}
