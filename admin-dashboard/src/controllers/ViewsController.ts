import { BaseUrl } from "./Constants";
import { ResponseMessage, ViewModel } from "./Types";

export class ViewsController {
  static GetView = async (model: ViewModel): Promise<ResponseMessage> => {
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
    var response = await fetch(BaseUrl + "/api/views", requestOptions).then(
      (response) => {
        return response.json();
      }
    );
    return response;
  };
}
