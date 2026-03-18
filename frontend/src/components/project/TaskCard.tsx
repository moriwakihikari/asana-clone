"use client";

import clsx from "clsx";
import { Calendar } from "lucide-react";
import { format, isPast, isToday, parseISO } from "date-fns";
import type { Task } from "@/types";
import Avatar from "@/components/common/Avatar";

interface TaskCardProps {
  task: Task;
  onClick?: (task: Task) => void;
  isDragging?: boolean;
}

const priorityBorderColors: Record<string, string> = {
  urgent: "#e8615a",
  high: "#f1914e",
  medium: "#f1bd6c",
  low: "#4573d2",
  none: "transparent",
};

const priorityDotColors: Record<string, string> = {
  urgent: "bg-[#e8615a]",
  high: "bg-[#f1914e]",
  medium: "bg-[#f1bd6c]",
  low: "bg-[#4573d2]",
  none: "bg-transparent",
};

export default function TaskCard({ task, onClick, isDragging }: TaskCardProps) {
  const isOverdue =
    task.due_date &&
    task.status !== "done" &&
    isPast(parseISO(task.due_date)) &&
    !isToday(parseISO(task.due_date));

  return (
    <div
      className={clsx(
        "group rounded-lg border border-border bg-surface cursor-pointer transition-all duration-100",
        isDragging
          ? "shadow-xl opacity-90 rotate-[1deg]"
          : "hover:bg-surface-hover"
      )}
      style={{ borderLeftWidth: "3px", borderLeftColor: priorityBorderColors[task.priority] }}
      onClick={() => onClick?.(task)}
    >
      <div className="px-3 py-2.5">
        <p
          className={clsx(
            "text-[14px] leading-snug",
            task.status === "done" ? "text-text-tertiary line-through" : "text-text-primary"
          )}
        >
          {task.title}
        </p>

        <div className="mt-2 flex items-center gap-2 flex-wrap">
          {task.priority !== "none" && (
            <span
              className={clsx(
                "inline-block h-1.5 w-1.5 rounded-full shrink-0",
                priorityDotColors[task.priority]
              )}
              title={task.priority}
            />
          )}

          {task.labels?.map((label) => (
            <span
              key={label.id}
              className="inline-flex items-center rounded-sm px-1.5 py-0 text-[10px] font-medium"
              style={{
                backgroundColor: label.color + "25",
                color: label.color,
              }}
            >
              {label.name}
            </span>
          ))}

          {task.due_date && (
            <span
              className={clsx(
                "inline-flex items-center gap-0.5 text-[11px] shrink-0",
                isOverdue ? "text-danger font-medium" : "text-text-tertiary"
              )}
            >
              <Calendar className="h-2.5 w-2.5" />
              {format(parseISO(task.due_date), "MMM d")}
            </span>
          )}

          <div className="flex-1" />

          {task.assignee && (
            <Avatar
              name={task.assignee.name}
              avatarUrl={task.assignee.avatar_url}
              size="sm"
            />
          )}
        </div>
      </div>
    </div>
  );
}
