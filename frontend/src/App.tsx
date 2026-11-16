import "./App.css";
import Map from "./components/Map";
import Navbar from "./components/Navbar";

function App() {
  return (
    <div className="layout">
      <div className="sidebar-container">
        <Navbar />
      </div>
      <div className="map-container">
        <Map />
      </div>
    </div>
  );
}

export default App;
