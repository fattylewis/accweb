package server

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	serverList []ServerSettings
	m          sync.Mutex
)

func LoadServerList() {
	logrus.Info("Loading server list...")
	dir, err := ioutil.ReadDir(os.Getenv("ACCWEB_CONFIG_PATH"))

	if err != nil {
		logrus.WithError(err).Fatal("Error opening config directory to initialize server list")
	}

	for _, entry := range dir {
		if entry.IsDir() {
			if err := loadServerSettings(entry.Name()); err != nil {
				logrus.WithError(err).WithField("name", entry.Name()).Fatal("Error loading server settings")
			}
		}
	}

	logrus.WithField("servers", len(serverList)).Info("Server list loaded successfully")
}

func loadServerSettings(name string) error {
	server := &ServerSettings{Id: parseServerId(name)}

	if err := loadConfigFromFile(&server.Configuration, filepath.Join(os.Getenv("ACCWEB_CONFIG_PATH"), name, configurationJsonName)); err != nil {
		return err
	}

	if err := loadConfigFromFile(&server.Settings, filepath.Join(os.Getenv("ACCWEB_CONFIG_PATH"), name, settingsJsonName)); err != nil {
		return err
	}

	if err := loadConfigFromFile(&server.Event, filepath.Join(os.Getenv("ACCWEB_CONFIG_PATH"), name, eventJsonName)); err != nil {
		return err
	}

	addServer(server)
	return nil
}

func parseServerId(name string) int {
	id, err := strconv.Atoi(name)

	if err != nil {
		logrus.WithError(err).Fatal("Error parsing server ID from directory name")
	}

	return id
}

func loadConfigFromFile(config interface{}, path string) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

func GetServerList() []ServerSettings {
	return serverList
}

func GetServerById(id int) *ServerSettings {
	m.Lock()
	defer m.Unlock()

	for _, server := range serverList {
		if server.Id == id {
			return &server
		}
	}

	return nil
}

func setServer(server *ServerSettings) {
	m.Lock()
	defer m.Unlock()

	for i, s := range serverList {
		if s.Id == server.Id {
			serverList[i] = *server
			break
		}
	}
}

func addServer(server *ServerSettings) {
	m.Lock()
	defer m.Unlock()
	serverList = append(serverList, *server)
}

func removeServer(id int) {
	m.Lock()
	defer m.Unlock()

	for i, s := range serverList {
		if s.Id == id {
			serverList = append(serverList[:i], serverList[i+1:]...)
			break
		}
	}
}