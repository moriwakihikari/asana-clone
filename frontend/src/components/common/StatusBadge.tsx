"use client";

import clsx from "clsx";
import type { TaskStatus } from "@/types";

const config: Record<string, { label: string; dotColor: string }> = {
  todo: { label: "To do", dotColor: "bg-text-tertiary" },
  in_progress: { label: "In progress", dotColor: "bg-[#4573d2]" },
  in_review: { label: "In review", dotColor: "bg-[#f1bd6c]" },
  done: { label: "Done", dotColor: "bg-[#5da283]" },
  cancelled: { label: "Cancelled", dotColor: "bg-[#e8615a]" },
};

export default function StatusBadge({ status }: { status: TaskStatus }) {
  const c = config[status] || config.todo;

  return (
    <span className="inline-flex items-center gap-1.5 text-[11px] text-text-secondary">
      <span className={clsx("h-1.5 w-1.5 rounded-full", c.dotColor)} />
      {c.label}
    </span>
  );
}
