package lsp

func (this *LspClient) Hover(file string, line int, character int) (HoverResponse, error) {
	this.id++

	request := BaseRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/hover",
		Params: Params {
			TextDocument: TextDocument { URI: "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.send(request)
	response, err := WaitForRequest[HoverResponse](this.userMessages, 1000)

	return response, err
}

func (this *LspClient) Completion(file string, line int, character int) (CompletionResponse, error) {
	this.id++

	request := BaseRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/completion",
		Params: Params{
			TextDocument: TextDocument { URI:  "file://" + file },
			Position: Position { Line: line, Character: character },
			Context: Context { TriggerKind: 1 },
		},
	}

	this.send(request)
	response, err := WaitForRequest[CompletionResponse](this.userMessages, 1000)

	return response, err
}

func (this *LspClient) Definition(file string, line int, character int) (DefinitionResponse, error) {
	this.id++

	request := DefinitionRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/definition",
		Params: DefinitionParams {
			TextDocument: TextDocument{ URI: "file://" + file },
			Position: Position{ Line: line, Character: character },
		},
	}

	this.send(request)
	response, err := WaitForRequest[DefinitionResponse](this.userMessages, 1000)

	return response, err
}

func (this *LspClient) SignatureHelp(file string, line int, character int) (SignatureHelpResponse, error) {
	this.id++

	request := BaseRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/signatureHelp",
		Params: Params {
			TextDocument: TextDocument { URI: "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.send(request)
	response, err := WaitForRequest[SignatureHelpResponse](this.userMessages, 1000)

	return response, err
}

func (this *LspClient) References(file string, line int, character int) (ReferencesResponse, error) {
	this.id++

	request := BaseRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/references",
		Params: Params{
			TextDocument: TextDocument{ URI: "file://" + file },
			Position: Position{ Line: line, Character: character },
			Context: Context{ IncludeDeclaration: false },
		},
	}

	this.send(request)
	response, err := WaitForRequest[ReferencesResponse](this.userMessages, 3000)

	return response, err
}


func (this *LspClient) PrepareRename(file string, line int, character int) (PrepareRenameResponse, error) {
	this.id++

	request := PrepareRenameRequest {
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/prepareRename",
		Params: Params{
			TextDocument: TextDocument { URI:  "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.send(request)
	response, err := WaitForRequest[PrepareRenameResponse](this.userMessages, 10000)

	return response, err
}

func (this *LspClient) Rename(file string, newname string, line int, character int) (RenameResponse, error) {
	this.id++

	request := RenameRequest{
		ID: this.id,  JSONRPC: "2.0", Method:  "textDocument/rename",
		Params: RenameParams {
			NewName: newname,
			Position: Position { Line: line, Character: character },
			TextDocument: TextDocument { URI:  "file://" + file },
		},
	}

	this.send(request)
	response, err := WaitForRequest[RenameResponse](this.userMessages, 10000)

	return response, err
}

func (this *LspClient) CodeAction(file string, spc int, spl int, epc int, epl int) (CodeActionResponse, error) {
	this.id++

	request := CodeActionRequest {
		ID: this.id,  JSONRPC: "2.0", Method: "textDocument/codeAction",
		Params: CodeActionParams {
			TextDocument: TextDocument { URI:  "file://" + file },
			Context: Context{ Only: []string{"refactor"}, TriggerKind: 1 },
			Range: RequestRange{
				Start: Position{ Line: spl, Character: spc},
				End: Position{ Line: epl, Character: epc},
			},
		},
	}

	this.send(request)
	response, err := WaitForRequest[CodeActionResponse](this.userMessages, 10000)

	return response, err
}
