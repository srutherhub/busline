import "./App.css";
import Map from "./components/Map";

function App() {
  return (
    <div className="layout">
      <div className="sidebar-container">
        <div>hello from app</div>
      </div>
      <div className="map-container">
        <Map />
      </div>
    </div>
  );
}

export default App;
