package collection

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xomatix/silly-syntax-backend-bonanza/database"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/permissions"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/types"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
)

func InitCollectionDeleteRoutes(w http.ResponseWriter, r *http.Request) {

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

	var query querygenerators.DeleteQueryCreator
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

	resolvedPermissionMacro := permissions.ResolvePermissionsMacro(collectionPermisionMacro.Delete, userID)

	query.Filter = resolvedPermissionMacro

	generatedDeleteQuery, err := query.DeleteQuery()
	if err != nil {
		fmt.Println(err)
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

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

	if query.CollectionName == "tables_permissions" {
		permissions.LoadTablesPermissions()
	}

	if userID != 0 {
		authentication.SetAuthenticationCookies(w, userID, username)
	}
	w.Write(jsonString)
}
