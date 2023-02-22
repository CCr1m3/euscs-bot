import { useState } from "react";
import reactLogo from "./assets/react.svg";
import "./App.css";
import { useCookies } from "react-cookie";
import { redirect, useNavigate } from "react-router-dom";

function App() {
  const [count, setCount] = useState(0);
  const callApi = () => {
    fetch(window.location.origin + "/api")
      .then((res) => res.json())
      .then((val) => alert(val.message));
  };
  const [cookies, setCookie, removeCookie] = useCookies(["sessionid"]);
  if (!cookies.sessionid) {
    window.location.href = window.location.origin + "/auth";
    return null;
  }
  return (
    <div className="App">
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src="/vite.svg" className="logo" alt="Vite logo" />
        </a>
        <a href="https://reactjs.org" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <button onClick={() => callApi()}>Click to call api</button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  );
}

export default App;
