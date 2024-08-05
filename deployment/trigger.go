package deployment

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/util"
)

// TriggerRequest is the request body for the trigger endpoint.
// It is used to send a deployment trigger to the controller
// from the web user interface or other api clients.
type TriggerRequest struct {
	DeploymentID string `json:"deploymentId"`
}

func Trigger(c echo.Context) error {
	// Deserialize the request body in to the request struct.
	req := &TriggerRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// TODO: Validate the deployment exists, etc.

	// Serialize the request for sending to the controller.
	str, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Publish the message to the controller over RabbitMQ.
	util.Producer.Publish(context.Background(), "controller", "controller", []byte(str))

	return c.NoContent(http.StatusAccepted)
}
