import { useState } from "react";
import CodeEditorWindow from "./components/CodeEditorWindow";
import LanguagesDropdown from "./components/LanguagesDropdown";
import { languageOptions } from "../constants/languageOptions";
import { defineTheme } from "./lib/defineTheme";

const App = () => {
  const [code, setCode] = useState("");
  const [language, setLanguage] = useState(languageOptions[0]);
  const [theme, setTheme] = useState("cobalt");

  const onChange = (data) => {
    setCode(data);
  };

  function handleThemeChange(th) {
    const theme = th;
    if (["light", "vs-dark"].includes(theme.value)) {
      setTheme(theme);
    } else {
      defineTheme(theme.value).then((_) => setTheme(theme));
    }
  }

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
