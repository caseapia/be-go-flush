package config

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type ServiceManager struct {
	configPath string
	config     *ManagerConfig
	mu         sync.RWMutex
	watchers   []chan ServiceEvent
}

type ManagerConfig struct {
	Version  int                       `json:"version"`
	Services map[string]*ServiceStatus `json:"services"`
}

type ServiceStatus struct {
	Enabled   bool              `json:"enabled"`
	UpdatedAt time.Time         `json:"updated_at"`
	UpdatedBy string            `json:"updated_by,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type ServiceEvent struct {
	ServiceName string
	Enabled     bool
	Timestamp   time.Time
}

func NewServiceManager(configPath string) (*ServiceManager, error) {
	sm := &ServiceManager{
		configPath: configPath,
		watchers:   make([]chan ServiceEvent, 0),
	}

	if err := sm.load(); err != nil {
		return nil, err
	}

	return sm, nil
}

func (sm *ServiceManager) load() error {
	data, err := os.ReadFile(sm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			sm.config = &ManagerConfig{
				Version:  1,
				Services: make(map[string]*ServiceStatus),
			}

			return sm.save()
		}

		return err
	}

	var config ManagerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	sm.config = &config
	return nil
}

func (sm *ServiceManager) save() error {
	sm.config.Version++

	data, err := json.MarshalIndent(sm.config, "", "")
	if err != nil {
		return err
	}

	tmpPath := sm.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, sm.configPath)
}

func (sm *ServiceManager) SetServiceEnabled(name string, enabled bool, updatedBy string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	status, exists := sm.config.Services[name]
	if !exists {
		status = &ServiceStatus{}
		sm.config.Services[name] = status
	}

	oldEnabled := status.Enabled
	status.Enabled = enabled
	status.UpdatedAt = time.Now()
	status.UpdatedBy = updatedBy

	if err := sm.save(); err != nil {
		return err
	}

	if oldEnabled != enabled {
		sm.notifyWatchers(ServiceEvent{
			ServiceName: name,
			Enabled:     enabled,
			Timestamp:   status.UpdatedAt,
		})
	}

	return nil
}

func (sm *ServiceManager) IsServiceEnabled(name string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	status, exists := sm.config.Services[name]
	return exists && status.Enabled
}

func (sm *ServiceManager) Subscribe() <-chan ServiceEvent {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	ch := make(chan ServiceEvent, 10)
	sm.watchers = append(sm.watchers, ch)
	return ch
}

func (sm *ServiceManager) notifyWatchers(event ServiceEvent) {
	for _, ch := range sm.watchers {
		select {
		case ch <- event:
		default:
		}
	}
}

func (sm *ServiceManager) PopulateServices() map[string]*ServiceStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	statuses := sm.config.Services

	return statuses
}

func (sm *ServiceManager) Reload() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.load()
}
