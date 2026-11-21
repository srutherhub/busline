import { useParams } from "react-router-dom";

export function useCity(): string {
  const { city } = useParams<{ city: string }>();
  if (!city) throw new Error("useCity must be used under a /:city route");
  return city;
}
