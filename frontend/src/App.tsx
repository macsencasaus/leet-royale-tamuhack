import { useState } from "react";
import Editor from "@monaco-editor/react";

function App() {
  const [_, setWs] = useState<WebSocket | undefined>(undefined);

  function connect() {
    const temp = new WebSocket(`ws://${window.location.host}/ws`);
    temp.onopen = () => console.log("Open");
    temp.onmessage = onMessage;
    setWs(temp);
  }

  function onMessage(messageJSON: MessageEvent) {
    console.log(messageJSON.data);
  }

  return (
    <>
      <button onClick={connect}>Connect</button>
      <Editor
        height="90vh"
        width="90dvw"
        defaultLanguage="javascript"
        theme="vs-dark"
        defaultValue="// some comment"
      />
    </>
  );
}

export default App;
