package ws

import (
	"context"
	"log"
	"net/http"

	"docker-ui/docker"
	"github.com/docker/docker/api/types"
	"github.com/gorilla/mux"
)

func TerminalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// In v26, ExecConfig and ExecStartCheck are in types
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          true,
		Cmd:          []string{"/bin/sh"},
	}

	execID, err := docker.Cli.ContainerExecCreate(context.Background(), id, execConfig)
	if err != nil {
		log.Printf("Failed to create exec: %v", err)
		return
	}

	attachConfig := types.ExecStartCheck{
		Tty: true,
	}

	resp, err := docker.Cli.ContainerExecAttach(context.Background(), execID.ID, attachConfig)
	if err != nil {
		log.Printf("Failed to attach exec: %v", err)
		return
	}
	defer resp.Close()

	// Handle input from WebSocket to container
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			resp.Conn.Write(msg)
		}
	}()

	// Handle output from container to WebSocket
	buf := make([]byte, 1024)
	for {
		n, err := resp.Reader.Read(buf)
		if n > 0 {
			if err := conn.WriteMessage(1, buf[:n]); err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}
