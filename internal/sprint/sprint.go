package sprint

// Sprint represents a sprint in the project management context.
type Sprint struct {
	Number int
	Path   string
	Status string
}

func NewSprint(number int, path string, status string) *Sprint {
	return &Sprint{
		Number: number,
		Path:   path,
		Status: status,
	}
}
