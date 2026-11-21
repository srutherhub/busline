import { queryOptions } from "@tanstack/react-query";
import { gtfsApi } from "../models/GtfsApi";
const HOUR = 1000 * 60 * 60;
export const gtfsQueries = {
  feedInfo: (city: string) =>
    queryOptions({
      queryKey: ["feed", city],
      queryFn: () => gtfsApi.getCityFeedInfo(city),
      staleTime: HOUR,
      retry: false,
    }),

  routes: (city: string) =>
    queryOptions({
      queryKey: ["routes", city],
      queryFn: () => gtfsApi.getActiveRoutes(city),
      retry: false,
    }),

  nearbyStops: (city: string, lat: number, lon: number) =>
    queryOptions({
      queryKey: ["nearbyStops", city, lat, lon],
      queryFn: () => gtfsApi.getNearbyStops(city, lat, lon),
      staleTime: HOUR,
      retry: false,
    }),

  tripsForStop: (city: string, stopId: string) =>
    queryOptions({
      queryKey: ["tripsForStop", city, stopId],
      queryFn: () => gtfsApi.getTripsForStop(city, stopId),
      staleTime: HOUR,
      retry: false,
    }),

  shapes: (city: string, shapeId: string) =>
    queryOptions({
      queryKey: ["shapes", city, shapeId],
      queryFn: () => gtfsApi.getShapesForTrip(city, shapeId),
      staleTime: HOUR,
      retry: false,
    }),
};
