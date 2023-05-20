import { useState } from "react";
import CodeEditorWindow from "./components/CodeEditorWindow";
import LanguagesDropdown from "./components/LanguagesDropdown";
import { languageOptions } from "./constants/languageOptions";
import { defineTheme } from "./lib/defineTheme";
import ThemesDropdown from "./components/ThemesDropdown";
import axios from "axios";

const App = () => {
  const [code, setCode] = useState("");
  const [language, setLanguage] = useState(languageOptions[0]);
  const [theme, setTheme] = useState("cobalt");
  const [processing, setProcessing] = useState(false);
  const [customInput, setCustomInput] = useState("");
  const [outputDetails, setOutputDetails] = useState(null);

  const onChange = (data) => {
    setCode(data);
  };

  const handleCompile = () => {
    setProcessing(true);
    const formData = {
      language_id: language.id,
      source_code: window.atob(code),
      stdin: window.atob(customInput),
    };
    const options = {
      method: "POST",
      url: import.meta.env.VITE_RAPID_API_URL,
      params: { base64_encoded: "true", fields: "*" },
      headers: {
        "content-type": "application/json",
        "Content-Type": "application/json",
        "X-RapidAPI-Host": import.meta.env.VITE_RAPID_API_HOST,
        "X-RapidAPI-Key": import.meta.env.VITE_RAPID_API_KEY,
      },
      data: formData,
    };

    axios
      .request(options)
      .then(function (response) {
        const token = response.data.token;
        checkStatus(token);
      })
      .catch((err) => {
        let error = err.response ? err.response.data : err;
        setProcessing(false);
      });
  };

  const checkStatus = async (token) => {
    const options = {
      method: "GET",
      url: import.meta.env.VITE_RAPID_API_URL + "/" + token,
      params: { base64_encoded: "true", fields: "*" },
      headers: {
        "X-RapidAPI-Host": import.meta.env.VITE_RAPID_API_HOST,
        "X-RapidAPI-Key": import.meta.env.VITE_RAPID_API_KEY,
      },
    };
    try {
      let response = await axios.request(options);
      let statusId = response.data.status?.id;

      if (statusId === 1 || statusId === 2) {
        setTimeout(() => {
          checkStatus(token);
        }, 2000);
        return;
      } else {
        setProcessing(false);
        setOutputDetails(response.data);
        return;
      }
    } catch (err) {
      setProcessing(false);
    }
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
      <ThemesDropdown />
      <CodeEditorWindow />
    </div>
  );
};

export default App;
