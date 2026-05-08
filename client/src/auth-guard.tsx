import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

interface AuthGuardProps {
    children: React.ReactNode;
}

export function AuthGuard({ children }: AuthGuardProps) {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        const checkAuth = async () => {
            try {
                // Replace this URL with your actual API endpoint
                const response = await fetch("/api/login", {
                    method: "GET", // Or POST depending on your backend setup
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                if (response.ok) {
                    setIsAuthenticated(true);
                } else {
                    setIsAuthenticated(false);
                    navigate("/login"); // Redirect to login page if unauthorized
                }
            } catch (error) {
                console.error("Auth check failed:", error);
                setIsAuthenticated(false);
                navigate("/login");
            }
        };

        checkAuth();
    }, [navigate]);

    // While checking, show a loading state (or a nice shadcn spinner!)
    if (isAuthenticated === null) {
        return (
            <div className="flex h-screen w-screen items-center justify-center">
                <p className="text-sm text-muted-foreground animate-pulse">
                    Authenticating...
                </p>
            </div>
        );
    }

    return isAuthenticated ? <>{children}</> : null;
}