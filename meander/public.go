package meander

// Facade exposes a Public method, which returns the
// public view of a struct.
type Facade interface {
	Public() interface{}
}

// Public takes any object and checks whether it implements
// the Facade interface. If it does, it calls the Public method
// and returns the result. Otherwise it just returns
// the original object.
func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}
