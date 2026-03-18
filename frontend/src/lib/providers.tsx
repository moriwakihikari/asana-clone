"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState, createContext, useContext, useEffect, useCallback } from "react";
import type { User, Workspace } from "@/types";
import { authApi, userApi, workspaceApi } from "./api";

interface AuthContextType {
  user: User | null;
  workspace: Workspace | null;
  workspaces: Workspace[];
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  selectWorkspace: (ws: Workspace) => void;
  refreshWorkspaces: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}

function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [workspace, setWorkspace] = useState<Workspace | null>(null);
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const refreshWorkspaces = useCallback(async () => {
    try {
      const wsList = await workspaceApi.list();
      setWorkspaces(wsList);
      if (wsList.length > 0 && !workspace) {
        setWorkspace(wsList[0]);
      }
    } catch {
      // ignore
    }
  }, [workspace]);

  useEffect(() => {
    const token = localStorage.getItem("access_token");
    if (token) {
      userApi
        .getMe()
        .then((u) => {
          setUser(u);
          return workspaceApi.list();
        })
        .then((wsList) => {
          setWorkspaces(wsList);
          if (wsList.length > 0) setWorkspace(wsList[0]);
        })
        .catch(() => {
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
        })
        .finally(() => setIsLoading(false));
    } else {
      setIsLoading(false);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const res = await authApi.login({ email, password });
    localStorage.setItem("access_token", res.tokens.access_token);
    localStorage.setItem("refresh_token", res.tokens.refresh_token);
    setUser(res.user);
    const wsList = await workspaceApi.list();
    setWorkspaces(wsList);
    if (wsList.length > 0) setWorkspace(wsList[0]);
  };

  const register = async (name: string, email: string, password: string) => {
    const res = await authApi.register({ name, email, password });
    localStorage.setItem("access_token", res.tokens.access_token);
    localStorage.setItem("refresh_token", res.tokens.refresh_token);
    setUser(res.user);
    const wsList = await workspaceApi.list();
    setWorkspaces(wsList);
    if (wsList.length > 0) setWorkspace(wsList[0]);
  };

  const logout = () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    setUser(null);
    setWorkspace(null);
    setWorkspaces([]);
  };

  const selectWorkspace = (ws: Workspace) => setWorkspace(ws);

  return (
    <AuthContext.Provider
      value={{ user, workspace, workspaces, isLoading, login, register, logout, selectWorkspace, refreshWorkspaces }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 30 * 1000,
            retry: 1,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>{children}</AuthProvider>
    </QueryClientProvider>
  );
}
