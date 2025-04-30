package langs

type Language interface {
	Query() string
}

var languages = map[string]Language{
	"javascript": &Javascript{},
	"typescript": &Typescript{},
	"python":     &Python{},
	"rust":       &Rust{},
	"go":         &Go{},
	"c":          &C{},
	"c++":        &Cpp{},
	"cpp":        &Cpp{},
	"html":       &Html{},
	"kotlin":     &Kotlin{},
	"yaml":       &Yaml{},
	"java":       &Java{},
	"bash":       &Bash{},
	"toml":       &Toml{},
	"lua":        &Lua{},
}

func MatchQueryLang(lang string) string {
	if l, exists := languages[lang]; exists {
		return l.Query()
	}
	return languages["toml"].Query()
}
