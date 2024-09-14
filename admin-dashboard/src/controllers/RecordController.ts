import { BaseUrl } from "./Constants";
import {
  DataDeleteModel,
  DataInsertModel,
  DataListModel,
  DataUpdateModel,
  ResponseMessage,
} from "./Types";

export class RecordController {
  static GetRecords = async (
    model: DataListModel
  ): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    const token = localStorage.getItem("bonanza_token");
    myHeaders.append("Authorization", token ? token : "");
    const requestOptions: any = {
      method: "POST",
      headers: myHeaders,
      redirect: "follow",
      body: JSON.stringify(model),
      Credentials: "include",
    };
    var response = await fetch(
      BaseUrl + "/api/collection/list",
      requestOptions
    ).then((response) => {
      return response.json();
    });
    return response;
  };

  static InsertData = async (
    model: DataInsertModel
  ): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    const token = localStorage.getItem("bonanza_token");
    myHeaders.append("Authorization", token ? token : "");
    const requestOptions: any = {
      method: "POST",
      headers: myHeaders,
      redirect: "follow",
      body: JSON.stringify(model),
      Credentials: "include",
    };
    var response = await fetch(
      BaseUrl + "/api/collection/insert",
      requestOptions
    ).then((response) => {
      return response.json();
    });
    return response;
  };

  static UpdateData = async (
    model: DataUpdateModel
  ): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    const token = localStorage.getItem("bonanza_token");
    myHeaders.append("Authorization", token ? token : "");
    const requestOptions: any = {
      method: "POST",
      headers: myHeaders,
      redirect: "follow",
      body: JSON.stringify(model),
      Credentials: "include",
    };
    var response = await fetch(
      BaseUrl + "/api/collection/update",
      requestOptions
    ).then((response) => {
      return response.json();
    });
    return response;
  };

  static DeleteData = async (
    model: DataDeleteModel
  ): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    const token = localStorage.getItem("bonanza_token");
    myHeaders.append("Authorization", token ? token : "");
    const requestOptions: any = {
      method: "POST",
      headers: myHeaders,
      redirect: "follow",
      body: JSON.stringify(model),
      Credentials: "include",
    };
    var response = await fetch(
      BaseUrl + "/api/collection/delete",
      requestOptions
    ).then((response) => {
      response.headers.getSetCookie().map((x) => {
        document.cookie = x;
      });
      return response.json();
    });
    return response;
  };
}
