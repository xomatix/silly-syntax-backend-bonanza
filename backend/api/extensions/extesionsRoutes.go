package extensions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/types"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	pluginfunctions "github.com/xomatix/silly-syntax-backend-bonanza/pluginManager/plugin_functions"
)

// InitExtensionsRoutes registers the /api/{custom_name} route, which
// handles POST requests and routes them to the appropriate handler from plugins.
func InitExtensionsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// region Basic read information
			body, err := io.ReadAll(r.Body)
			resp := types.ResponseMessage{
				Success: true,
				Data:    "",
				Message: "success",
			}
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonStr, _ := json.Marshal(resp)
				w.Write(jsonStr)
				return
			}

			userID, username, err := authentication.ResolveCookiesWithUserFromHeader(r)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonStr, _ := json.Marshal(resp)
				w.Write(jsonStr)
				return
			}

			requestedApiMethod := strings.Replace(r.URL.Path, "/api/", "", 1)

			if requestedApiMethod == "" {
				resp.Success = false
				resp.Message = fmt.Sprintf("Requested API method not found: '%s'", requestedApiMethod)
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			// endregion

			pathHandler, ok := pluginfunctions.GetPluginLoader().Api[requestedApiMethod]
			if !ok {
				resp.Success = false
				resp.Message = fmt.Sprintf("Requested API method not found: '%s'", requestedApiMethod)
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}

			var bodyDecompiled map[string]interface{}
			err = json.Unmarshal(body, &bodyDecompiled)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonStr, _ := json.Marshal(resp)
				w.Write(jsonStr)
				return
			}

			handlerResponse, err := pathHandler(bodyDecompiled, userID)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonStr, _ := json.Marshal(resp)
				w.Write(jsonStr)
				return
			}

			//region returning data with auth
			resp.Data = handlerResponse
			jsonString, err := json.Marshal(resp)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonStr, _ := json.Marshal(resp)
				w.Write(jsonStr)
				return
			}

			if userID != 0 {
				authentication.SetAuthenticationCookies(w, userID, username)
			}
			w.Write(jsonString)
			//endregion

		}
	})

}
