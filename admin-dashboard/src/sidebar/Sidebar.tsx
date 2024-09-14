import collectionIcon from "../assets/collection.svg";
import logsIcon from "../assets/logs.svg";
import settingsIcon from "../assets/settings.svg";

const Sidebar = ({
  SelectedPage,
  SetSelectedPage,
}: {
  SelectedPage: string;
  SetSelectedPage: Function;
}) => {
  const icons = [
    { src: collectionIcon, alt: "collections", hint: "Collections" },
    { src: logsIcon, alt: "logs", hint: "Logs" },
    { src: settingsIcon, alt: "settings", hint: "Settings" },
  ];
  return (
    <div className="min-w-16 bg-white border-r border-gray-300 h-screen flex flex-col items-center py-4">
      {icons.map((icon, index) => (
        <div
          key={index}
          className="group relative my-2"
          onClick={() => SetSelectedPage(icon.alt)}
        >
          <img
            src={icon.src}
            alt={icon.alt}
            className={`w-10 h-10 mb-4 p-1 ${
              SelectedPage == icon.alt ? "border-2" : ""
            } border-solid rounded-lg border-gray-950 
             ${SelectedPage != icon.alt ? "hover:bg-gray-200" : ""}`}
          />
          <span className="absolute left-12 top-1/2 transform -translate-y-1/2 bg-gray-800 text-white text-xs rounded-md px-2 py-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
            {icon.hint}
          </span>
        </div>
      ))}
    </div>
  );
};

export default Sidebar;
