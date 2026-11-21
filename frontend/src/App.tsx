import "./App.css";
import Map from "./components/Map";
import Navbar from "./components/Navbar";

function App() {
  return (
    <div className="app">
      <Navbar />
      <div className="layout">
        <div className="sidebar-container"></div>
        <div className="map-container">
          <Map />
        </div>
      </div>
    </div>
  );
}

export default App;
