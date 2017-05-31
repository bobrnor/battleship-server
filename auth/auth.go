package auth

import (
	"net/http"

	jsonep "git.nulana.com/bobrnor/json-ep.git"
	"go.uber.org/zap"

	"git.nulana.com/bobrnor/battleship-server/db/client"
)

type authData struct {
	ClientUID string `json:"client_uid"`
}

func Handler() http.HandlerFunc {
	var data authData
	return jsonep.Decorate(handle, &data)
}

func handle(data interface{}) interface{} {
	zap.S().Infow("received",
		"data", data,
	)

	auth, ok := data.(*authData)
	if !ok {
		zap.S().Fatalw("Auth data has wrong struct", data)
		return map[string]interface{}{
			"status": -1,
		}
	}

	if len(auth.ClientUID) == 0 {
		zap.S().Warnw("Expected not empty `client_uid` field", auth)
		return map[string]interface{}{
			"status": -1,
		}
	}

	c, err := client.FindByUID(auth.ClientUID)
	if err != nil {
		zap.S().Errorw("Error during finding client", err)
		return map[string]interface{}{
			"status": -1,
		}
	}

	if c == nil {
		newClient := client.Client{
			UID: auth.ClientUID,
		}
		if err := newClient.Save(); err != nil {
			zap.S().Errorw("Error during saving client", err)
		}
	}

	return map[string]interface{}{
		"status": 0,
	}
}
