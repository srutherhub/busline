import { queryOptions } from "@tanstack/react-query";
import { gtfsApi } from "../models/GtfsApi";

export const gtfsQueries = {
  feedInfo: (city: string) =>
    queryOptions({
      queryKey: ["feed", city],
      queryFn: () => gtfsApi.getCityFeedInfo(city),
    }),

  routes: (city: string) =>
    queryOptions({
      queryKey: ["routes", city],
      queryFn: () => gtfsApi.getActiveRoutes(city),
    }),

  nearbyStops: (city: string, lat: number, lon: number) =>
    queryOptions({
      queryKey: ["nearbyStops", city, lat, lon],
      queryFn: () => gtfsApi.getNearbyStops(city, lat, lon),
    }),

  tripsForStop: (city: string, stopId: string) =>
    queryOptions({
      queryKey: ["tripsForStop", city, stopId],
      queryFn: () => gtfsApi.getTripsForStop(city, stopId),
    }),

  shapes: (city: string, shapeId: string) =>
    queryOptions({
      queryKey: ["shapes", city, shapeId],
      queryFn: () => gtfsApi.getShapesForTrip(city, shapeId),
    }),
};
