import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import NavLayout from "./layout/NavLayout";

function App() {
  return (
    <>
      <NavLayout>
        <Outlet />
      </NavLayout>
      <TanStackRouterDevtools position="bottom-right" />
      <ReactQueryDevtools position="bottom" />
    </>
  );
}

export default App;
