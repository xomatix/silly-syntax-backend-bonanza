import { useEffect, useState } from "react";
import Drawer from "../components/Drawer";

import AddCollection from "./AddCollection";
import { CollectionController } from "../controllers/CollectionController";
import { TableConfig } from "../controllers/Types";
import CollectionView from "./CollectionView";

function CollectionsWrapper() {
  const [selectedCollection, setSelectedCollection] = useState<TableConfig>();
  const [isDrawerOpen, setDrawerOpen] = useState(false);
  const [collections, setCollections] = useState<TableConfig[]>([]);

  const downloadCollections = async () => {
    const response = await CollectionController.collectionsList();
    await setCollections(response.data);
    if (response.data.length > 0 && selectedCollection == null)
      await setSelectedCollection(response.data[0]);
  };

  useEffect(() => {
    downloadCollections();
  }, []);

  const toggleDrawer = async () => {
    setDrawerOpen(!isDrawerOpen);
    downloadCollections();
  };

  return (
    <div className="flex">
      {/* Collections sidebar */}
      <div className="min-w-56 bg-white text-black border-r border-gray-300 h-screen flex flex-col items-center py-4">
        <input
          type="text"
          className="border-b-2 mb-2 border-gray-300 p-2 bg-white focus:outline-none"
        />
        {collections.map((col, index) => (
          <div
            onClick={() => setSelectedCollection(col)}
            key={index}
            className={`group relative cursor-pointer mb-2 hover:bg-slate-200 w-52 rounded-md p-2 transition duration-300 text-center
               ${selectedCollection === col ? "bg-slate-200" : ""}`}
          >
            {col.name}
          </div>
        ))}
        <button
          onClick={toggleDrawer}
          className=" bg-white border border-2 border-gray-950 p-2 rounded-lg w-52 mx--auto"
        >
          Add Collection
        </button>

        <Drawer isOpen={isDrawerOpen} onClose={toggleDrawer}>
          <div className="w-[32rem]">
            {isDrawerOpen && (
              <AddCollection
                onClose={toggleDrawer}
                refresh={downloadCollections}
              />
            )}
          </div>
        </Drawer>
      </div>

      <div className="flex-1 ">
        {selectedCollection ? (
          <CollectionView
            collectionConfig={selectedCollection}
            reloadCollctionList={downloadCollections}
          />
        ) : (
          <h1>Select or create a collection</h1>
        )}
      </div>
    </div>
  );
}

export default CollectionsWrapper;
