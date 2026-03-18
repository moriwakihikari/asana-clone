"use client";

import { useState } from "react";
import { useAuth } from "@/lib/providers";
import Sidebar from "./Sidebar";
import { Menu, Search, Bell } from "lucide-react";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { user } = useAuth();

  return (
    <div className="flex h-screen overflow-hidden bg-bg">
      <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Top header bar */}
        <header className="flex h-11 shrink-0 items-center gap-3 border-b border-border px-4 bg-bg">
          <button
            onClick={() => setSidebarOpen(true)}
            className="rounded p-1 text-text-secondary hover:bg-surface-hover lg:hidden"
          >
            <Menu size={18} />
          </button>

          <div className="relative flex-1 max-w-sm">
            <Search size={14} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-text-tertiary" />
            <input
              type="text"
              placeholder="Search"
              className="h-7 w-full rounded-md border border-border bg-surface pl-8 pr-10 text-[13px] text-text-primary placeholder:text-text-tertiary focus:border-accent focus:outline-none"
            />
            <kbd className="absolute right-2 top-1/2 -translate-y-1/2 text-[10px] text-text-tertiary bg-surface-hover rounded px-1 py-0.5 border border-border">
              ⌘K
            </kbd>
          </div>

          <div className="flex-1" />

          <button className="relative rounded p-1.5 text-text-secondary hover:bg-surface-hover transition-colors">
            <Bell size={16} />
            <span className="absolute top-1 right-1 h-1.5 w-1.5 rounded-full bg-danger" />
          </button>
          <button className="flex h-7 w-7 items-center justify-center rounded-full bg-accent text-[11px] font-medium text-white">
            {user?.name?.charAt(0)?.toUpperCase() || "U"}
          </button>
        </header>

        <main className="flex-1 overflow-y-auto bg-bg">{children}</main>
      </div>
    </div>
  );
}
