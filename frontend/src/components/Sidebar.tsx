import { useQuery } from "@tanstack/react-query";
import Skeleton from "react-loading-skeleton";
import "react-loading-skeleton/dist/skeleton.css";
import { Link, Navigate } from "react-router-dom";
import { gtfsQueries } from "../queries/gtfsQueries";

type SidebarProps = {
  city: string;
};

const capitalize = (str: string) => str.charAt(0).toUpperCase() + str.slice(1);

export default function Sidebar({ city }: SidebarProps) {
  const feed = useQuery(gtfsQueries.feedInfo(city));
  const routes = useQuery(gtfsQueries.routes(city));

  if (feed.isError) {
    return <Navigate to="/error" replace />;
  }

  return (
    <div className="sidebar-container">
      <div className="feed-container">
        <h1>
          {feed.isPending ? (
            <Skeleton width="50%" />
          ) : (
            capitalize(feed.data.city)
          )}
        </h1>
        {feed.isPending ? (
          <Skeleton width="60%" />
        ) : (
          <Link
            to={feed.data.url}
            target="_blank"
            className="interactable secondary"
          >
            {feed.data.publisher_name}
          </Link>
        )}
      </div>
      {routes.isPending ? (
        <Skeleton className="routes-container" />
      ) : (
        <div className="routes-container">
          <h3>Routes</h3>
          <div className="routes secondary">
            {routes.isError ? (
              <p>Couldn't load routes.</p>
            ) : (
              routes.data.map((r) => (
                <Link to={`?route=${r.id}`} className="interactable secondary">
                  {r.long_name}
                </Link>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}
