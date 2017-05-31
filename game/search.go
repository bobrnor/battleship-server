package game

import (
	"net/http"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	jsonep "git.nulana.com/bobrnor/json-ep.git"
	"go.uber.org/zap"
)

type searchData struct {
	ClientUID string `json:"client_uid"`
}

func Handler() http.HandlerFunc {
	var data searchData
	return jsonep.Decorate(handle, &data)
}

func handle(data interface{}) interface{} {
	zap.S().Infow("received",
		"data", data,
	)

	search, ok := data.(*searchData)
	if !ok {
		zap.S().Fatalw("Search data has wrong struct", data)
		return map[string]interface{}{
			"status": -1,
		}
	}

	if len(search.ClientUID) == 0 {
		zap.S().Warnw("Expected not empty `client_uid` field", search)
		return map[string]interface{}{
			"status": -1,
		}
	}

	c, err := client.FindByUID(search.ClientUID)
	if err != nil {
		zap.S().Errorw("Error during finding client", err)
		return map[string]interface{}{
			"status": -1,
		}
	}

	if c == nil {
		zap.S().Errorw("Client not found", err)
		return map[string]interface{}{
			"status": -1,
		}
	}

	// TODO: find request or create new (!!!: concurrency)
	// if request found:
	// - delete it
	// - create room
	// - return room uid

	return map[string]interface{}{
		"status": 0,
	}
}
