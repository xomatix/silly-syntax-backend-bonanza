package collection

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/types"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_functions"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/permissions"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
	pluginfunctions "github.com/xomatix/silly-syntax-backend-bonanza/pluginManager/plugin_functions"
)

func InitCollectionUpdateRoutes(w http.ResponseWriter, r *http.Request) {
	//region Basic read information
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
	}

	userID, username, err := authentication.ResolveCookiesWithUserFromHeader(r)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	var query querygenerators.UpdateQueryCreator
	err = json.Unmarshal(body, &query)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	collectionPermisionMacro, err := permissions.GetTablePermissions(query.CollectionName)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	resolvedPermissionMacro := permissions.ResolvePermissionsMacro(collectionPermisionMacro.Update, userID)
	query.Filter = resolvedPermissionMacro

	generatedReturningQuery, err := query.SelectQuery()
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}
	//endregion

	//region update exec with trigger
	originalRecord, err := database_functions.ExecuteQuery(generatedReturningQuery)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	//so we have copy here
	updatedRecord := make(map[string]interface{})
	for k, v := range originalRecord[0] {
		updatedRecord[k] = v
	}
	for k, v := range query.Values {
		updatedRecord[k] = v
	}

	err = triggerBeforeUpdate(query.CollectionName, originalRecord[0], &updatedRecord)
	if err != nil {
		resp.Success = false
		resp.Message = "Record not updated. " + err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	generatedUpdateQuery, err := query.UpdateQuery()
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	_, err = database_functions.ExecuteNonQuery(generatedUpdateQuery)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	triggerAfterUpdate(query.CollectionName, originalRecord[0], &updatedRecord)
	//endregion

	//region returning data with auth
	jsonString, err := json.Marshal(resp)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	if query.CollectionName == "tables_permissions" {
		permissions.LoadTablesPermissions()
	}

	if userID != 0 {
		authentication.SetAuthenticationCookies(w, userID, username)
	}
	w.Write(jsonString)
	//endregion
}

func triggerBeforeUpdate(collectionName string, originalObj map[string]interface{}, obj *map[string]interface{}) error {
	funcsToCall := pluginfunctions.GetPluginLoader().Triggers[collectionName]["before_update"]
	for _, f := range funcsToCall {
		err := f(originalObj, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func triggerAfterUpdate(collectionName string, originalObj map[string]interface{}, obj *map[string]interface{}) {
	funcsToCall := pluginfunctions.GetPluginLoader().Triggers[collectionName]["after_update"]
	for _, f := range funcsToCall {
		f(originalObj, obj)
	}
}
