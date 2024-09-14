import { useEffect, useState } from "react";
import "./App.css";
import Sidebar from "./sidebar/Sidebar";
import CollectionsWrapper from "./collections/CollectionsWrapper";
import LoginComponent from "./authentication/LoginComponent";

function App() {
  const [selectedPage, setSelectedPage] = useState("collections");

  useEffect(() => {
    checkLoginStatus();
  }, []);

  const checkLoginStatus = async () => {
    const token = localStorage.getItem("bonanza_token");
    if (token == null || token == "") setSelectedPage("login");
    else setSelectedPage("collections");
  };

  return selectedPage == "login" ? (
    <LoginComponent onClose={checkLoginStatus} />
  ) : (
    <>
      <div className="flex">
        <Sidebar
          SelectedPage={selectedPage}
          SetSelectedPage={setSelectedPage}
        />
        <div className="flex-1 overflow-x-hidden">
          {selectedPage == "collections" && <CollectionsWrapper />}
          {selectedPage == "logs" && <div>Logs</div>}
          {selectedPage == "settings" && <div>Settings</div>}
        </div>
      </div>
    </>
  );
}

export default App;
