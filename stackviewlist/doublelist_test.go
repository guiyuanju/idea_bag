package stackviewlist

import (
	"testing"
)

func TestDoubleListAppend(t *testing.T) {
	d := newEmptyDoubleList[int]()
	if d.length() != 0 {
		t.Error("empty length should be zero: ", d.length())
	}

	d = newDoubleList([]int{1, 2, 3})
	if d.length() != 3 {
		t.Error("empty length should be three: ", d.length())
	}

	d.clear()
	d.append(1).append(2).append(3)
	if d.length() != 3 {
		t.Error("empty length should be three: ", d.length())
	}

	d.prepend(0)
	if d.length() != 4 {
		t.Error("empty length should be four: ", d.length())
	}
	if d.head.value != 0 {
		t.Error("head should be 0: ", d.head.value)
	}
}

func TestDoubleListRemove(t *testing.T) {
	d := newEmptyDoubleList[int]()
	d.remove(nil)

	d.append(1)
	d.remove(d.head)
	if d.length() != 0 {
		t.Error("expect zero length after removing: ", d.length())
	}

	d.append(1).append(2)
	d.remove(d.end)
	if d.length() != 1 || d.head.value != 1 || d.head != d.end {
		t.Error("expect 1 left after removing: ", d.String())
	}

	d.clear()
	if d.length() != 0 {
		t.Error("expect zero length after clearing: ", d.String())
	}

	d.clear().append(1).append(2).append(3)
	d.remove(d.head.next)
	if d.head.value != 1 || d.end.value != 3 || d.length() != 2 {
		t.Error("expect 1, 3 left after removing 2: ", d.String())
	}
}
