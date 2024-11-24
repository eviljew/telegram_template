package tg_template

type Data struct {
	Pattern     string
	Replacement any
}

func NewData(pattern string, replacement any) *Data {
	return &Data{
		Pattern:     pattern,
		Replacement: replacement,
	}
}
