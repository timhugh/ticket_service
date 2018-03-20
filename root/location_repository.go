package root

type LocationRepository interface {
	Create(Location) error
	Find(string) (*Location, error)
}
