package main

import (
	"context"
	"fmt"

	"github.com/project-flogo/contrib/activity/log"
	"github.com/project-flogo/contrib/trigger/rest"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/engine"
)

func main() {

	app := myApp()

	e, err := api.NewEngine(app)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	engine.RunEngine(e)
}

func myApp() *api.App {
	app := api.NewApp()

	trg := app.NewTrigger(&rest.Trigger{}, &rest.Settings{Port: 8080})
	h, _ := trg.NewHandler(&rest.HandlerSettings{Method: "GET", Path: "/blah/:num"})
	h.NewAction(RunActivities)

	//store in map to avoid activity instance recreation
	logAct, _ := api.NewActivity(&log.Activity{})
	activities = map[string]activity.Activity{"log": logAct}

	return app
}

var activities map[string]activity.Activity

func RunActivities(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {

	trgOut := &rest.Output{}
	trgOut.FromMap(inputs)

	msg, _ := coerce.ToString(trgOut.PathParams)
	_, err := api.EvalActivity(activities["log"], &log.Input{Message: msg})
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{})

	response["id"] = "123"
	response["amount"] = "1"
	response["balance"] = "500"
	response["currency"] = "USD"

	reply := &rest.Reply{Code: 200, Data: response}
	return reply.ToMap(), nil
}
