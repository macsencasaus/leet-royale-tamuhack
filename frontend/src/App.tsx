import { useState } from "react"


function App() {
  const [ws, setWs] = useState<WebSocket | undefined>(undefined);

  function connect() {
    const temp = new WebSocket("ws://10.246.176.24:6969/ws");
    temp.onopen = () => console.log("Open");
    temp.onmessage = onMessage;
    setWs(temp);
  }

  function onMessage(messageJSON: MessageEvent) {
    console.log(messageJSON.data);
  }

  return (
    <button onClick={connect}>
      Connect
    </button>
  )
}

export default App
