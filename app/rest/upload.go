package rest

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"

	ps "github.com/denniswon/tcex/app/pubsub"
)

type wsConn struct {
	conn *websocket.Conn
}

func (ws *wsConn) sendUploadHeader(header *ps.UploadHeader) {
	ws.conn.WriteMessage(websocket.TextMessage, []byte(header.String()))
}

func (ws *wsConn) requestNextBlock() {
	ws.conn.WriteMessage(websocket.TextMessage, []byte("NEXT"))
}

// HandleUpload - Handles file upload from client to server and sends the result response back to client
// returns error if any during upload
func HandleUpload(conn *websocket.Conn, header *ps.UploadHeader, tempDir string) error {

	ws := &wsConn{conn: conn}
	var err error

	_filepath := filepath.Join(tempDir, header.Filepath)

	header.Filepath = _filepath
	header.Generate() // assign a new request id

	// Check if file already exists
	if _, err = os.Stat(_filepath); err == nil {

		log.Printf("Upload file already exists : %s %d\n", header.Filepath, header.Size)

		ws.sendUploadHeader(header)

		return nil

	}

	// Create temp file to save file.
	var namedFile *os.File
	if namedFile, err = os.Create(header.Filepath); err != nil {
		return err
	}
	defer func() {
		namedFile.Close()
	}()

	// Read file blocks until all bytes are received.
	var bytesRead int64 = 0

	for {

		mt, message, _err := ws.conn.ReadMessage()

		if _err != nil {
			err = _err
			break
		}

		if mt != websocket.BinaryMessage {

			if mt == websocket.TextMessage && string(message) == "CANCEL" {
				err = errors.New("upload cancelled")
				break
			}

			err = errors.New("invalid file block received")
			break
		}

		namedFile.Write(message)

		bytesRead += int64(len(message))

		if bytesRead == header.Size {

			namedFile.Close()
			break
		}

		ws.requestNextBlock()
	}

	if err != nil {
		os.Remove(namedFile.Name())
		return err
	} else {

		log.Printf("Upload finished with request id %s: %s %d\n", header.ID, header.Filepath, header.Size)

		ws.sendUploadHeader(header)
	}

	return nil
}
