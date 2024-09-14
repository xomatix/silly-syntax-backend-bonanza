import { BaseUrl } from "./Constants";
import { LoginModel, ResponseMessage } from "./Types";

export class AuthController {
  static Login = async (model: LoginModel): Promise<ResponseMessage> => {
    const myHeaders = new Headers();

    const raw = JSON.stringify(model);

    myHeaders.append("Content-Type", "application/json");
    const requestOptions: any = {
      method: "POST",
      headers: myHeaders,
      body: raw,
      redirect: "follow",
      Credentials: "include",
    };
    var response = await fetch(
      BaseUrl + "/api/auth/login",
      requestOptions
    ).then((response) => {
      return response.json();
    });

    localStorage.setItem("bonanza_token", response.data);
    return response;
  };
}
