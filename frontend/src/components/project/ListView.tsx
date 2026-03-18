"use client";

import { useState } from "react";
import { ChevronDown, ChevronRight, ArrowUpDown, CheckCircle2, Circle } from "lucide-react";
import clsx from "clsx";
import { format, isPast, isToday, parseISO } from "date-fns";
import type { Section, Task } from "@/types";
import Avatar from "@/components/common/Avatar";
import PriorityBadge from "@/components/common/PriorityBadge";
import StatusBadge from "@/components/common/StatusBadge";

interface ListViewProps {
  sections: Section[];
  tasksBySection: Record<string, Task[]>;
  onClickTask: (task: Task) => void;
  onToggleStatus: (task: Task) => void;
}

type SortKey = "title" | "priority" | "due_date" | "status";
type SortDir = "asc" | "desc";

const priorityOrder = { urgent: 0, high: 1, medium: 2, low: 3, none: 4 };
const statusOrder = { todo: 0, in_progress: 1, in_review: 2, done: 3, cancelled: 4 };

export default function ListView({ sections, tasksBySection, onClickTask, onToggleStatus }: ListViewProps) {
  const [collapsedSections, setCollapsedSections] = useState<Set<string>>(new Set());
  const [sortKey, setSortKey] = useState<SortKey>("title");
  const [sortDir, setSortDir] = useState<SortDir>("asc");

  const toggleSection = (id: string) => {
    setCollapsedSections((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(key);
      setSortDir("asc");
    }
  };

  const sortTasks = (tasks: Task[]): Task[] => {
    return [...tasks].sort((a, b) => {
      let cmp = 0;
      switch (sortKey) {
        case "title":
          cmp = a.title.localeCompare(b.title);
          break;
        case "priority":
          cmp = priorityOrder[a.priority] - priorityOrder[b.priority];
          break;
        case "due_date": {
          const da = a.due_date || "9999";
          const db = b.due_date || "9999";
          cmp = da.localeCompare(db);
          break;
        }
        case "status":
          cmp = statusOrder[a.status] - statusOrder[b.status];
          break;
      }
      return sortDir === "asc" ? cmp : -cmp;
    });
  };

  const SortHeader = ({ label, field }: { label: string; field: SortKey }) => (
    <button
      className="inline-flex items-center gap-1 text-[12px] font-medium text-text-tertiary hover:text-text-secondary transition-colors"
      onClick={() => handleSort(field)}
    >
      {label}
      <ArrowUpDown
        className={clsx(
          "h-3 w-3",
          sortKey === field ? "text-accent" : "text-text-tertiary"
        )}
      />
    </button>
  );

  return (
    <div className="rounded-lg border border-border bg-surface overflow-hidden">
      {/* Table Header */}
      <div className="grid grid-cols-[auto_1fr_120px_100px_100px_100px] gap-3 items-center px-4 py-2.5 bg-surface-hover border-b border-border text-[12px] font-medium text-text-tertiary">
        <div className="w-5" />
        <SortHeader label="Task Name" field="title" />
        <span>Assignee</span>
        <SortHeader label="Due Date" field="due_date" />
        <SortHeader label="Priority" field="priority" />
        <SortHeader label="Status" field="status" />
      </div>

      {/* Sections */}
      {sections.map((section) => {
        const tasks = sortTasks(tasksBySection[section.id] || []);
        const isCollapsed = collapsedSections.has(section.id);

        return (
          <div key={section.id}>
            {/* Section Header */}
            <button
              className="flex w-full items-center gap-2 px-4 py-2 bg-surface-hover border-b border-border hover:bg-surface-active transition-colors"
              onClick={() => toggleSection(section.id)}
            >
              {isCollapsed ? (
                <ChevronRight className="h-3.5 w-3.5 text-text-tertiary" />
              ) : (
                <ChevronDown className="h-3.5 w-3.5 text-text-tertiary" />
              )}
              <span className="text-[13px] font-semibold text-text-primary">{section.name}</span>
              <span className="text-[12px] text-text-tertiary">({tasks.length})</span>
            </button>

            {/* Tasks */}
            {!isCollapsed &&
              tasks.map((task, i) => {
                const isDone = task.status === "done";
                const isOverdue =
                  task.due_date &&
                  !isDone &&
                  isPast(parseISO(task.due_date)) &&
                  !isToday(parseISO(task.due_date));

                return (
                  <div
                    key={task.id}
                    className={clsx(
                      "grid grid-cols-[auto_1fr_120px_100px_100px_100px] gap-3 items-center px-4 py-2.5 border-b border-border cursor-pointer transition-colors",
                      "hover:bg-surface-hover",
                      i % 2 === 1 && "bg-bg-raised/30"
                    )}
                    onClick={() => onClickTask(task)}
                  >
                    {/* Checkbox */}
                    <button
                      className="w-5"
                      onClick={(e) => {
                        e.stopPropagation();
                        onToggleStatus(task);
                      }}
                    >
                      {isDone ? (
                        <CheckCircle2 className="h-4 w-4 text-success" />
                      ) : (
                        <Circle className="h-4 w-4 text-text-tertiary hover:text-text-secondary" />
                      )}
                    </button>

                    {/* Title */}
                    <span
                      className={clsx(
                        "text-[13px] truncate",
                        isDone ? "text-text-tertiary line-through" : "text-text-primary"
                      )}
                    >
                      {task.title}
                    </span>

                    {/* Assignee */}
                    <div>
                      {task.assignee ? (
                        <div className="flex items-center gap-1.5">
                          <Avatar name={task.assignee.name} size="sm" />
                          <span className="text-[12px] text-text-secondary truncate">
                            {task.assignee.name.split(" ")[0]}
                          </span>
                        </div>
                      ) : (
                        <span className="text-[12px] text-text-tertiary">--</span>
                      )}
                    </div>

                    {/* Due Date */}
                    <span
                      className={clsx(
                        "text-[12px]",
                        isOverdue ? "text-danger font-medium" : "text-text-secondary"
                      )}
                    >
                      {task.due_date
                        ? format(parseISO(task.due_date), "MMM d")
                        : "--"}
                    </span>

                    {/* Priority */}
                    <PriorityBadge priority={task.priority} />

                    {/* Status */}
                    <StatusBadge status={task.status} />
                  </div>
                );
              })}
          </div>
        );
      })}
    </div>
  );
}
