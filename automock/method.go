package automock

import "go/types"

// Method contains the details from an interface method
type Method struct {
	gen *generator
	fn  *types.Func
}

// Name returns the method name
func (m Method) Name() string {
	return m.fn.Name()
}

// ParamTypes returns the list of types for the params
func (m Method) ParamTypes() []string {
	sig := m.signature()
	return m.listTypes(sig.Params())
}

// ReturnTypes returns the list of types for the params
func (m Method) ReturnTypes() []string {
	sig := m.signature()
	return m.listTypes(sig.Results())
}

func (m Method) listTypes(t *types.Tuple) []string {
	num := t.Len()
	list := make([]string, num)
	for i := 0; i < num; i++ {
		list[i] = types.TypeString(t.At(i).Type(), m.gen.qf)
	}
	return list
}

func (m Method) signature() *types.Signature {
	return m.fn.Type().(*types.Signature)
}
