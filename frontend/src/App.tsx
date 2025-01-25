import { useState } from "react"


function App() {
  const [ws, setWs] = useState<WebSocket | undefined>(undefined);

  function connect() {
    
  }

  return (
    <button onClick={connect}>
      Connect
    </button>
  )
}

export default App
