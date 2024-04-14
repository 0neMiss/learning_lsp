package main

import (
	"bufio"
	"build_lsp/analysis"
	"build_lsp/lsp"
	"build_lsp/rpc"
	"encoding/json"
	"log"
	"os"
)

func main() {

	logger := getLogger("/home/jordan/repos/build_lsp/log.txt")
	logger.Println("Main.go has started!")
	// attach a scanner to stdin to break up the messages from the editor
	scanner := bufio.NewScanner(os.Stdin)
	// Attach the function to stdin that will be used to split the header, contentLenth, and content
	scanner.Split(rpc.Split)
	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		// grab the []byte the scanner is currently reading and attempt to decode
		// we pull out the method of the message, and the paylaod(contents)
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		// Pass the decoded message to our handler
		handleMessage(logger, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, state analysis.State, method string, content []byte) {
	logger.Printf("Received msg with method: %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}
		logger.Printf("Connected to %s on version %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)

		// Try and reply to the initialize method. Communcation happens through stdout
		msg := lsp.NewInitializeResponse(request.ID)
		// First we need to put the message in the proper format per the spec
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		writer.Write([]byte(reply))
		logger.Print("Sent the reply")

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}
		logger.Printf("Opened: %s %s ",
			request.Params.TextDocument.URI,
			request.Params.TextDocument.Text,
		)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	}
}

// We cant log to stdout because thats how we communicate with the editor, for now just logging to a file.
func getLogger(filePath string) *log.Logger {
	// create a file, at the path provided, truncate it, and make it read write, 0666 is who can do it and it means pretty much anybody
	logfile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return log.New(logfile, "[mylsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
