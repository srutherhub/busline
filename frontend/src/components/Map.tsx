import { useDarkMode } from "../hooks/useDarkMode";
import { Map as MapLibre } from "react-map-gl/maplibre";
import { DeckGL } from "@deck.gl/react";
import "maplibre-gl/dist/maplibre-gl.css";

const LIGHT = "https://basemaps.cartocdn.com/gl/positron-gl-style/style.json";
const DARK = "https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json";

export default function MapView() {
  const isDarkMode = useDarkMode();
  return (
    <DeckGL
      initialViewState={{
        longitude: -89.384171,
        latitude: 43.074699,
        zoom: 13,
        pitch: 45,
        bearing: 0,
      }}
      controller
      layers={[]}
    >
      <MapLibre mapStyle={isDarkMode ? DARK : LIGHT} />
    </DeckGL>
  );
}
