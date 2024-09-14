import closeIcon from "../assets/close.svg";

const Drawer = ({ isOpen, onClose, children }: any) => {
  return (
    <div
      className={`fixed inset-0 z-50 flex ${
        isOpen ? "pointer-events-auto" : "pointer-events-none"
      }`}
    >
      {/* Overlay */}
      <div
        className={`fixed inset-0 bg-black transition-opacity duration-300 ${
          isOpen ? "opacity-50" : "opacity-0"
        }`}
        onClick={onClose}
      ></div>

      {/* Drawer */}
      <div
        className={`fixed overflow-y-scroll inset-y-0 right-0 w-fit bg-white shadow-lg transform transition-transform duration-300 ${
          isOpen ? "translate-x-0" : "translate-x-full"
        }`}
      >
        <div className="flex justify-end p-4 mb-2">
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-800  "
          >
            <img src={closeIcon} alt="" className="w-6 h-6 " />
          </button>
        </div>
        <div className="flex flex-col items-center justify-center h-full p-4">
          {children}
        </div>
      </div>
    </div>
  );
};

export default Drawer;
