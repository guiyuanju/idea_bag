package stackviewlist

type ViewTransformer[T any] func(doubleList[Item[T]]) []*Item[T]

type StackViewList[T any] struct {
	ViewStackBottom doubleList[Item[T]]
	ViewStack       []doubleList[Item[T]]
	viewTransformer []ViewTransformer[T]
}

type Item[T any] struct {
	value         T
	prevStackItem *Item[T]
}

func New[T any](items []T) StackViewList[T] {
	its := []Item[T]{}
	for _, item := range items {
		its = append(its, Item[T]{value: item, prevStackItem: nil})
	}
	d := newDoubleList(its)
	return StackViewList[T]{ViewStackBottom: d}
}

func (m *StackViewList[T]) TopView() *doubleList[Item[T]] {
	if len(m.ViewStack) == 0 {
		return nil
	}
	return &m.ViewStack[len(m.ViewStack)-1]
}

func (m *StackViewList[T]) PushViewStack(view ViewTransformer[T]) {
	m.viewTransformer = append(m.viewTransformer, view)
}

func (m *StackViewList[T]) cascadeTransfrom() {
	if len(m.viewTransformer) == 0 {
		return
	}

	m.ViewStack = make([]doubleList[Item[T]], len(m.viewTransformer))

	for i, tf := range m.viewTransformer {
		var from doubleList[Item[T]]
		if i == 0 {
			from = m.ViewStackBottom
		} else {
			from = m.ViewStack[i-1]
		}
		to := tf(from)
		newList := newEmptyDoubleList[Item[T]]()
		for _, n := range to {
			newItem := Item[T]{value: n.value, prevStackItem: n}
			newList.append(newItem)
		}
		m.ViewStack[i] = newList
	}
}
