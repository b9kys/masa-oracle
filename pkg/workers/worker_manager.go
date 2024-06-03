package workers

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	pubsub2 "github.com/libp2p/go-libp2p-pubsub"

	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
)

var (
	instance *WorkHandlerManager
	once     sync.Once
)

func GetWorkHandlerManager() *WorkHandlerManager {
	once.Do(func() {
		instance = NewWorkHandlerManager()
	})
	return instance
}

// ErrHandlerNotFound is an error returned when a work handler cannot be found.
var ErrHandlerNotFound = errors.New("work handler not found")

// WorkHandler defines the interface for handling different types of work.
type WorkHandler interface {
	HandleWork(data map[string]interface{}) (interface{}, error)
}

// WorkHandlerInfo contains information about a work handler, including metrics.
type WorkHandlerInfo struct {
	Handler      WorkHandler
	CallCount    int64
	TotalRuntime time.Duration
}

// WorkHandlerManager manages work handlers and tracks their execution metrics.
type WorkHandlerManager struct {
	handlers map[string]*WorkHandlerInfo
	mu       sync.RWMutex
}

// NewWorkHandlerManager creates a new instance of WorkHandlerManager.
func NewWorkHandlerManager() *WorkHandlerManager {
	return &WorkHandlerManager{
		handlers: make(map[string]*WorkHandlerInfo),
	}
}

// AddWorkHandler registers a new work handler under a specific name.
func (whm *WorkHandlerManager) AddWorkHandler(name string, handler WorkHandler) {
	whm.mu.Lock()
	defer whm.mu.Unlock()
	whm.handlers[name] = &WorkHandlerInfo{Handler: handler}
}

// GetWorkHandler retrieves a registered work handler by name.
func (whm *WorkHandlerManager) GetWorkHandler(name string) (WorkHandler, bool) {
	whm.mu.RLock()
	defer whm.mu.RUnlock()
	info, exists := whm.handlers[name]
	if !exists {
		return nil, false
	}
	return info.Handler, true
}

// ExecuteWork finds and executes the work handler associated with the given name.
// It tracks the call count and execution duration for the handler.
func (whm *WorkHandlerManager) ExecuteWork(name, requestId, messageId string, data map[string]interface{}) (*messages.Response, error) {
	handler, exists := whm.GetWorkHandler(name)
	if !exists {
		return nil, ErrHandlerNotFound
	}

	startTime := time.Now()

	result, err := handler.HandleWork(data)
	if err != nil {
		return nil, err
	}
	chanResponse := ChanResponse{
		Response:  map[string]interface{}{"data": result},
		ChannelId: requestId,
	}
	val := &pubsub2.Message{
		ValidatorData: chanResponse,
		ID:            messageId,
	}
	jsn, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	response := &messages.Response{RequestId: requestId, Value: string(jsn)}

	duration := time.Since(startTime)

	whm.mu.Lock()
	handlerInfo := whm.handlers[name]
	handlerInfo.CallCount++
	handlerInfo.TotalRuntime += duration
	whm.mu.Unlock()

	return response, err
}
