package workers

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
)

// HandleConnect is a method of the Worker struct that handles the connection of a worker.
// It takes in an actor context and a Connect message as parameters.
func (a *Worker) HandleConnect(ctx actor.Context, m *messages.Connect) {
	logrus.Infof("[+] Worker %v connected", m.Sender)
	clients.Add(m.Sender)
}

// HandleLog is a method of the Worker struct that handles logging.
// It takes in an actor context and a string message as parameters.
func (a *Worker) HandleLog(ctx actor.Context, l string) {
	logrus.Info(l)
}

// HandleWork is a method of the Worker struct that handles the work assigned to a worker.
// It takes in an actor context and a Work message as parameters.
// @todo fire data to masa sdk
func (a *Worker) HandleWork(ctx actor.Context, m *messages.Work) {
	var workData map[string]string
	err := json.Unmarshal([]byte(m.Data), &workData)
	if err != nil {
		logrus.Errorf("Error parsing work data: %v", err)
		return
	}

	var bodyData map[string]interface{}
	if workData["body"] != "" {
		if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
			logrus.Errorf("Error unmarshalling body: %v", err)
			return
		}
	}
	response, err := GetWorkHandlerManager().ExecuteWork(workData["request"], workData["request_id"], m.Id, bodyData)
	if err != nil {
		logrus.Errorf("Error processing request: %v", err)
		return
	}
	ctx.Respond(response)
	ctx.Poison(ctx.Self())
}
