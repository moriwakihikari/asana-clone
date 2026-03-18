"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { useAuth } from "@/lib/providers";
import { projectApi } from "@/lib/api";
import clsx from "clsx";
import {
  Home,
  CheckSquare,
  Inbox,
  Plus,
  ChevronDown,
  LogOut,
  X,
} from "lucide-react";

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function Sidebar({ isOpen, onClose }: SidebarProps) {
  const pathname = usePathname();
  const { user, workspace, workspaces, selectWorkspace, logout } = useAuth();
  const [wsOpen, setWsOpen] = useState(false);

  const { data: projects = [] } = useQuery({
    queryKey: ["projects", workspace?.id],
    queryFn: () => projectApi.list(workspace!.id),
    enabled: !!workspace?.id,
  });

  const nav = [
    { href: "/dashboard", label: "Home", icon: Home },
    { href: "/my-tasks", label: "My Tasks", icon: CheckSquare },
    { href: "/inbox", label: "Inbox", icon: Inbox },
  ];

  return (
    <>
      {isOpen && (
        <div className="fixed inset-0 z-40 bg-black/50 lg:hidden" onClick={onClose} />
      )}

      <aside
        className={clsx(
          "fixed top-0 left-0 z-50 flex h-full w-[220px] flex-col bg-sidebar border-r border-border lg:relative lg:z-auto lg:translate-x-0 transition-transform duration-200",
          isOpen ? "translate-x-0" : "-translate-x-full"
        )}
      >
        <button
          onClick={onClose}
          className="absolute top-2.5 right-2.5 rounded p-1 text-text-secondary hover:bg-sidebar-hover lg:hidden"
        >
          <X size={16} />
        </button>

        {/* Workspace name */}
        <div className="flex h-11 items-center px-4 border-b border-border">
          <span className="text-[15px] font-bold text-text-primary tracking-tight truncate">
            {workspace?.name || "Workspace"}
          </span>
        </div>

        {/* Navigation */}
        <nav className="mt-2 space-y-px px-2">
          {nav.map((item) => {
            const active = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={clsx(
                  "flex items-center gap-2.5 rounded-md px-2.5 py-[6px] text-[13px] transition-colors",
                  active
                    ? "bg-sidebar-active text-text-primary font-medium"
                    : "text-text-secondary hover:bg-sidebar-hover hover:text-text-primary"
                )}
              >
                <item.icon size={15} strokeWidth={active ? 2.2 : 1.8} />
                {item.label}
              </Link>
            );
          })}
        </nav>

        {/* Divider */}
        <div className="mx-4 my-3 h-px bg-border" />

        {/* Projects section */}
        <div className="flex items-center justify-between px-4 mb-1">
          <span className="text-[11px] font-semibold uppercase tracking-wider text-text-tertiary">
            Projects
          </span>
          <button className="rounded p-0.5 text-text-tertiary hover:bg-sidebar-hover hover:text-text-primary transition-colors">
            <Plus size={14} />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto sidebar-scroll px-2">
          {projects.filter((p) => !p.is_archived).map((project) => {
            const active = pathname === `/projects/${project.id}`;
            return (
              <Link
                key={project.id}
                href={`/projects/${project.id}`}
                className={clsx(
                  "flex items-center gap-2.5 rounded-md px-2.5 py-[6px] text-[13px] transition-colors",
                  active
                    ? "bg-sidebar-active text-text-primary font-medium"
                    : "text-text-secondary hover:bg-sidebar-hover hover:text-text-primary"
                )}
              >
                <span
                  className="h-2 w-2 shrink-0 rounded-full"
                  style={{ backgroundColor: project.color || "#4573d2" }}
                />
                <span className="truncate">{project.name}</span>
              </Link>
            );
          })}
        </div>

        {/* Bottom user section */}
        <div className="border-t border-border p-2">
          <div className="relative">
            <button
              onClick={() => setWsOpen(!wsOpen)}
              className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-[13px] text-text-secondary hover:bg-sidebar-hover transition-colors"
            >
              <div className="flex h-6 w-6 items-center justify-center rounded-full bg-accent text-[10px] font-semibold text-white">
                {user?.name?.charAt(0)?.toUpperCase() || "U"}
              </div>
              <div className="flex-1 text-left min-w-0">
                <p className="truncate text-[12px] text-text-primary">{user?.name}</p>
              </div>
              <ChevronDown size={12} className={clsx("text-text-tertiary transition-transform", wsOpen && "rotate-180")} />
            </button>

            {wsOpen && (
              <div className="absolute bottom-full left-0 mb-1 w-full rounded-md border border-border-strong bg-surface py-1 shadow-xl">
                {workspaces.map((ws) => (
                  <button
                    key={ws.id}
                    onClick={() => { selectWorkspace(ws); setWsOpen(false); }}
                    className={clsx(
                      "flex w-full items-center gap-2 px-3 py-1.5 text-left text-[12px] hover:bg-surface-hover transition-colors",
                      ws.id === workspace?.id ? "text-text-primary" : "text-text-secondary"
                    )}
                  >
                    {ws.name}
                  </button>
                ))}
                <div className="my-1 h-px bg-border" />
                <button
                  onClick={() => { logout(); setWsOpen(false); }}
                  className="flex w-full items-center gap-2 px-3 py-1.5 text-[12px] text-danger hover:bg-surface-hover transition-colors"
                >
                  <LogOut size={12} />
                  Log out
                </button>
              </div>
            )}
          </div>
        </div>
      </aside>
    </>
  );
}
