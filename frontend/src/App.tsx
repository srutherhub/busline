import "./App.css";
import Map from "./components/Map";
import { useQuery } from "@tanstack/react-query";
import { gtfsQueries } from "./queries/gtfsQueries";
import { useCity } from "./hooks/useCity";
import Sidebar from "./components/Sidebar";
import { SkeletonTheme } from "react-loading-skeleton";

function App() {
  const city = useCity();
  const { data: routes } = useQuery(gtfsQueries.routes(city));
  console.log(routes);
  return (
    <div className="layout">
      <div>
        <SkeletonTheme
          baseColor="var(--surface)"
          highlightColor="var(--border)"
        >
          <Sidebar city={city} />
        </SkeletonTheme>
      </div>
      <div className="map-container">
        <Map />
      </div>
    </div>
  );
}

export default App;
