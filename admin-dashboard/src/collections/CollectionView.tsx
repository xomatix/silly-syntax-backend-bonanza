import { useEffect, useState } from "react";
import { DataListModel, DataType, TableConfig } from "../controllers/Types";
import settingsIcon from "../assets/gear.svg";
import refreshIcon from "../assets/refresh.svg";
import addIcon from "../assets/add.svg";
import EditCollection from "./EditCollection";
import Drawer from "../components/Drawer";
import { CollectionController } from "../controllers/CollectionController";
import { RecordController } from "../controllers/RecordController";
import AddEditRecord from "../records/AddEditRecord";

const CollectionView = ({
  collectionConfig,
  reloadCollctionList: reloadCollectionList,
}: {
  collectionConfig: TableConfig;
  reloadCollctionList: () => void;
}) => {
  const [collectionData, setCollectionData] = useState<TableConfig>();
  const [isColEditDrawerOpen, setColEditDrawerOpen] = useState(false);
  const [isRecordAddEditDrawerOpen, setRecordAddEditDrawerOpen] =
    useState(false);

  const [records, setRecords] = useState<object[]>();
  const [selectedRecordID, setSelectRecordID] = useState<number>(0);
  const [filter, setFilter] = useState<string>("");

  useEffect(() => {
    setRecords([]);
    setCollectionData(collectionConfig);
    loadRecords();
  }, [collectionConfig]);

  //#region Handlers
  const loadRecords = async () => {
    let model: DataListModel = {
      collectionName: collectionConfig.name,
      filter: filter,
    };
    const response = await RecordController.GetRecords(model);
    let objs: any[] = [];
    if (response.data != null) {
      objs = [...response.data];
    }
    setRecords(objs);
  };

  const reloadCollectionConfig = async () => {
    const response = await CollectionController.collectionsList();
    let respArr: TableConfig[] = response.data;
    let found = respArr.find((col) => col.name == collectionConfig.name);
    if (found != undefined) setCollectionData(found);
  };

  const onClose = () => {
    setColEditDrawerOpen(false);
    reloadCollectionConfig();
  };

  const handleOpenEdit = (id: number) => {
    setSelectRecordID(id);
    setRecordAddEditDrawerOpen(true);
  };

  const onCloseAddEditRecord = async () => {
    setRecordAddEditDrawerOpen(false);
    setSelectRecordID(0);
    loadRecords();
  };

  const handleChangeFilter = (e: any) => {
    setFilter(e.target.value);
  };

  const handleFilterApply = (event: any) => {
    if (event.key === "Enter") {
      loadRecords();
    }
  };

  //#endregion

  return (
    <div className="container bg-white min-w-full">
      {/* Header */}
      <div className="flex px-4 pt-4 justify-left items-center mb-4 space-x-2 ">
        <h1 className="text-2xl text-black font-semibold">
          Collection / {collectionConfig.name}
        </h1>
        <div className="flex space-x-2 ">
          <button
            onClick={() => setColEditDrawerOpen(true)}
            className=" text-white p-2 my-auto w-9 h-9 rounded hover:rotate-90 transition duration-300"
          >
            <img src={settingsIcon} alt="settings" />
          </button>
          <button
            onClick={() => alert("Settings button clicked")}
            className=" text-white p-2 my-auto w-9 h-9 rounded hover:rotate-180 transition duration-300"
          >
            <img src={refreshIcon} alt="refresh" />
          </button>
          <button
            onClick={() => setRecordAddEditDrawerOpen(true)}
            className=" text-white p-2 my-auto w-9 h-9 rounded hover:bg-slate-200 transition duration-300"
          >
            <img src={addIcon} alt="insert" />
          </button>
        </div>
      </div>

      {/* Border */}
      <div className="border-b border-gray-300 mb-2"></div>

      {/* filter */}
      <div className="rounded-md bg-slate-200 focus:bg-slate-300 m-2 p-2 h-full">
        <input
          type="text"
          name="filter"
          className="border text-black rounded-md w-full bg-slate-200 focus:outline-none transition duration-300 "
          placeholder="..."
          value={filter}
          onChange={handleChangeFilter}
          onKeyUp={handleFilterApply}
        />
      </div>

      {/* Border */}
      <div className="border-b border-gray-300 mb-2"></div>

      {/* Table */}
      <div className="relative overflow-x-auto text-sm">
        <table className="w-full bg-white">
          <thead className="bg-white">
            <tr>
              {collectionData?.columns.map((column, index) => (
                <th
                  key={index + column.name}
                  className="py-2 px-4 text-black font-semibold border-b border-gray-200 text-left min-w-[13rem]"
                >
                  {column.name == "id" ? "🔑" : ""}
                  {column.dataType == DataType.DATETIME ? "🕒" : ""}
                  {column.dataType == DataType.TEXT && column.name != "id"
                    ? "📝"
                    : ""}
                  {column.name != "id" &&
                  (column.dataType == DataType.DOUBLE ||
                    column.dataType == DataType.INTEGER)
                    ? "🔢"
                    : ""}
                  {column.dataType == DataType.BOOLEAN ? "✅" : ""}
                  &nbsp;
                  {column.name}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {records?.map((item, index) => (
              <tr
                key={index}
                onClick={() =>
                  handleOpenEdit(item["id" as keyof typeof item] as number)
                }
                className="hover:bg-slate-200 transition ease-in-out duration-300 "
              >
                {Object.values(item).map((_, i) => (
                  <td
                    key={i}
                    className="py-2 px-4 text-black font-semibold border-b border-gray-200 text-left max-w-[12rem] overflow-x-hidden"
                  >
                    {collectionData?.columns[i]?.dataType != DataType.BOOLEAN
                      ? item[
                          collectionData?.columns[i]?.name as keyof typeof item
                        ]
                      : item[
                          collectionData?.columns[i]?.name as keyof typeof item
                        ] == true
                      ? "✅"
                      : "❌"}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Edit collection drawer */}
      <Drawer isOpen={isColEditDrawerOpen} onClose={onClose}>
        <div className="w-[40rem]">
          {isColEditDrawerOpen && collectionData != undefined && (
            <EditCollection
              CollectionConfig={collectionData}
              onClose={onClose}
              reloadCollctionList={reloadCollectionList}
            />
          )}
        </div>
      </Drawer>

      {/* Add record drawer */}
      <Drawer isOpen={isRecordAddEditDrawerOpen} onClose={onCloseAddEditRecord}>
        <div className="w-[40rem]">
          {isRecordAddEditDrawerOpen && collectionData != undefined && (
            <AddEditRecord
              CollectionConfig={collectionData}
              onClose={() => {
                onCloseAddEditRecord();
              }}
              recordID={selectedRecordID}
            />
          )}
        </div>
      </Drawer>
    </div>
  );
};

export default CollectionView;
