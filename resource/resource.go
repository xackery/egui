package resource

// Reader implements
type Reader interface {
	Type() string
	Data() interface{}
}

// Writer implements
type Writer interface {
	SetType() string
	SetData() interface{}
}
