package langs

type Language interface {
	Query() string
}

var languages = map[string]Language{
	"javascript": &Javascript{},
	"typescript": &Javascript{},
	"python":     &Python{},
	"rust":       &Rust{},
	"go":         &Go{},
	"c":          &C{},
	"c++":        &Cpp{},
	"cpp":        &Cpp{},
	"css":        &Css{},
	"html":       &Html{},
	"java":       &Java{},
	"bash":       &Bash{},
}

func MatchQueryLang(lang string) string {
	if l, exists := languages[lang]; exists {
		return l.Query()
	}
	return ""
}
