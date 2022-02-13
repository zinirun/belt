package belt

type I interface{}

type Mapper func(I, int) (I, error)
type Reducer func(I, int) (I, error)
type Filter func(I, int) (bool, error)
type Piper func([]I) error
