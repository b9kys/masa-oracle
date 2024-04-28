package gpumanager

import (
	"errors"
	"log"
	"sync"

	cuda "your_project/go-cuda" // Placeholder for actual go-cuda library
)

type GPUManager struct {
	gpus []*cuda.Device
	lock sync.Mutex
}

// NewGPUManager initializes a new GPU manager
func NewGPUManager() *GPUManager {
	return &GPUManager{}
}

// DetectGPUs detects and initializes CUDA-capable GPUs
func (m *GPUManager) DetectGPUs() error {
	devices, err := cuda.GetDevices()
	if err != nil {
		return err
	}
	for _, device := range devices {
		if err := device.Init(); err != nil {
			log.Printf("Failed to initialize device %v: %v", device.ID, err)
			continue
		}
		m.gpus = append(m.gpus, device)
		log.Printf("Initialized GPU: %v", device.ID)
	}
	if len(m.gpus) == 0 {
		return errors.New("no CUDA-capable GPUs found or initialized")
	}
	return nil
}

// ProcessTask attempts to process a task using an available GPU
func (m *GPUManager) ProcessTask(task *Task) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, gpu := range m.gpus {
		if gpu.IsAvailable() {
			return gpu.Process(task.Data)
		}
	}
	return errors.New("no available GPUs")
}
