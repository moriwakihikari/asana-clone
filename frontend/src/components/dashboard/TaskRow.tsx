"use client";

import clsx from "clsx";
import { Calendar, Circle, CheckCircle2 } from "lucide-react";
import { format, isPast, parseISO } from "date-fns";
import type { Task } from "@/types";
import PriorityBadge from "@/components/common/PriorityBadge";
import Avatar from "@/components/common/Avatar";

interface TaskRowProps {
  task: Task;
  projectName?: string;
  onToggleStatus?: (task: Task) => void;
  onClick?: (task: Task) => void;
}

export default function TaskRow({ task, projectName, onToggleStatus, onClick }: TaskRowProps) {
  const isDone = task.status === "done";
  const isOverdue = task.due_date && !isDone && isPast(parseISO(task.due_date));

  return (
    <div
      className="group flex items-center gap-3 px-4 py-2.5 hover:bg-surface-hover cursor-pointer transition-colors"
      onClick={() => onClick?.(task)}
    >
      <button
        className="shrink-0"
        onClick={(e) => { e.stopPropagation(); onToggleStatus?.(task); }}
      >
        {isDone ? (
          <CheckCircle2 size={16} className="text-success" />
        ) : (
          <Circle size={16} className="text-text-tertiary group-hover:text-text-secondary" />
        )}
      </button>

      <span className={clsx("flex-1 text-[13px] truncate", isDone ? "text-text-tertiary line-through" : "text-text-primary")}>
        {task.title}
      </span>

      <div className="flex items-center gap-2 shrink-0">
        {task.priority !== "none" && <PriorityBadge priority={task.priority} />}

        {task.due_date && (
          <span className={clsx("flex items-center gap-1 text-[12px]", isOverdue ? "text-danger" : "text-text-tertiary")}>
            <Calendar size={11} />
            {format(parseISO(task.due_date), "MMM d")}
          </span>
        )}

        {task.assignee && <Avatar name={task.assignee.name} size="sm" />}

        {projectName && (
          <span className="text-[11px] text-text-tertiary bg-surface-hover rounded px-1.5 py-0.5 max-w-[100px] truncate">
            {projectName}
          </span>
        )}
      </div>
    </div>
  );
}
