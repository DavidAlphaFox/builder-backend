// Copyright 2022 The ILLA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"errors"

	ws "github.com/illa-family/builder-backend/internal/websocket"
	"github.com/illa-family/builder-backend/pkg/user"
)

func SignalEnter(hub *ws.Hub, message *ws.Message, ai *user.AuthenticatorImpl) error {
	// init
	currentClient := hub.Clients[message.ClientID]
	var ok bool
	if len(message.Payload) == 0 {
		err := errors.New("[websocket-server] websocket protocol syntax error.")
		currentClient.Feedback(message, ws.ERROR_CODE_LOGIN_FAILED, err)
		return err
	}
	var authToken map[string]interface{}
	if authToken, ok = message.Payload[0].(map[string]interface{}); !ok {
		err := errors.New("[websocket-server] websocket protocol syntax error.")
		currentClient.Feedback(message, ws.ERROR_CODE_LOGIN_FAILED, err)
		return err
	}
	token, _ := authToken["authToken"].(string)

	// convert authToken to uid
	userID, userUID, extractErr := ai.ExtractUserIDFromToken(token)
	if extractErr != nil {
		return extractErr
	}
	validAccessToken, validaAccessErr := ai.ValidateAccessToken(token)
	validUser, validUserErr := ai.ValidateUser(userID, userUID)
	if validaAccessErr != nil || validUserErr != nil || !validAccessToken || !validUser {
		err := errors.New("[websocket-server] access token invalid.")
		currentClient.Feedback(message, ws.ERROR_CODE_LOGIN_FAILED, err)
		return err
	}
	// assign logged in and mapped user id
	currentClient.IsLoggedIn = true
	currentClient.MappedUserID = userID
	currentClient.Feedback(message, ws.ERROR_CODE_LOGGEDIN, nil)
	return nil

}
