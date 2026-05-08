import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthGuard } from "./auth-guard";
import { Button } from "@/components/ui/button";

// A dummy Dashboard
const Dashboard = () => (
  <div className="p-8">
    <h1 className="text-2xl font-bold italic text-primary">NudgeBuddy Dashboard</h1>
    <Button className="mt-4">Start Outreach</Button>
  </div>
);

// A dummy Login Page
const LoginPage = () => (
  <div className="flex h-screen items-center justify-center">
    <Button variant="outline">Login with Google</Button>
  </div>
);

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />

        {/* Protect this route */}
        <Route
          path="/"
          element={
            <AuthGuard>
              <Dashboard />
            </AuthGuard>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}