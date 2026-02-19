import React, { createContext, useContext, useEffect, useMemo, useState } from "react";
import { apiClient } from "../api/client";

export interface AuthUser {
  id: number;
  email?: string;
  username?: string;
}

interface AuthContextValue {
  user: AuthUser | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const refresh = async () => {
    try {
      const response = await apiClient.get("/auth/me");
      setUser({ id: response.data.id });
    } catch (error) {
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    refresh();
  }, []);

  const login = async (email: string, password: string) => {
    const response = await apiClient.post("/auth/login", { email, password });
    setUser({ id: response.data.id, email: response.data.email });
  };

  const register = async (email: string, username: string, password: string) => {
    await apiClient.post("/auth/register", { email, username, password });
    await login(email, password);
  };

  const logout = async () => {
    try {
      await apiClient.post("/auth/logout");
    } finally {
      setUser(null);
    }
  };

  const value = useMemo(
    () => ({ user, isLoading, login, register, logout, refresh }),
    [user, isLoading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
