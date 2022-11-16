package copier

type Copier[S, D any] interface {
	Copy(src *S) (*D, error)
	CopyTo(src *S, dst *D) error
}
