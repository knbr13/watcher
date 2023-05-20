import { useState } from "react";
import CodeEditorWindow from "./components/CodeEditorWindow";
import LanguagesDropdown from "./components/LanguagesDropdown";
import { languageOptions } from "../constants/languageOptions";

const App = () => {
  const [code, setCode] = useState("");
  const [language, setLanguage] = useState(languageOptions[0]);

  const onChange = (data) => {
    setCode(data);
  };

  const onSelectChange = (sl) => {
    setLanguage(sl);
  };

  return (
    <div>
      Hello Code!
      <LanguagesDropdown />
      <CodeEditorWindow />
    </div>
  );
};

export default App;
