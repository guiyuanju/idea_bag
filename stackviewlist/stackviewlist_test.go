package stackviewlist

import (
	"testing"
)

func TestNewStackViewList(t *testing.T) {
	items := []int{1, 2, 3, 4}
	s := New(items)
	if len(s.ViewStack) != 0 {
		t.Errorf("view stack not zero")
	}

	if s.ViewStackBottom.head.value.value != 1 {
		t.Errorf("base head value incorrect")
	}
}

func TestTransformStackViewList(t *testing.T) {
	items := []int{1, 2, 3, 4}
	s := New(items)
	transformer := func(d doubleList[Item[int]]) []*Item[int] {
		var res []*Item[int]
		for item := range d.All() {
			if item.value >= 2 && item.value < 4 {
				res = append(res, item)
			}
		}
		return res
	}
	s.PushViewStack(transformer)
	if len(s.viewTransformer) != 1 {
		t.Error("expect viewTransformer has lenght one after one push: ", len(s.viewTransformer))
	}
	s.cascadeTransfrom()
	if len(s.ViewStack) != 1 {
		t.Error("expect view stack has lenght one transform: ", len(s.ViewStack))
	}
	if s.ViewStack[0].length() != 2 || s.ViewStack[0].head.value.value != 2 || s.ViewStack[0].end.value.value != 3 {
		t.Error("wrong transform result, expect 2 and 3, got: ", s.ViewStack[0].head.value.value, s.ViewStack[0].end.value.value)
	}
	if s.ViewStack[0].head.value.prevStackItem.value != s.ViewStack[0].head.value.value {
		t.Error("wrong prev stack item conncection, expect 2 got: ", s.ViewStack[0].head.value.prevStackItem.value)
	}
}
