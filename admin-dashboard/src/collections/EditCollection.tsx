import { useEffect, useState } from "react";
import {
  ColumnConfig,
  DataType,
  ResponseMessage,
  TableConfig,
} from "../controllers/Types";
import { CollectionController } from "../controllers/CollectionController";
import deleteIcon from "../assets/delete.svg";

function EditCollection({
  CollectionConfig,
  onClose,
  reloadCollctionList,
}: {
  CollectionConfig: TableConfig;
  onClose: () => void;
  reloadCollctionList: () => void;
}) {
  const [openSelectDataType, setOpenSelectDataType] = useState<boolean>(false);
  const [fieldsToAdd, setFieldsToAdd] = useState<ColumnConfig[]>([]);
  const [fieldsToDelte, setFieldsToDelete] = useState<string[]>([]);
  const [presentFields, setPresent] = useState<ColumnConfig[]>([]);
  const [avaliableCollections, setAvaliableCollections] = useState<string[]>(
    []
  );

  const loadCollections = async () => {
    const response = await CollectionController.collectionsList();
    setAvaliableCollections(
      response.data
        .filter((col: TableConfig) => col.name != CollectionConfig.name)
        .map((col: TableConfig) => col.name)
    );
  };

  const handleAddField = (type: DataType) => {
    const newField: ColumnConfig = {
      name: "",
      dataType: type,
      refTable: type == DataType.REFERENCE ? avaliableCollections[0] : "",
      notNull: false,
      unique: false,
    };
    setFieldsToAdd([...fieldsToAdd, newField]);
    setOpenSelectDataType(false);
  };

  useEffect(() => {
    setPresent(CollectionConfig.columns);
    loadCollections();
  }, []);

  const handleSaveCollectionChanges = async () => {
    let addedNames = fieldsToAdd.map((field) => field.name);

    for (let index = 0; index < fieldsToAdd.length; index++) {
      const element = presentFields[index];
      if (addedNames.includes(element.name)) {
        alert(`Field already exists: ${element.name}`);
        return;
      }
    }

    for (let index = 0; index < fieldsToDelte.length; index++) {
      const element: ColumnConfig = {
        name: fieldsToDelte[index],
        refTable: "",
        dataType: DataType.TEXT,
        notNull: false,
        unique: false,
      };
      let tabConf: TableConfig = {
        name: CollectionConfig.name,
        columns: [element],
      };
      await CollectionController.deleteColumn(tabConf);
    }

    for (let index = 0; index < fieldsToAdd.length; index++) {
      const element = fieldsToAdd[index];
      let tabConf: TableConfig = {
        name: CollectionConfig.name,
        columns: [element],
      };
      let result: ResponseMessage = await CollectionController.addColumn(
        tabConf
      );
      if (result.success == false) {
        alert(result.message);
        return;
      }
    }

    onClose();
  };

  const handleAddFiledToBeDelete = (name: string) => {
    if (
      presentFields.map((col) => col.name).includes(name) &&
      !fieldsToDelte.includes(name)
    ) {
      setFieldsToDelete([...fieldsToDelte, name]);
    } else if (fieldsToDelte.includes(name)) {
      setFieldsToDelete(fieldsToDelte.filter((field) => field != name));
    } else if (fieldsToAdd.map((col) => col.name).includes(name)) {
      setFieldsToAdd(fieldsToAdd.filter((field) => field.name != name));
    }
  };

  const handleDeleteCollection = async () => {
    let result = confirm(
      `Do you want to delete ${CollectionConfig.name} collection?`
    );

    if (result) {
      await CollectionController.deleteCollection(CollectionConfig);
      onClose();
      reloadCollctionList();
    }
  };

  const handleEditField = (
    event: string,
    index: number,
    refTable: boolean = false
  ) => {
    let newFields = [...fieldsToAdd];
    if (refTable) newFields[index].refTable = event;
    else newFields[index].name = event;
    setFieldsToAdd(newFields);
  };

  const handleEditFieldFlag = (value: boolean, index: number, flag: string) => {
    let newFields = [...fieldsToAdd];
    if (flag == "UNIQUE") newFields[index].unique = value;
    if (flag == "NOTNULL") newFields[index].notNull = value;
    setFieldsToAdd(newFields);
  };

  return (
    <div
      className="p-4 mx-auto absolute top-[2rem] justify-start text-black  w-[40rem]"
      //onSubmit={() => {}}
    >
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold mb-4">Edit collection</h2>
        <div className="mb-4 p-2 rounded-md bg-slate-200 focus:bg-slate-300">
          <label className="block text-gray-700 relative">
            Collection Name
            <img
              src={deleteIcon}
              alt="delete"
              className={`hover:cursor-pointer p-1 rounded-lg hover:border-2 hover:bg-slate-200 border-gray-950 absolute transition w-8 top-0 right-0`}
              onClick={() => handleDeleteCollection()}
            />
          </label>
          <input
            disabled
            type="text"
            name="collectionName"
            value={CollectionConfig.name}
            className="border  py-2 rounded-md w-full  bg-slate-200 focus:outline-none transition duration-300 "
            placeholder="Enter collection name"
            //onChange={(event) => () => {}}
          />
        </div>
        <div className="mb-4">
          <p>
            Collection fields that cannot be deleted or updated{" "}
            <a className="bg-slate-200 rounded-md px-1">id</a>,
            <a className="bg-slate-200 rounded-md px-1">created</a>,
            <a className="bg-slate-200 rounded-md px-1">updated</a>
          </p>
        </div>
      </div>

      {/* Table fields */}
      {presentFields.map(
        (field, index) =>
          field.name !== "id" &&
          field.name !== "created" &&
          field.name !== "updated" && (
            <div
              className="mb-4 p-2 rounded-md bg-slate-200 focus:bg-slate-300"
              key={index + field.name}
            >
              <label className="block text-gray-700 relative">
                {field.dataType} {field.unique ? "UNIQUE " : ""}
                {field.notNull ? "NOT NULL" : ""}
                <img
                  src={deleteIcon}
                  alt="delete"
                  className={`${
                    fieldsToDelte.includes(field.name) ? "border-2" : ""
                  }
                    hover:cursor-pointer p-1 rounded-lg hover:p-1 hover:bg-slate-200 border-gray-950 absolute transition w-8 top-0 right-0`}
                  onClick={() => handleAddFiledToBeDelete(field.name)}
                />
              </label>
              <div className="flex">
                <input
                  disabled
                  type="text"
                  name="collectionName"
                  value={field.name}
                  className="border  py-2 rounded-md w-full bg-slate-200 focus:outline-none transition duration-300 "
                />
                {field.dataType === DataType.REFERENCE && (
                  <input
                    disabled
                    type="text"
                    name="collectionName"
                    value={field.refTable}
                    className="border  py-2 rounded-md w-full bg-slate-200 focus:outline-none transition duration-300 "
                  />
                )}
              </div>
            </div>
          )
      )}

      {/* Added fields */}
      {fieldsToAdd.map((field, index) => (
        <div className="mb-4 p-2 rounded-md bg-slate-200 focus:bg-slate-300">
          <label className="block text-gray-700 relative">
            <div className="flex">
              {field.dataType}
              <div
                onClick={() =>
                  handleEditFieldFlag(!field.unique, index, "UNIQUE")
                }
                className={`ml-2 px-1 ${
                  field.unique ? "bg-slate-400 rounded-md text-slate-200" : ""
                }`}
              >
                UNIQUE
              </div>
              <div
                onClick={() =>
                  handleEditFieldFlag(!field.notNull, index, "NOTNULL")
                }
                className={`ml-2 px-1 ${
                  field.notNull ? "bg-slate-400 rounded-md text-slate-200" : ""
                }`}
              >
                NOT NULL
              </div>
              <img
                src={deleteIcon}
                alt="delete"
                className="hover:cursor-pointer hover:border-2 p-1 hover:rounded-lg  hover:p-1 hover:bg-slate-200 hover:border-gray-950 absolute transition w-8 top-0 right-0"
                onClick={() => handleAddFiledToBeDelete(field.name)}
              />
            </div>
          </label>
          <div className="flex">
            <input
              type="text"
              name="collectionName"
              required
              value={field.name}
              className="border  py-2 rounded-md w-full  bg-slate-200 focus:outline-none transition duration-300 "
              placeholder="Column name"
              onChange={(event) => handleEditField(event.target.value, index)}
            />
            {field.dataType === DataType.REFERENCE && (
              <select
                className="w-full bg-slate-200 py-2 rounded-md "
                onChange={(event) =>
                  handleEditField(event.target.value, index, true)
                }
              >
                {avaliableCollections.map((col, idx) => (
                  <option key={idx}>{col}</option>
                ))}
              </select>
            )}
          </div>
        </div>
      ))}

      {/* Add field */}
      <div className="mb-4">
        <button
          type="button"
          className="bg-white  border border-2 border-gray-950 p-2 hover:bg-slate-200 transition duration-300 rounded-lg w-full"
          onClick={() => setOpenSelectDataType(!openSelectDataType)}
        >
          Add field ⬆️
        </button>
        {openSelectDataType && (
          <div className="grid  grid-cols-4 gap-2 border border-2 mt-2 rounded-md p-2 border-gray-950 ">
            {Object.values(DataType).map((type) => (
              <div
                key={type}
                className="cursor-pointer w-full p-2 hover:bg-slate-200 transition duration-300 rounded-md "
                onClick={() => {
                  handleAddField(type);
                }}
              >
                {type}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="flex space-x-4">
        <button
          type="button"
          className="bg-white border border-2 border-gray-950 p-2 hover:bg-gray-950 hover:text-white transition duration-300 rounded-lg w-full"
          onClick={onClose}
        >
          Cancel
        </button>
        <button
          type="submit"
          className="bg-white border border-2 border-gray-950 p-2 hover:bg-gray-950 hover:text-white transition duration-300 rounded-lg w-full"
          onClick={handleSaveCollectionChanges}
        >
          Save
        </button>
      </div>
    </div>
  );
}

export default EditCollection;
