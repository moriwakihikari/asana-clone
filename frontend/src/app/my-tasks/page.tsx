"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuth } from "@/lib/providers";
import { taskApi } from "@/lib/api";
import { CheckCircle2, Circle, Calendar } from "lucide-react";
import { format, isPast, parseISO } from "date-fns";
import clsx from "clsx";
import PriorityBadge from "@/components/common/PriorityBadge";
import Avatar from "@/components/common/Avatar";
import type { Task } from "@/types";

export default function MyTasksPage() {
  const { workspace } = useAuth();
  const queryClient = useQueryClient();

  const { data: tasks = [] } = useQuery({
    queryKey: ["my-tasks", workspace?.id],
    queryFn: () => taskApi.getMyTasks(workspace!.id),
    enabled: !!workspace?.id,
  });

  const toggleStatus = useMutation({
    mutationFn: (task: Task) =>
      taskApi.changeStatus(task.id, task.status === "done" ? "todo" : "done"),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["my-tasks"] }),
  });

  const activeTasks = tasks.filter((t) => t.status !== "done" && t.status !== "cancelled");
  const completedTasks = tasks.filter((t) => t.status === "done");

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <h1 className="text-xl font-bold text-text-primary mb-6">My Tasks</h1>

      {/* Active tasks */}
      <div className="rounded-lg border border-border bg-surface divide-y divide-border">
        {activeTasks.length === 0 ? (
          <div className="py-10 text-center">
            <CheckCircle2 size={24} className="mx-auto text-text-tertiary mb-2" />
            <p className="text-[13px] text-text-tertiary">No active tasks</p>
          </div>
        ) : (
          activeTasks.map((task) => (
            <TaskItem key={task.id} task={task} onToggle={() => toggleStatus.mutate(task)} />
          ))
        )}
      </div>

      {/* Completed */}
      {completedTasks.length > 0 && (
        <div className="mt-6">
          <h2 className="text-[13px] font-semibold text-text-secondary mb-2">
            Completed ({completedTasks.length})
          </h2>
          <div className="rounded-lg border border-border bg-surface divide-y divide-border">
            {completedTasks.map((task) => (
              <TaskItem key={task.id} task={task} onToggle={() => toggleStatus.mutate(task)} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function TaskItem({ task, onToggle }: { task: Task; onToggle: () => void }) {
  const isDone = task.status === "done";
  const isOverdue = task.due_date && !isDone && isPast(parseISO(task.due_date));

  return (
    <div className="flex items-center gap-3 px-4 py-2.5 hover:bg-surface-hover transition-colors">
      <button onClick={onToggle} className="shrink-0">
        {isDone ? (
          <CheckCircle2 size={16} className="text-success" />
        ) : (
          <Circle size={16} className="text-text-tertiary hover:text-text-secondary" />
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
      </div>
    </div>
  );
}
