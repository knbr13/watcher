import { useState } from "react";
import CodeEditorWindow from "./components/CodeEditorWindow";

const App = () => {
  const [code, setCode] = useState("");

  const onChange = (data) => {
    setCode(data);
  };

  return (
    <div>
      Hello Code!
      <CodeEditorWindow />
    </div>
  );
};

export default App;
