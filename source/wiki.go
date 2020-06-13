package source

type DocHandler func(title, body string)

func LoadWiki(wikipedia string, maxIndexCount int, handler DocHandler) error{
	return nil
}
