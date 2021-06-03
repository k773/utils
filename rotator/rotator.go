package rotator

import "sync"

type Rotator struct {
	Items struct {
		M map[int]interface{}
		i int
		s sync.RWMutex
	}
	itemsUnused struct {
		M map[int]interface{}
		s sync.RWMutex
	}
	autoRefillUnused     bool
	customRefillFuncMeta string                  // customRefillFuncMeta wll be passed to customRefillFunc
	customRefillFunc     func(meta string) error // customRefillFunc will be called instead of simple refilling
	s                    sync.RWMutex
}

//
func NewRotator(autoRefillUnused bool, customRefillFunc func(meta string) error, customRefillFuncMeta string) *Rotator {
	return &Rotator{
		Items: struct {
			M map[int]interface{}
			i int
			s sync.RWMutex
		}{M: map[int]interface{}{}},
		itemsUnused: struct {
			M map[int]interface{}
			s sync.RWMutex
		}{M: map[int]interface{}{}},
		autoRefillUnused:     autoRefillUnused,
		s:                    sync.RWMutex{},
		customRefillFunc:     customRefillFunc,
		customRefillFuncMeta: customRefillFuncMeta,
	}
}

func (r *Rotator) AddItem(pushToUnused bool, item interface{}) int {
	r.Items.s.Lock()
	defer r.Items.s.Unlock()

	r.Items.M[r.Items.i] = item
	if pushToUnused {
		r.itemsUnused.s.Lock()
		defer r.itemsUnused.s.Unlock()
		r.itemsUnused.M[r.Items.i] = item
	}

	r.Items.i++
	return r.Items.i - 1
}

func (r *Rotator) AddItems(alreadyLocked, pushToUnused bool, items ...interface{}) (ids []int) {
	if !alreadyLocked {
		r.Items.s.Lock()
		r.itemsUnused.s.Lock()
		defer r.Items.s.Unlock()
		defer r.itemsUnused.s.Unlock()
	}

	for _, item := range items {
		ids = append(ids, r.Items.i)
		r.Items.M[r.Items.i] = item
		if pushToUnused {
			r.itemsUnused.M[r.Items.i] = item
		}

		r.Items.i++
	}

	return
}

func (r *Rotator) RemoveItem(id int, removeFromUnused bool) {
	r.Items.s.Lock()
	defer r.Items.s.Unlock()

	delete(r.Items.M, id)
	if removeFromUnused {
		r.itemsUnused.s.Lock()
		defer r.itemsUnused.s.Unlock()
		delete(r.itemsUnused.M, id)
	}
}

func (r *Rotator) GetItemByID(id int) (item interface{}, unused bool) {
	r.Items.s.Lock()
	r.itemsUnused.s.Lock()
	defer r.itemsUnused.s.Unlock()
	defer r.Items.s.Unlock()

	_, unused = r.itemsUnused.M[id]
	item = r.Items.M[id]
	return
}

func (r *Rotator) GetRandomUnusedItem() interface{} {
	r.itemsUnused.s.RLock()
	if len(r.itemsUnused.M) == 0 {
		r.itemsUnused.s.RUnlock()
		if r.autoRefillUnused {
			r.RefillUnused()
		}
		r.itemsUnused.s.RLock()
	}
	var item interface{}
	for _, a := range r.itemsUnused.M {
		item = a
		break
	}
	r.itemsUnused.s.RUnlock()
	return item
}

func (r *Rotator) PullRandomUnusedItem() interface{} {
	r.itemsUnused.s.RLock()
	if len(r.itemsUnused.M) == 0 {
		r.itemsUnused.s.RUnlock()
		if r.autoRefillUnused {
			r.RefillUnused()
		}
		r.itemsUnused.s.RLock()
	}
	var item interface{}
	for i, a := range r.itemsUnused.M {
		delete(r.itemsUnused.M, i)
		item = a
		break
	}
	r.itemsUnused.s.RUnlock()
	return item
}

func (r *Rotator) RefillUnused() {
	if r.customRefillFunc != nil {
		_ = r.customRefillFunc(r.customRefillFuncMeta)
		return
	} else {
		r.Items.s.Lock()
		r.itemsUnused.s.Lock()
		defer r.itemsUnused.s.Unlock()
		defer r.Items.s.Unlock()

		r.itemsUnused.M = map[int]interface{}{}
		for i, item := range r.Items.M {
			r.itemsUnused.M[i] = item
		}
	}
}

// It items.len > 0 then AddItems() will be called
func (r *Rotator) Clear(pushToUnused bool, items ...interface{}) {
	r.Items.s.Lock()
	r.itemsUnused.s.Lock()
	defer r.itemsUnused.s.Unlock()
	defer r.Items.s.Unlock()

	r.Items.i = 0

	r.Items.M = map[int]interface{}{}
	r.itemsUnused.M = map[int]interface{}{}

	if len(items) > 0 {
		r.AddItems(true, pushToUnused, items...)
	}
}

func (r *Rotator) EnableAutoRefillUnused() {
	r.s.Lock()
	defer r.s.Unlock()
	r.autoRefillUnused = true
}

func (r *Rotator) DisableAutoRefillUnused() {
	r.s.Lock()
	defer r.s.Unlock()
	r.autoRefillUnused = false
}

func (r *Rotator) CountUnused() (i int) {
	r.itemsUnused.s.RLock()
	r.itemsUnused.s.RUnlock()
	return len(r.itemsUnused.M)
}
