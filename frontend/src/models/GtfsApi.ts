import { ApiClient } from "./ApiClient";

export type FeedInfo = {
  id: string;
  publisher_name: string;
  url: string;
  language: string;
  start_date: string;
  end_date: string;
  imported_at: string;
  version: string;
  city: string;
};

export type Route = {
  id: string;
  route_id: string;
  short_name: string;
  long_name: string;
  type: number;
  url: string;
  color: string;
  text_color: string;
};

export type Stop = {
  id: string;
  stop_id: string;
  code: string;
  name: string;
  latitude: number;
  longitude: number;
  wheelchair_boarding: number;
  distance_to: number;
};

export type StopTrip = {
  id: string;
  trip_id: string;
  service_id: string;
  headsign: string;
  shape_id: string;
  wheelchair_accessible: number;
  route_id: string;
  route_short_name: string;
  route_long_name: string;
  route_color: string;
  route_text_color: string;
  route_type: number;
  arrival_time: string;
  departure_time: string;
  stop_sequence: number;
  shape_dist_traveled: number;
};

export type Shape = {
  id: string;
  latitude: string;
  longitude: string;
  sequence: string;
  distance_travelled: string;
};

export class GtfsApi {
  apiClient = new ApiClient();

  getCityFeedInfo(city: string): Promise<FeedInfo> {
    return this.apiClient.get(`/${city}/feed`);
  }

  getActiveRoutes(city: string): Promise<Route[]> {
    return this.apiClient.get(`/${city}/routes`);
  }

  getNearbyStops(city: string, lat: number, lon: number): Promise<Stop[]> {
    return this.apiClient.get(`/${city}/stops?lat=${lat}&long=${lon}`);
  }

  getTripsForStop(city: string, stopId: string): Promise<StopTrip[]> {
    return this.apiClient.get(`/${city}/stops/${stopId}`);
  }

  getShapesForTrip(city: string, shapeId: string): Promise<Shape[]> {
    return this.apiClient.get(`/${city}/shapes/${shapeId}`);
  }
}

export const gtfsApi = new GtfsApi();
