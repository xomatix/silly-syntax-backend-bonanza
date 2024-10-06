package collection

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/types"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_functions"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/permissions"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
)

func InitCollectionGetRoutes(w http.ResponseWriter, r *http.Request) {

	resp := types.ResponseMessage{
		Success: true,
		Data:    "",
		Message: "success",
	}

	userID, username, err := authentication.ResolveCookiesWithUserFromHeader(r)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	var query querygenerators.SelectQueryCreator
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

	resolvedPermissionMacro := permissions.ResolvePermissionsMacro(collectionPermisionMacro.Read, userID)

	if len(query.Filter) > 0 && len(resolvedPermissionMacro) > 0 {
		query.Filter = fmt.Sprintf("(%s) AND (%s)", resolvedPermissionMacro, query.Filter)
	} else if len(resolvedPermissionMacro) > 0 {
		query.Filter = resolvedPermissionMacro
	}

	data, err := ControllerGet(query)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	resp.Data = data

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
}

func ControllerGet(query querygenerators.SelectQueryCreator) ([]map[string]interface{}, error) {
	generatedSelectQuery, _ := query.GetQuery()

	result, err := database_functions.ExecuteQuery(generatedSelectQuery)
	if err != nil {
		return nil, err
	}

	var rows []map[string]interface{}
	rows = append(rows, result...)

	return rows, err
}
