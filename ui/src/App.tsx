import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { AuthLayout } from "./layout/AuthLayout";
import NavLayout from "./layout/NavLayout";

function App() {
  return (
    <>
      <AuthLayout>
        <NavLayout>
          <Outlet />
        </NavLayout>
      </AuthLayout>
      <TanStackRouterDevtools position="bottom-right" />
      <ReactQueryDevtools position="bottom" />
    </>
  );
}

export default App;
