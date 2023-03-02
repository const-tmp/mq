package queue

import "testing"

func TestList(t *testing.T) {
	l := new(List[int])
	for i := 0; i < num; i++ {
		l.Push(i)
	}
	for i := 0; i < num; i++ {
		if item := l.Pop(); item == nil || item.value != i {
			t.Error()
		}
	}
}

func TestList_Remove(t *testing.T) {
	for j := 0; j < num; j++ {
		l := new(List[int])
		var pushed, popped []*ListItem[int]

		for i := 0; i < num; i++ {
			pushed = append(pushed, l.Push(i))
		}

		l.Remove(pushed[j])

		for i := 0; i < num; i++ {
			if v := l.Pop(); v != nil {
				popped = append(popped, v)
			}
		}
		if contains(popped, pushed[j]) {
			t.Errorf("j=%d popped contains %d", j, pushed[j].value)
		}
	}
}

func contains(slice []*ListItem[int], item *ListItem[int]) bool {
	for _, li := range slice {
		if li == item {
			return true
		}
	}
	return false
}
