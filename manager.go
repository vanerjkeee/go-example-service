package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Manager struct {
	rep Repository
}

func (m *Manager) Add(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body []Task
	err := decoder.Decode(&body)
	if err != nil {
		m.responseBadRequestError(w)
		return
	}
	for _, b := range body {
		if b.NumberOfRequests <= 0 || b.Url == "" {
			m.responseBadRequestError(w)
			return
		}
	}

	tasks, err := m.rep.AddTasks(body)
	if err != nil {
		log.Println(err)
		m.responseInternalError(w)
		return
	}

	for _, task := range tasks {
		go m.processTask(task)
	}

	m.responseOK(w)
}

func (m *Manager) Status(w http.ResponseWriter, req *http.Request) {
	status, err := m.rep.GetStatus()
	if err != nil {
		log.Println(err)
		m.responseInternalError(w)
		return
	}

	body, err := json.Marshal(status)
	if err != nil {
		log.Println(err)
		m.responseInternalError(w)
		return
	}

	m.responseOKWithBody(w, body)
}

func (m *Manager) processTask(task Task) {
	for i := 0; i < task.NumberOfRequests; i++ {
		resp, err := http.Get(task.Url)
		if err == nil {
			log.Println(task.Url, " success")
			err = m.rep.IncSuccessTask(task)
		} else {
			err = m.rep.IncErrorTask(task)
			log.Println(task.Url, " error")
		}

		if resp != nil {
			_ = resp.Body.Close()
		}
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (m *Manager) responseOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (m *Manager) responseOKWithBody(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (m *Manager) responseInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}

func (m *Manager) responseBadRequestError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
}
