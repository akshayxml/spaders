package EntityState

type EntityState int

const (
	Alive EntityState = iota
	Dead  EntityState = iota
	Dying EntityState = iota
)
