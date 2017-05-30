package auth

import (
	"net/http"

	jsonep "git.nulana.com/bobrnor/json-ep.git"
	"github.com/hashicorp/packer/common/uuid"
	"go.uber.org/zap"

	"git.nulana.com/bobrnor/battleship-server/db"
)

type authData struct {
	ClientID string `json:"client_id"`
}

func Handler() http.HandlerFunc {
	var data authData
	return jsonep.Decorate(handle, &data)
}

func handle(data interface{}) interface{} {
	zap.S().Infow("received",
		"data", data,
	)

	mapData := data.(*authData)

	session := db.Session{
		ClientID:  mapData.ClientID,
		SessionID: uuid.TimeOrderedUUID(),
	}

	if err := session.Save(); err != nil {
		zap.S().Errorw(err.Error())
	}

	return map[string]interface{}{
		"session_id": session.SessionID,
	}
}
