package paginator

type Founder interface {
	Iterator
}

type found struct {
	*iterator
}

var _ Founder = (*found)(nil)
