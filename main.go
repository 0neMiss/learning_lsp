package main

import (
	"bufio"
	"build_lsp/lsp"
	"build_lsp/rpc"
	"encoding/json"
	"log"
	"os"
)

func main() {
	logger := getLogger("/home/jordan/repos/build_lsp/log.txt")
	logger.Println("Hey, I started!")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		handleMessage(logger, method, contents)
	}
}

func handleMessage(logger *log.Logger, method string, content []byte) {
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
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		writer.Write([]byte(reply))
		logger.Print("Sent the reply")
	}
}

func getLogger(filePath string) *log.Logger {
	// create a file, at the path provided, truncate it, and make it read write, 0666 is who can do it and it means pretty much anybody
	logfile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return log.New(logfile, "[mylsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
