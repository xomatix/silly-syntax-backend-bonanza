import { BaseUrl } from "./Constants";
import { ResponseMessage, TableConfig } from "./Types";

export class CollectionController {
  static collectionsList = async (): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    const requestOptions: any = {
      method: "GET",
      headers: myHeaders,
      //body: raw,
      redirect: "follow",
    };
    var response = await fetch(BaseUrl + "/api/tables", requestOptions).then(
      (response) => {
        return response.json();
      }
    );
    return response;
  };

  static addCollection = async (name: string): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    var raw: TableConfig = {
      name: name,
      columns: [],
    };
    const requestOptions: any = {
      method: "POST",
      headers: myHeaders,
      body: JSON.stringify(raw),
      redirect: "follow",
    };
    var response = await fetch(BaseUrl + "/api/tables", requestOptions).then(
      (response) => {
        return response.json();
      }
    );
    return response;
  };

  /**
   * Adds a new column to the table based on the provided TableConfig.
   *
   * @param {TableConfig} column - The table configuration to be added.
   * {ColumnConfig} column - The column configuration to be added to the table.
   * Can add multiple columns 1 recommended.
   * @return {Promise<ResponseMessage>} A Promise that resolves with the response message after adding the column.
   */
  static addColumn = async (column: TableConfig): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    var raw: TableConfig = column;
    const requestOptions: any = {
      method: "PUT",
      headers: myHeaders,
      body: JSON.stringify(raw),
      redirect: "follow",
    };
    var response = await fetch(BaseUrl + "/api/tables", requestOptions).then(
      (response) => {
        return response.json();
      }
    );
    return response;
  };

  static deleteColumn = async (
    column: TableConfig
  ): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    var raw: TableConfig = column;
    const requestOptions: any = {
      method: "DELETE",
      headers: myHeaders,
      body: JSON.stringify(raw),
      redirect: "follow",
    };
    var response = await fetch(BaseUrl + "/api/tables", requestOptions).then(
      (response) => {
        return response.json();
      }
    );
    return response;
  };

  /**
   * Deletes the collection with the given name.
   *
   * @param {TableConfig} tableConfig - The name of the collection to be deleted is needed as parameter.
   * @return {Promise<ResponseMessage>} A Promise that resolves with the response message after deleting the collection.
   */
  static deleteCollection = async (
    tableConfig: TableConfig
  ): Promise<ResponseMessage> => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    var raw: TableConfig = tableConfig;
    const requestOptions: any = {
      method: "DELETE",
      headers: myHeaders,
      body: JSON.stringify(raw),
      redirect: "follow",
    };
    var response = await fetch(
      BaseUrl + "/api/tables/delete",
      requestOptions
    ).then((response) => {
      return response.json();
    });
    return response;
  };
}
