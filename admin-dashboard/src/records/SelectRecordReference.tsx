import { useEffect, useState } from "react";
import { RecordController } from "../controllers/RecordController";
import { DataListModel, TableConfig } from "../controllers/Types";
import { CollectionController } from "../controllers/CollectionController";

function SelectRecordReference({
  collectionName,
  onClose,
}: {
  collectionName: string;
  onClose: (s: string) => void;
}) {
  const [collectionData, setCollectionData] = useState<TableConfig>();
  const [records, setRecords] = useState<object[]>([]);
  const [filter, setFilter] = useState<string>("");

  useEffect(() => {
    handleCollectionDataLoad();
    handleRefresh();
  }, [collectionName]);

  const handleChangeFilter = (e: any) => {
    setFilter(e.target.value);
  };

  const handleRefresh = async () => {
    let model: DataListModel = {
      collectionName: collectionName,
    };
    if (filter != "") {
      model.filter = filter;
    }

    const response = await RecordController.GetRecords(model);
    if (response.data.length > 0) {
      setRecords(response.data);
    } else {
      setRecords([]);
    }
  };

  const handleCollectionDataLoad = async () => {
    const response = await CollectionController.collectionsList();
    let TableConfigs: TableConfig[] = response.data;

    let found = TableConfigs.find((col) => col.name == collectionName);
    if (!found) return;

    setCollectionData(found);
  };

  const handleSelectSaveRefenrence = async (selectedId: string) => {
    onClose(selectedId);
  };

  const handleFilterApply = (event: any) => {
    if (event.key === "Enter") {
      handleRefresh();
    }
  };

  return (
    <div className="relative text-sm">
      {/* filter */}
      <div className="rounded-md bg-slate-200 focus:bg-slate-300 mb-2 p-2 h-full">
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

      {/* table with records*/}
      <table className="min-w-full bg-white text-sm text-left text-black font-light select-none">
        <thead className="bg-white">
          <tr>
            {collectionData?.columns.map(
              (column, i) =>
                column.name != "created" &&
                column.name != "updated" && (
                  <th
                    key={i}
                    className="py-2 px-4 border-b border-gray-200 text-left min-w-[13rem]"
                  >
                    {column.name}
                  </th>
                )
            )}
          </tr>
        </thead>
        <tbody>
          {records?.map((item, index) => (
            <tr
              key={index}
              className="hover:bg-slate-200 transition ease-in-out duration-300 border-b border-gray-200"
              onDoubleClick={() =>
                handleSelectSaveRefenrence(
                  item["id" as keyof typeof item] as string
                )
              }
            >
              {Object.values(item).map(
                (_, i) =>
                  collectionData?.columns[i]?.name != "created" &&
                  collectionData?.columns[i]?.name != "updated" && (
                    <td
                      key={i}
                      className="py-2 cursor-pointer px-4
                    min-w-[13rem] overflow-x-hidden"
                    >
                      <p
                        className={` ${
                          collectionData?.columns[i]?.name == "id"
                            ? "bg-slate-200 rounded-md max-w-fit p-1"
                            : ""
                        }`}
                      >
                        {
                          item[
                            collectionData?.columns[i]
                              ?.name as keyof typeof item
                          ]
                        }
                      </p>
                    </td>
                  )
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default SelectRecordReference;
