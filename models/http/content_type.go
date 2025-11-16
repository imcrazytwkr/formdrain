package http

type ContentType int8

const (
	ContentTypeUndefined ContentType = iota
	ContentTypeHTML
	ContentTypeJSON
)

var toString = map[ContentType]string{
	ContentTypeUndefined: "undefined",
	ContentTypeHTML:      "text/html",
	ContentTypeJSON:      "application/json",
}

var fromString map[string]ContentType

func init() {
	fromString = make(map[string]ContentType, len(toString))
	for k, v := range toString {
		fromString[v] = k
	}
}

var formFromString = map[string]ContentType{
	"application/x-www-form-urlencoded": ContentTypeHTML,
	"application/json":                  ContentTypeJSON,
}

func (c ContentType) String() string {
	return toString[c]
}

func ParseContentType(str string) ContentType {
	val, ok := fromString[str]
	if !ok {
		return ContentTypeUndefined
	}

	return val
}

func ParseFormContentType(str string) ContentType {
	val, ok := formFromString[str]
	if !ok {
		return ContentTypeUndefined
	}

	return val
}
