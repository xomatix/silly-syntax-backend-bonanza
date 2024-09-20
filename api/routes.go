package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/auth"
	"github.com/xomatix/silly-syntax-backend-bonanza/api/collection"
	"github.com/xomatix/silly-syntax-backend-bonanza/api/types"
	"github.com/xomatix/silly-syntax-backend-bonanza/api/views"
	"github.com/xomatix/silly-syntax-backend-bonanza/database"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/permissions"
)

func InitApiRoutes(mux *http.ServeMux) {
	initDatabaseManipulationRoutes(mux)
	initCollectionManipulationRoutes(mux)
	initAuthenticationRoutes(mux)
	initViewsRoutes(mux)
}

func initCollectionManipulationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/collection/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			collection.InitCollectionGetRoutes(w, r)
		}

	})

	mux.HandleFunc("/api/collection/insert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			collection.InitCollectionPostRoutes(w, r)
		}
	})

	mux.HandleFunc("/api/collection/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			collection.InitCollectionUpdateRoutes(w, r)
		}
	})

	mux.HandleFunc("/api/collection/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			collection.InitCollectionDeleteRoutes(w, r)
		}
	})
}

func initAuthenticationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			auth.InitAuthenticationLoginRoutes(w, r)
		}
	})
}

func initDatabaseManipulationRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/api/tables", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tabConfArr := make([]database.TableConfig, 0)
			for _, v := range database.GetTablesConfig() {
				// if v.Name == "tables_permissions" {
				// 	continue
				// }
				tabConfArr = append(tabConfArr, v)
			}

			resp := types.ResponseMessage{
				Success: true,
				Data:    tabConfArr,
				Message: "success",
			}
			jsonString, err := json.Marshal(resp)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				w.Write(jsonString)
			}
			w.Write(jsonString)
		}
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			resp := types.ResponseMessage{
				Success: true,
				Data:    nil,
				Message: "success",
			}
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			var tableConf database.TableConfig
			err = json.Unmarshal(body, &tableConf)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			err = database.CrateTable(tableConf.Name)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}

			permissions.LoadTablesPermissions()

			resp.Message = "created table: " + tableConf.Name
			jsonString, _ := json.Marshal(resp)
			w.Write(jsonString)
		}
		if r.Method == http.MethodPut {
			body, err := io.ReadAll(r.Body)
			resp := types.ResponseMessage{
				Success: true,
				Data:    nil,
				Message: "success",
			}
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			var tabConf database.TableConfig
			err = json.Unmarshal(body, &tabConf)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			for _, v := range tabConf.Columns {
				err = database.AddColumnToTable(tabConf.Name, v)
				if err != nil {
					resp.Success = false
					resp.Message = err.Error()
					jsonString, _ := json.Marshal(resp)
					w.Write(jsonString)
					return
				}
			}

			resp.Message = "added column to table: " + tabConf.Name + " successfully"
			jsonString, _ := json.Marshal(resp)
			w.Write(jsonString)
		}
		if r.Method == http.MethodDelete {
			body, err := io.ReadAll(r.Body)
			resp := types.ResponseMessage{
				Success: true,
				Data:    nil,
				Message: "success",
			}
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			var tabConf database.TableConfig
			err = json.Unmarshal(body, &tabConf)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			for _, v := range tabConf.Columns {
				err = database.RemoveColumnFromTable(tabConf.Name, v)
				if err != nil {
					resp.Success = false
					resp.Message = err.Error()
					jsonString, _ := json.Marshal(resp)
					w.Write(jsonString)
					return
				}
			}

			permissions.LoadTablesPermissions()

			resp.Message = "removed column from table: " + tabConf.Name + " successfully"
			jsonString, _ := json.Marshal(resp)
			w.Write(jsonString)
		}
	})

	mux.HandleFunc("/api/tables/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			body, err := io.ReadAll(r.Body)
			resp := types.ResponseMessage{
				Success: true,
				Data:    nil,
				Message: "success",
			}
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			var tabConf database.TableConfig
			err = json.Unmarshal(body, &tabConf)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}
			err = database.DeleteTable(tabConf.Name)
			if err != nil {
				resp.Success = false
				resp.Message = err.Error()
				jsonString, _ := json.Marshal(resp)
				w.Write(jsonString)
				return
			}

			permissions.LoadTablesPermissions()

			resp.Message = "removed table: " + tabConf.Name + " successfully"
			jsonString, _ := json.Marshal(resp)
			w.Write(jsonString)

		}
	})
}

func initViewsRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/api/views", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			views.InitViewsRoutes(w, r)
		}
	})
}
