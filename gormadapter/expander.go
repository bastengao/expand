package gormadapter

type Expander interface {
	Expand(expand []string) (Scope, error)
}

type expander struct {
	whitelist any
	options   []Option
}

func New(whitelist any, options ...Option) Expander {
	return &expander{
		whitelist: whitelist,
		options:   options,
	}
}

func (e *expander) Expand(expand []string) (Scope, error) {
	return Expand(expand, e.whitelist, e.options...)
}
