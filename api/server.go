package api

import (
	"github.com/assetto-corsa-web/accweb/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	maxMemory = 10000000 // 10 MB
)

func SaveServerSetttingsHandler(w http.ResponseWriter, r *http.Request) {
	req := &server.ServerSettings{}

	if err := decodeJSON(r, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	if err := server.SaveServerSettings(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	writeResponse(w, nil)
}

func CopyServerSetttingsHandler(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Id int `json:"id"`
	}{}

	if err := decodeJSON(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	if err := server.CopyServerSettings(req.Id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	writeResponse(w, nil)
}

func GetServerHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeResponse(w, server.GetServerList())
	} else {
		idInt, err := strconv.Atoi(id)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		writeResponse(w, server.GetServerById(idInt))
	}
}

func DeleteServerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	if err := server.DeleteServer(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	writeResponse(w, nil)
}

func ImportServerHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	configuration, configurationHeader, err := r.FormFile("configuration")

	if err != nil || configurationHeader.Size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	defer func() {
		if err := configuration.Close(); err != nil {
			logrus.WithError(err).Error("Error closing file on import")
		}
	}()
	settings, settingsHeader, err := r.FormFile("settings")

	if err != nil || settingsHeader.Size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	defer func() {
		if err := settings.Close(); err != nil {
			logrus.WithError(err).Error("Error closing file on import")
		}
	}()
	event, eventHeader, err := r.FormFile("event")

	if err != nil || eventHeader.Size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	defer func() {
		if err := event.Close(); err != nil {
			logrus.WithError(err).Error("Error closing file on import")
		}
	}()

	if err := server.ImportServer(configuration, settings, event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	writeResponse(w, nil)
}

func ExportServerHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if !isValidToken(token, false, false) {
		w.WriteHeader(http.StatusUnauthorized)
		writeResponse(w, nil)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, nil)
		return
	}

	data, err := server.ExportServer(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, nil)
		return
	}

	if _, err := w.Write(data); err != nil {
		logrus.WithError(err).Error("Error writing zip response")
		w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, nil)
		return
	}
}