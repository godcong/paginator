package paginator

type Conditioner interface {
	Conditions() []string
}
