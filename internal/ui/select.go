package ui

import (
	highlighter "github.com/vipmax/edgo/internal/highlighter"
)

func (e *Editor) OnSelectMoreAtCursor() {
	var node highlighter.NodeRange

	if !e.Selection.IsSelected || e.TreePath == nil || (e.TreePath.Aty != e.Row || e.TreePath.Atx != e.Col) {
		treepath := e.treeSitterHighlighter.GetNodePathAt(e.Row, e.Col, e.Row, e.Col)
		e.TreePath = &treepath
		node = e.TreePath.CurrentNode()
	} else {
		node = e.TreePath.Next()
	}

	e.Selection.Ssx = node.Ssx; e.Selection.Ssy = node.Ssy
	e.Selection.Sex = node.Sex; e.Selection.Sey = node.Sey
	e.Selection.IsSelected = true
}

func (e *Editor) OnSelectLessAtCursor() {
	if e.TreePath == nil { return }
	node := e.TreePath.Prev()
	e.Selection.Ssx = node.Ssx; e.Selection.Ssy = node.Ssy
	e.Selection.Sex = node.Sex; e.Selection.Sey = node.Sey
	e.Selection.IsSelected = true
}

func (e *Editor) OnSelectAll() {
	if len(e.Lines) == 0 { return }
	e.Selection.Ssx = 0; e.Selection.Ssy = 0
	e.Selection.Sey = len(e.Lines)
	lastElement := len(e.Lines[len(e.Lines)-1].Buf)
	e.Selection.Sex = lastElement
	e.Selection.Sey = len(e.Lines)
	e.Selection.IsSelected = true
}

