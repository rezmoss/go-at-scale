// Example 95
//go:generate stringer -type=Status
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusSuspended
    StatusCancelled
)

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
    Find(id string) (*Entity, error)
    Save(entity *Entity) error
}