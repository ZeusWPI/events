import { Navigate } from "@tanstack/react-router";

function Index() {
  return <Navigate to="/events" mask={{ to: "/" }} />;
}

export default Index;
