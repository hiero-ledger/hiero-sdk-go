package hiero

// SPDX-License-Identifier: Apache-2.0

type _LockableSlice struct {
	slice  []interface{}
	locked bool
	index  int
}

func _NewLockableSlice() *_LockableSlice {
	return &_LockableSlice{
		slice: []interface{}{},
	}
}

func (ls *_LockableSlice) _RequireNotLocked() {
	if ls.locked {
		panic(errLockedSlice)
	}
}

func (ls *_LockableSlice) _SetLocked(locked bool) *_LockableSlice { // nolint
	ls.locked = locked
	return ls
}

func (ls *_LockableSlice) _SetSlice(slice []interface{}) *_LockableSlice { //nolint
	ls._RequireNotLocked()
	ls.slice = slice
	ls.index = 0
	return ls
}

func (ls *_LockableSlice) _Push(items ...interface{}) *_LockableSlice {
	ls._RequireNotLocked()
	ls.slice = append(ls.slice, items...)
	return ls
}

func (ls *_LockableSlice) _Clear() *_LockableSlice { //nolint
	ls._RequireNotLocked()
	ls.slice = []interface{}{}
	return ls
}

func (ls *_LockableSlice) _Get(index int) interface{} { //nolint
	return ls.slice[index]
}

func (ls *_LockableSlice) _Set(index int, item interface{}) *_LockableSlice { //nolint
	ls._RequireNotLocked()

	if len(ls.slice) == index {
		ls.slice = append(ls.slice, item)
	} else {
		ls.slice[index] = item
	}

	return ls
}

func (ls *_LockableSlice) _SetIfAbsent(index int, item interface{}) *_LockableSlice { //nolint
	ls._RequireNotLocked()
	if len(ls.slice) == index || ls.slice[index] == nil {
		ls._Set(index, item)
	}
	return ls
}

func (ls *_LockableSlice) _GetNext() interface{} { //nolint
	return ls._Get(ls._Advance())
}

func (ls *_LockableSlice) _GetCurrent() interface{} { //nolint
	return ls._Get(ls.index)
}

func (ls *_LockableSlice) _Advance() int { //nolint
	index := ls.index
	if len(ls.slice) != 0 {
		ls.index = (ls.index + 1) % len(ls.slice)
	}
	return index
}

func (ls *_LockableSlice) _IsEmpty() bool { //nolint
	return len(ls.slice) == 0
}

func (ls *_LockableSlice) _Length() int { //nolint
	return len(ls.slice)
}
