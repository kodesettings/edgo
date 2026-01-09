package operations

type Action string

const (
	Insert Action = "insert"
	AddChar Action = "addchar"
	Delete Action = "delete"
	Start Action = "start"
	End Action = "end"
)

type Operation struct {
	Action Action
	Text   []byte
	Offset int
	Cursor CursorMove
}

type EditOperation []Operation

type CursorMove struct {
	Line int
	Column int
}
