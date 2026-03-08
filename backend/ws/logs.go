package ws

import (
	"context"
	"io"
	"log"
	"net/http"

	"docker-ui/docker"
	"github.com/docker/docker/api/types"
	"github.com/gorilla/mux"
)

func LogsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	}

	out, err := docker.Cli.ContainerLogs(context.Background(), id, options)
	if err != nil {
		log.Printf("Failed to get container logs: %v", err)
		return
	}
	defer out.Close()

	buf := make([]byte, 1024)
	for {
		n, err := out.Read(buf)
		if n > 0 {
			if err := conn.WriteMessage(1, buf[:n]); err != nil {
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}
}
