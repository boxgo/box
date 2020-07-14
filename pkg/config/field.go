package config

type Field struct {
	name string
	path string
	desc string
	def  interface{}
}

func NewField(name, path, desc string, def interface{}) *Field {
	return &Field{
		name: name,
		path: path,
		desc: desc,
		def:  def,
	}
}

func (f Field) String() string {
	if f.name != "" {
		return f.name + "." + f.path
	}

	return f.path
}
