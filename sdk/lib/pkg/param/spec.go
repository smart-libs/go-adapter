package sdkparam

type (
	Spec interface {
		Name() string
		Options() []Option
	}

	defaultSpec struct {
		name    string
		options []Option
	}
)

func (d defaultSpec) Name() string      { return d.name }
func (d defaultSpec) Options() []Option { return d.options }

func NewSpec(name string, options ...Option) Spec {
	return defaultSpec{
		name:    name,
		options: options,
	}
}
