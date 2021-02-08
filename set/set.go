package set

type Set map[interface{}]struct{}

func New(elements ...interface{}) Set {
	set := make(Set)
	set.Insert(elements...)
	return set
}

func (s Set) Insert(elements ...interface{}) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

func (s Set) HasElement(element interface{}) bool {
	_, hasElement := s[element]
	return hasElement
}

func (s Set) Remove(elements ...interface{}) {
	for _, element := range elements {
		delete(s, element)
	}
}

func (s Set) Each(f func(interface{})) {
	for _, element := range s {
		f(element)
	}
}
