"use client";

import clsx from "clsx";
import type { TaskPriority } from "@/types";

const config: Record<string, { label: string; dotColor: string }> = {
  urgent: { label: "Urgent", dotColor: "bg-[#e8615a]" },
  high: { label: "High", dotColor: "bg-[#f1914e]" },
  medium: { label: "Medium", dotColor: "bg-[#f1bd6c]" },
  low: { label: "Low", dotColor: "bg-[#4573d2]" },
};

export default function PriorityBadge({ priority }: { priority: TaskPriority }) {
  const c = config[priority];
  if (!c) return null;

  return (
    <span className="inline-flex items-center gap-1.5 text-[11px] text-text-secondary">
      <span className={clsx("h-1.5 w-1.5 rounded-full", c.dotColor)} />
      {c.label}
    </span>
  );
}
