import { createPortal } from "react-dom";

const Dialog = ({ isOpen, onClose, children }: any) => {
  return createPortal(
    <div
      className={`${
        isOpen ? "block" : "hidden"
      } fixed inset-0 z-50 flex items-center justify-center `}
    >
      <div className="bg-white rounded-lg z-51 shadow-lg w-md p-6">
        <div className=" overflow-y-auto overflow-x-auto">{children}</div>
        <div className="flex justify-end">
          <button
            onClick={onClose}
            className="mt-4 px-4 py-2 bg-white border border-2 border-gray-950 rounded-md text-black rounded hover:bg-white "
          >
            Close
          </button>
        </div>
      </div>
    </div>,
    document.body
  );
};

export default Dialog;
