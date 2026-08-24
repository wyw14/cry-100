package dispatch

import "github.com/wyw14/cry-100/internal/model"

func operationIdentity(request model.DispatchRequest) model.OperationID {
	if request.OperationID != "" {
		return request.OperationID
	}
	return model.NewOperationID()
}
