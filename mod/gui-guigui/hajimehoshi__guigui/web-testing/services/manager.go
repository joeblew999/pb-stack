package services

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// ServiceInfo represents information about a managed service
type ServiceInfo struct {
	Name    string
	URL     string
	PID     int32
	Status  string
	Healthy bool
}

// ServiceManager manages multiple services for testing
type ServiceManager struct {
	services map[string]*ServiceInfo
	mutex    sync.RWMutex
}

// NewServiceManager creates a new service manager
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: make(map[string]*ServiceInfo),
	}
}

// RegisterService registers a service for monitoring
func (sm *ServiceManager) RegisterService(name, url string, pid int32) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.services[name] = &ServiceInfo{
		Name:    name,
		URL:     url,
		PID:     pid,
		Status:  "registered",
		Healthy: false,
	}
}

// GetService returns information about a service
func (sm *ServiceManager) GetService(name string) (*ServiceInfo, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	service, exists := sm.services[name]
	return service, exists
}

// GetAllServices returns all registered services
func (sm *ServiceManager) GetAllServices() map[string]*ServiceInfo {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	result := make(map[string]*ServiceInfo)
	for name, service := range sm.services {
		// Create a copy to avoid race conditions
		result[name] = &ServiceInfo{
			Name:    service.Name,
			URL:     service.URL,
			PID:     service.PID,
			Status:  service.Status,
			Healthy: service.Healthy,
		}
	}
	
	return result
}

// CheckHealth checks the health of a specific service
func (sm *ServiceManager) CheckHealth(name string) error {
	sm.mutex.RLock()
	service, exists := sm.services[name]
	sm.mutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("service %s not registered", name)
	}
	
	// Check if process is running
	if service.PID > 0 {
		exists, err := process.PidExists(service.PID)
		if err != nil || !exists {
			sm.updateServiceStatus(name, "process_dead", false)
			return fmt.Errorf("service %s process (PID %d) not running", name, service.PID)
		}
	}
	
	// Check HTTP endpoint if URL is provided
	if service.URL != "" {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		
		resp, err := client.Get(service.URL)
		if err != nil {
			sm.updateServiceStatus(name, "http_error", false)
			return fmt.Errorf("service %s HTTP check failed: %w", name, err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode >= 400 {
			sm.updateServiceStatus(name, "http_error", false)
			return fmt.Errorf("service %s returned status %d", name, resp.StatusCode)
		}
	}
	
	sm.updateServiceStatus(name, "healthy", true)
	return nil
}

// CheckAllHealth checks the health of all registered services
func (sm *ServiceManager) CheckAllHealth() map[string]error {
	sm.mutex.RLock()
	services := make([]string, 0, len(sm.services))
	for name := range sm.services {
		services = append(services, name)
	}
	sm.mutex.RUnlock()
	
	results := make(map[string]error)
	for _, name := range services {
		results[name] = sm.CheckHealth(name)
	}
	
	return results
}

// WaitForHealth waits for all services to become healthy
func (sm *ServiceManager) WaitForHealth(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for services to become healthy")
		case <-ticker.C:
			results := sm.CheckAllHealth()
			allHealthy := true
			
			for name, err := range results {
				if err != nil {
					allHealthy = false
					fmt.Printf("Service %s not healthy: %v\n", name, err)
				}
			}
			
			if allHealthy && len(results) > 0 {
				return nil
			}
		}
	}
}

// WaitForService waits for a specific service to become healthy
func (sm *ServiceManager) WaitForService(name string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for service %s to become healthy", name)
		case <-ticker.C:
			if err := sm.CheckHealth(name); err == nil {
				return nil
			}
		}
	}
}

// updateServiceStatus updates the status of a service
func (sm *ServiceManager) updateServiceStatus(name, status string, healthy bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	if service, exists := sm.services[name]; exists {
		service.Status = status
		service.Healthy = healthy
	}
}

// GetHealthyServices returns a list of healthy services
func (sm *ServiceManager) GetHealthyServices() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	var healthy []string
	for name, service := range sm.services {
		if service.Healthy {
			healthy = append(healthy, name)
		}
	}
	
	return healthy
}

// GetUnhealthyServices returns a list of unhealthy services
func (sm *ServiceManager) GetUnhealthyServices() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	var unhealthy []string
	for name, service := range sm.services {
		if !service.Healthy {
			unhealthy = append(unhealthy, name)
		}
	}
	
	return unhealthy
}

// RemoveService removes a service from monitoring
func (sm *ServiceManager) RemoveService(name string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	delete(sm.services, name)
}

// Clear removes all services from monitoring
func (sm *ServiceManager) Clear() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.services = make(map[string]*ServiceInfo)
}

// PrintStatus prints the status of all services
func (sm *ServiceManager) PrintStatus() {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	fmt.Println("Service Status:")
	fmt.Println("==============")
	
	if len(sm.services) == 0 {
		fmt.Println("No services registered")
		return
	}
	
	for name, service := range sm.services {
		status := "❌"
		if service.Healthy {
			status = "✅"
		}
		
		fmt.Printf("%s %s - %s (%s)\n", status, name, service.URL, service.Status)
		if service.PID > 0 {
			fmt.Printf("   PID: %d\n", service.PID)
		}
	}
}

// AutoDiscoverServices attempts to discover running services on common ports
func (sm *ServiceManager) AutoDiscoverServices(ports []int) {
	for _, port := range ports {
		url := fmt.Sprintf("http://localhost:%d", port)
		
		client := &http.Client{
			Timeout: 2 * time.Second,
		}
		
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			name := fmt.Sprintf("auto-discovered-%d", port)
			sm.RegisterService(name, url, 0)
			fmt.Printf("Auto-discovered service: %s at %s\n", name, url)
		}
	}
}
