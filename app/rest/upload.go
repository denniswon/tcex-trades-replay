package rest

import (
	"errors"
	"os"

	"github.com/gorilla/websocket"
)

type UploadHeader struct {
	Filename  string
	Size      int
	RequestID string
}

type wsConn struct {
	conn *websocket.Conn
}

func (ws *wsConn) requestNextBlock() {
	ws.conn.WriteMessage(websocket.TextMessage, []byte("NEXT"))
}

// HandleUpload - Handles file upload from client to server
// and returns named file name and bytes receive created in the temp directory
func HandleUpload(conn *websocket.Conn, header *UploadHeader) (string, int, error) {

	ws := &wsConn{conn: conn}
	var err error

	// Create temp file to save file.
	var namedFile *os.File
	if namedFile, err = os.Create(header.Filename); err != nil {
		return "", 0, err
	}
	defer func() {
		namedFile.Close()
	}()

	// Read file blocks until all bytes are received.
	var bytesRead int = 0

	for {

		mt, message, _err := ws.conn.ReadMessage()

		if _err != nil {
			err = _err
			break
		}

		if mt != websocket.BinaryMessage {

			if mt == websocket.TextMessage {

				if string(message) == "CANCEL" {
					err = errors.New("upload cancelled")
					break

				}

			}

			err = errors.New("invalid file block received")
			break
		}

		namedFile.Write(message)

		bytesRead += len(message)

		if bytesRead == header.Size {
			namedFile.Close()
			break
		}

		ws.requestNextBlock()
	}

	if err != nil {
		os.Remove(namedFile.Name())
		return "", 0, err
	}

	return namedFile.Name(), bytesRead, nil
}
