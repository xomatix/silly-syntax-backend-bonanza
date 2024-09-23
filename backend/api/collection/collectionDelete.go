package collection

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xomatix/silly-syntax-backend-bonanza/database"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/permissions"
	pluginmanager "github.com/xomatix/silly-syntax-backend-bonanza/pluginManager"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/types"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
)

func InitCollectionDeleteRoutes(w http.ResponseWriter, r *http.Request) {
	// region basic data collection
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

	var deleteQuery querygenerators.DeleteQueryCreator
	err = json.Unmarshal(body, &deleteQuery)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	collectionPermisionMacro, err := permissions.GetTablePermissions(deleteQuery.CollectionName)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	resolvedPermissionMacro := permissions.ResolvePermissionsMacro(collectionPermisionMacro.Delete, userID)

	deleteQuery.Filter = resolvedPermissionMacro

	generatedDeleteQuery, err := deleteQuery.DeleteQuery()
	if err != nil {
		fmt.Println(err)
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}
	// endregion

	// region trigger and delete
	returningQuery := querygenerators.SelectQueryCreator{
		CollectionName: deleteQuery.CollectionName,
		ID:             []int64{deleteQuery.ID},
	}
	generatedReturningQuery, err := returningQuery.GetQuery()
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	originalRecord, err := database.ExecuteQuery(generatedReturningQuery)
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
	// no need to copy this because we will not change it we just delete

	err = triggerBeforeDelete(deleteQuery.CollectionName, originalRecord[0], &updatedRecord)
	if err != nil {
		resp.Success = false
		resp.Message = "Record not deleted. " + err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}
	// endregion

	_, err = database.ExecuteNonQuery(generatedDeleteQuery)
	if err != nil {
		fmt.Println(generatedDeleteQuery)
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	jsonString, err := json.Marshal(resp)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	if deleteQuery.CollectionName == "tables_permissions" {
		permissions.LoadTablesPermissions()
	}

	if userID != 0 {
		authentication.SetAuthenticationCookies(w, userID, username)
	}
	w.Write(jsonString)
}

func triggerBeforeDelete(collectionName string, originalObj map[string]interface{}, obj *map[string]interface{}) error {
	funcsToCall := pluginmanager.GetPluginLoader().Triggers[collectionName]["before_delete"]
	for _, f := range funcsToCall {
		err := f(originalObj, obj)
		if err != nil {
			return err
		}
	}
	return nil
}
