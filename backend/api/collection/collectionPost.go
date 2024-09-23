package collection

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/types"
	"github.com/xomatix/silly-syntax-backend-bonanza/database"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/permissions"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
	pluginmanager "github.com/xomatix/silly-syntax-backend-bonanza/pluginManager"
)

func InitCollectionPostRoutes(w http.ResponseWriter, r *http.Request) {

	// region Basic read information and insert into database
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

	var query querygenerators.InsertQueryCreator
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

	resolvedPermissionMacro := permissions.ResolvePermissionsMacro(collectionPermisionMacro.Create, userID)
	query.Filter = resolvedPermissionMacro

	generatedInsertQuery, _ := query.InsertQuery()
	insertedID, err := database.ExecuteNonQuery(generatedInsertQuery)
	if err != nil {
		fmt.Println(generatedInsertQuery)
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	// endregion
	resp.Data = insertedID

	// region Trigger before insert
	returningQuery := querygenerators.SelectQueryCreator{
		CollectionName: query.CollectionName,
		ID:             []int64{insertedID},
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
	for k, v := range query.Values {
		updatedRecord[k] = v
	}

	err = triggerBeforeSave(query.CollectionName, originalRecord[0], &updatedRecord)
	if err != nil {
		deleteQuery := querygenerators.DeleteQueryCreator{
			CollectionName: query.CollectionName,
			ID:             insertedID,
		}
		generatedDeleteQuery, _ := deleteQuery.DeleteQuery()
		database.ExecuteNonQuery(generatedDeleteQuery)

		resp.Success = false
		resp.Message = "Record not inserted. " + err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	triggerAfterSave(query.CollectionName, originalRecord[0], &updatedRecord)

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
}

func triggerBeforeSave(collectionName string, originalObj map[string]interface{}, obj *map[string]interface{}) error {
	funcsToCall := pluginmanager.GetPluginLoader().Triggers[collectionName]["before_insert"]
	for _, f := range funcsToCall {
		err := f(originalObj, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func triggerAfterSave(collectionName string, originalObj map[string]interface{}, obj *map[string]interface{}) {
	funcsToCall := pluginmanager.GetPluginLoader().Triggers[collectionName]["after_insert"]
	for _, f := range funcsToCall {
		f(originalObj, obj)
	}
}
