package cselect

import "reflect"

type Selector struct {
	cases     []reflect.SelectCase
	selecting []reflect.SelectCase
	predict   []func() bool
	toSend    []*reflect.Value
	cbs       []interface{}
}

func New() *Selector {
	return &Selector{}
}

var emptyCase = reflect.SelectCase{
	Dir: reflect.SelectRecv,
}

func (s *Selector) Add(ch interface{}, cb interface{}, predict func() bool, toSend interface{}) {
	dir := reflect.SelectRecv
	if toSend != nil { // send
		dir = reflect.SelectSend
	}
	selectCase := reflect.SelectCase{
		Dir:  dir,
		Chan: reflect.ValueOf(ch),
	}
	s.cases = append(s.cases, selectCase)
	s.selecting = append(s.selecting, selectCase)
	s.predict = append(s.predict, predict)
	if toSend != nil {
		v := reflect.ValueOf(toSend)
		s.toSend = append(s.toSend, &v)
	} else {
		s.toSend = append(s.toSend, nil)
	}
	s.cbs = append(s.cbs, cb)
}

func (s *Selector) Select() {
	for i := 0; i < len(s.cases); i++ {
		if s.predict[i] != nil {
			if !s.predict[i]() {
				s.selecting[i] = emptyCase
			} else {
				s.selecting[i] = s.cases[i]
			}
		}
		if s.toSend[i] != nil {
			ret := s.toSend[i].Call(nil)
			s.selecting[i].Send = ret[0]
		}
	}
	n, recv, ok := reflect.Select(s.selecting)
	if s.selecting[n].Dir == reflect.SelectRecv {
		if ok {
			if s.cbs[n] != nil {
				switch cb := s.cbs[n].(type) {
				case func():
					cb()
				default:
					reflect.ValueOf(cb).Call([]reflect.Value{
						recv,
					})
				}
			}
		}
	} else {
		if s.cbs[n] != nil {
			(s.cbs[n].(func()))()
		}
	}
}
