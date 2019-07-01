package atomic

type Ivolatile interface {
	Set(value interface{})
	Read() interface{}
}
type Volatile struct {
	value chan interface{}
}

func (v *Volatile) Set(value interface{}) {
	<-v.value
	v.value <- value
}
func (v *Volatile) Read() interface{} {
	t := <-v.value
	v.value <- t
	return t
}
func NewVolatile(init interface{}) Volatile {
	v := Volatile{
		value: make(chan interface{},1),
	}
	v.value <- init
	return v
}
