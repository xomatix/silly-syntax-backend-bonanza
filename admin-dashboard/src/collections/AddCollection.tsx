import { useState } from "react";
import { CollectionController } from "../controllers/CollectionController";

function AddCollection({
  onClose,
  refresh,
}: {
  onClose: () => void;
  refresh: () => void;
}) {
  const [name, setName] = useState("");
  const handleSubmit = async (event: any) => {
    event.preventDefault();

    console.log("name", name);

    await CollectionController.addCollection(name);

    refresh();
    onClose();
  };

  return (
    <form
      className="p-4 mx-auto absolute top-[2rem] justify-start  w-[32rem]"
      onSubmit={handleSubmit}
    >
      <h2 className="text-2xl font-bold mb-4">New collection</h2>
      <div className="mb-4 p-2 rounded-md bg-slate-200 focus:bg-slate-300">
        <label className="block text-gray-700">
          Collection Name<i className="text-red-500">*</i>
        </label>
        <input
          required
          type="text"
          name="collectionName"
          className="border  py-2 rounded-md w-full  bg-slate-200 focus:outline-none transition duration-300 "
          placeholder="Enter collection name"
          onChange={(event) => setName(event.target.value)}
        />
      </div>
      <div className="mb-4">
        <p>
          Collection will have predefined fields{" "}
          <a className="bg-slate-200 rounded-md px-1">id</a>,
          <a className="bg-slate-200 rounded-md px-1">created</a>,
          <a className="bg-slate-200 rounded-md px-1">updated</a>
        </p>
      </div>
      <button
        type="submit"
        className="bg-white border border-2 border-gray-950 p-2 rounded-lg w-full"
      >
        Save
      </button>
    </form>
  );
}

export default AddCollection;
