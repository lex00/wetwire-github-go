package workflow

// Strategy configures the job execution strategy.
type Strategy struct {
	// Matrix defines the matrix configuration.
	Matrix *Matrix `yaml:"matrix,omitempty"`

	// FailFast stops all jobs if any matrix job fails.
	FailFast *bool `yaml:"fail-fast,omitempty"`

	// MaxParallel limits the number of parallel matrix jobs.
	MaxParallel int `yaml:"max-parallel,omitempty"`
}

// Matrix defines a build matrix for running jobs with different configurations.
type Matrix struct {
	// Values defines the matrix dimensions and their values.
	// Each key is a dimension name, each value is a list of possible values.
	//
	// Example:
	//   Values: map[string][]any{
	//       "go":   {"1.22", "1.23"},
	//       "os":   {"ubuntu-latest", "macos-latest"},
	//   }
	Values map[string][]any `yaml:",inline"`

	// Include adds additional matrix combinations.
	Include []map[string]any `yaml:"include,omitempty"`

	// Exclude removes specific matrix combinations.
	Exclude []map[string]any `yaml:"exclude,omitempty"`
}

// NewMatrix creates a new Matrix with the given dimensions.
func NewMatrix(values map[string][]any) *Matrix {
	return &Matrix{Values: values}
}

// WithInclude adds include combinations to the matrix.
func (m *Matrix) WithInclude(include ...map[string]any) *Matrix {
	m.Include = append(m.Include, include...)
	return m
}

// WithExclude adds exclude combinations to the matrix.
func (m *Matrix) WithExclude(exclude ...map[string]any) *Matrix {
	m.Exclude = append(m.Exclude, exclude...)
	return m
}
