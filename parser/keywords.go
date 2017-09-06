package parser

// keyword lookup
//type int int

const (
	kwInvalid int = iota
	kwMerge
	kwStyle
	kwSet
	kwUse
)

var keywords = map[string]int{
	"merge": kwMerge,
	"style": kwStyle,
	"set":   kwSet,
	"use":   kwUse,
}

func lookupKeyword(name string) (kw int, isKeyWord bool) {
	kw, isKeyWord = keywords[name]
	return
}
