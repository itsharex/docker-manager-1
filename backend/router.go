package main

import (
	"encoding/json"
	"net/http"
	"docker-ui/docker"
	"docker-ui/ws"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Container routes
	r.HandleFunc("/api/containers", ListContainersHandler).Methods("GET")
	r.HandleFunc("/api/containers/{id}/start", StartContainerHandler).Methods("POST")
	r.HandleFunc("/api/containers/{id}/stop", StopContainerHandler).Methods("POST")
	r.HandleFunc("/api/containers/{id}/remove", RemoveContainerHandler).Methods("DELETE")
	r.HandleFunc("/api/containers/{id}/inspect", InspectContainerHandler).Methods("GET")

	// Image routes
	r.HandleFunc("/api/images", ListImagesHandler).Methods("GET")
	r.HandleFunc("/api/images/{id}", RemoveImageHandler).Methods("DELETE")

	// Volume routes
	r.HandleFunc("/api/volumes", ListVolumesHandler).Methods("GET")
	r.HandleFunc("/api/volumes/{id}", RemoveVolumeHandler).Methods("DELETE")

	// Network routes
	r.HandleFunc("/api/networks", ListNetworksHandler).Methods("GET")
	r.HandleFunc("/api/networks/{id}", RemoveNetworkHandler).Methods("DELETE")

	// Stats routes
	r.HandleFunc("/api/info", SystemInfoHandler).Methods("GET")

	// WebSocket routes
	r.HandleFunc("/ws/logs/{id}", ws.LogsHandler)
	r.HandleFunc("/ws/terminal/{id}", ws.TerminalHandler)

	return r
}

func ListContainersHandler(w http.ResponseWriter, r *http.Request) {
	containers, err := docker.ListContainers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(containers)
}

func StartContainerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := docker.StartContainer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func StopContainerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := docker.StopContainer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RemoveContainerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := docker.RemoveContainer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func InspectContainerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	info, err := docker.InspectContainer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func ListImagesHandler(w http.ResponseWriter, r *http.Request) {
	images, err := docker.ListImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(images)
}

func RemoveImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := docker.RemoveImage(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ListVolumesHandler(w http.ResponseWriter, r *http.Request) {
	volumes, err := docker.ListVolumes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(volumes)
}

func RemoveVolumeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := docker.RemoveVolume(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ListNetworksHandler(w http.ResponseWriter, r *http.Request) {
	networks, err := docker.ListNetworks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(networks)
}

func RemoveNetworkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := docker.RemoveNetwork(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SystemInfoHandler(w http.ResponseWriter, r *http.Request) {
	info, err := docker.GetSystemInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(info)
}
