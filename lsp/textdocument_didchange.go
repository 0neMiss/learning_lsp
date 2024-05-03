package lsp

type TextDocumentDidChangeNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}

type DidChangeTextDocumentParams struct {
	TextDocument VersionTextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}
