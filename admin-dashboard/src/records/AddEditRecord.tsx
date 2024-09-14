import { useEffect, useState } from "react";
import {
  DataDeleteModel,
  DataInsertModel,
  DataListModel,
  DataType,
  DataUpdateModel,
  RecordModel,
  ResponseMessage,
  TableConfig,
} from "../controllers/Types";
import Dialog from "../components/Dialog";
import SelectRecordReference from "./SelectRecordReference";
import { RecordController } from "../controllers/RecordController";

function AddEditRecord({
  CollectionConfig,
  onClose,
  recordID,
}: {
  CollectionConfig: TableConfig;
  onClose: () => void;
  recordID: number;
}) {
  const [recordData, setRecordData] = useState<RecordModel>({});
  const [isSelectReferenceOpen, setSelectReferenceOpen] =
    useState<boolean>(false);
  const [referenceCollection, setReferenceCollection] = useState<string>("");
  const [referenceField, setReferenceField] = useState<string>("");

  useEffect(() => {
    if (recordID != 0 && recordID != undefined && recordID != null) {
      downloadRecord();
    } else {
      let blankRecord: RecordModel = {};

      for (let index = 0; index < CollectionConfig.columns.length; index++) {
        const column = CollectionConfig.columns[index];
        blankRecord[column.name] =
          column.dataType == DataType.BOOLEAN
            ? false
            : column.dataType == DataType.INTEGER ||
              column.dataType == DataType.DOUBLE
            ? 0
            : "";
      }

      setRecordData(blankRecord);
    }
  }, []);

  const downloadRecord = async () => {
    let model: DataListModel = {
      collectionName: CollectionConfig.name,
      ID: [recordID],
    };
    var result = await RecordController.GetRecords(model);
    if (result.data.length > 0) {
      const downloadedRecord = result.data[0];
      if (
        downloadedRecord["password"] != null &&
        downloadedRecord["password"] != undefined
      ) {
        downloadedRecord["password"] = null;
      }
      setRecordData(downloadedRecord);
    }
  };

  const handleSubmit = async () => {
    let response: ResponseMessage = {
      success: false,
      message: "",
      data: null,
    };
    if (recordID != 0 && recordID != undefined && recordID != null) {
      response = await handleUpdateRecord();
    } else {
      response = await handleNewRecord();
    }
    if (!response.success) {
      alert(`Record not inserted! ${response.message}`);
    } else onClose();
  };

  const handleDelete = async () => {
    if (recordID != 0 && recordID != undefined && recordID != null) {
      if (
        confirm(
          "Are you sure you want to delete this record?\nThere may be other records that depend on this one"
        )
      ) {
        let model: DataDeleteModel = {
          collectionName: CollectionConfig.name,
          ID: recordID,
        };
        let result: ResponseMessage = await RecordController.DeleteData(model);
        if (result.success == false) {
          alert(result.message);
          return;
        }
      }
    }
    onClose();
  };

  const handleNewRecord = async (): Promise<ResponseMessage> => {
    let recordModel: RecordModel = {
      ...recordData,
    };
    let model: DataInsertModel = {
      collectionName: CollectionConfig.name,
      values: recordModel,
    };
    let response = await RecordController.InsertData(model);
    return response;
  };

  const handleUpdateRecord = async (): Promise<ResponseMessage> => {
    let recordModel: RecordModel = {
      ...recordData,
    };
    let model: DataUpdateModel = {
      collectionName: CollectionConfig.name,
      values: recordModel,
      ID: recordID,
    };
    console.log(model);
    let response = await RecordController.UpdateData(model);
    return response;
  };

  const handleSelectReferenceOpen = (collectionName: string, field: string) => {
    setReferenceCollection(collectionName);
    setReferenceField(field);
    setSelectReferenceOpen(true);
  };

  const handleSelectReferenceClose = async (val: string) => {
    if (val != "") {
      let record: RecordModel = {
        ...recordData,
      };
      record[referenceField] = val;
      setRecordData(record);
    }
    setSelectReferenceOpen(false);
  };

  return (
    <form className="p-4 mx-auto absolute top-[3rem] justify-start text-black w-[40rem]">
      <h2 className="text-2xl font-bold mb-4">
        Add {CollectionConfig.name} Record
      </h2>

      {/* Input fields */}
      {CollectionConfig.columns.map((col, index) => (
        <div key={index} className="mb-4">
          {col.dataType == DataType.BOOLEAN && (
            <label className="inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={recordData[col.name]}
                className="sr-only peer"
                onChange={(e) => {
                  setRecordData({
                    ...recordData,
                    [col.name]: e.target.checked,
                  });
                }}
              />
              <div className="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-green-300 dark:peer-focus:ring-green-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-green-600"></div>
              <span className="ms-3 text-sm font-medium ">{col.name}</span>
            </label>
          )}
          {["created", "updated"].indexOf(col.name) == -1 &&
            (col.dataType == DataType.DOUBLE ||
              col.dataType == DataType.INTEGER ||
              col.dataType == DataType.DATETIME ||
              col.dataType == DataType.TEXT) && (
              <div className="rounded-md bg-slate-200 focus:bg-slate-300  p-2">
                <label className="block text-gray-700">{col.name}</label>
                <input
                  disabled={col.name == "id"}
                  type={
                    col.dataType == DataType.DOUBLE ||
                    col.dataType == DataType.INTEGER
                      ? "number"
                      : col.dataType == DataType.DATETIME
                      ? "datetime-local"
                      : "text"
                  }
                  name="record field"
                  className="border mb-2 rounded-md w-full bg-slate-200 focus:outline-none transition duration-300 "
                  placeholder="..."
                  value={
                    col.dataType == DataType.DATETIME
                      ? (recordData[col.name] as string)?.slice(0, 16)
                      : (recordData[col.name] as string)
                  }
                  onChange={(e) => {
                    setRecordData({
                      ...recordData,
                      [col.name]: e.target.value,
                    });
                  }}
                />
              </div>
            )}
          {col.dataType == DataType.REFERENCE && (
            <div className="rounded-md bg-slate-200 focus:bg-slate-300 p-2">
              <label className="block text-gray-700">{col.name}</label>
              <input
                value={recordData[col.name] as string}
                name="record field"
                className="border mb-2 rounded-md w-full bg-slate-200 focus:outline-none transition duration-300 "
                placeholder="..."
                onClick={() =>
                  handleSelectReferenceOpen(col.refTable, col.name)
                }
              />
            </div>
          )}
        </div>
      ))}

      {recordID != 0 && recordID != undefined && recordID != null && (
        <button
          type="button"
          onClick={handleDelete}
          className="bg-white border border-2 border-gray-950 p-2 mb-2 rounded-lg w-full hover:bg-red-600 transition"
        >
          Delete
        </button>
      )}
      <button
        type="button"
        onClick={handleSubmit}
        className="bg-white border border-2 border-gray-950 p-2 rounded-lg w-full"
      >
        Save
      </button>

      {/* Select Record Reference */}
      <Dialog
        isOpen={isSelectReferenceOpen}
        onClose={() => handleSelectReferenceClose("")}
      >
        <div className="w-[50rem] h-[30rem]">
          {isSelectReferenceOpen && referenceCollection != "" && (
            <SelectRecordReference
              collectionName={referenceCollection}
              onClose={(s: string) => handleSelectReferenceClose(s)}
            />
          )}
        </div>
      </Dialog>
    </form>
  );
}

export default AddEditRecord;
