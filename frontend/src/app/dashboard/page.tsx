"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import {
  CheckSquare,
  CheckCircle2,
  Clock,
  AlertTriangle,
  ArrowRight,
} from "lucide-react";
import { useAuth } from "@/lib/providers";
import { projectApi, taskApi } from "@/lib/api";
import StatsCard from "@/components/dashboard/StatsCard";
import TaskRow from "@/components/dashboard/TaskRow";
import ProjectCard from "@/components/dashboard/ProjectCard";
import type { Task, Project } from "@/types";

export default function DashboardPage() {
  const { user, workspace } = useAuth();
  const queryClient = useQueryClient();
  const today = format(new Date(), "EEEE, MMMM d, yyyy");

  const { data: myTasks = [] } = useQuery({
    queryKey: ["my-tasks", workspace?.id],
    queryFn: () => taskApi.getMyTasks(workspace!.id),
    enabled: !!workspace?.id,
  });

  const { data: projects = [] } = useQuery({
    queryKey: ["projects", workspace?.id],
    queryFn: () => projectApi.list(workspace!.id),
    enabled: !!workspace?.id,
  });

  const toggleStatus = useMutation({
    mutationFn: (task: Task) =>
      taskApi.changeStatus(task.id, task.status === "done" ? "todo" : "done"),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["my-tasks"] }),
  });

  const totalTasks = myTasks.length;
  const completedTasks = myTasks.filter((t) => t.status === "done").length;
  const inProgressTasks = myTasks.filter((t) => t.status === "in_progress").length;
  const overdueTasks = myTasks.filter(
    (t) => t.due_date && t.status !== "done" && new Date(t.due_date) < new Date()
  ).length;

  const activeTasks = myTasks.filter((t) => t.status !== "done" && t.status !== "cancelled");
  const greeting = new Date().getHours() < 12 ? "Good morning" : new Date().getHours() < 17 ? "Good afternoon" : "Good evening";

  return (
    <div className="mx-auto max-w-4xl px-6 py-8">
      {/* Header */}
      <div className="mb-6">
        <p className="text-[13px] text-text-tertiary">{today}</p>
        <h1 className="text-xl font-semibold text-text-primary mt-0.5">
          {greeting}, {user?.name?.split(" ")[0] || "there"}
        </h1>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-4 mb-8">
        <StatsCard icon={CheckSquare} label="Total Tasks" value={totalTasks} color="#6f6e6f" />
        <StatsCard icon={CheckCircle2} label="Completed" value={completedTasks} color="#5da283" />
        <StatsCard icon={Clock} label="In Progress" value={inProgressTasks} color="#4573d2" />
        <StatsCard icon={AlertTriangle} label="Overdue" value={overdueTasks} color="#e8615a" />
      </div>

      {/* My Tasks */}
      <section className="mb-8">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-[15px] font-semibold text-text-primary">My Tasks</h2>
          <button className="flex items-center gap-1 text-[13px] text-accent hover:text-accent-hover font-medium transition-colors">
            View all <ArrowRight size={13} />
          </button>
        </div>
        <div className="rounded-lg border border-border bg-surface divide-y divide-border">
          {activeTasks.length === 0 ? (
            <div className="py-10 text-center">
              <CheckCircle2 size={24} className="mx-auto text-text-tertiary mb-2" />
              <p className="text-[13px] text-text-tertiary">All caught up!</p>
            </div>
          ) : (
            activeTasks.slice(0, 6).map((task) => (
              <TaskRow
                key={task.id}
                task={task}
                projectName={projects.find((p) => p.id === task.project_id)?.name}
                onToggleStatus={(t) => toggleStatus.mutate(t)}
              />
            ))
          )}
        </div>
      </section>

      {/* Recent Projects */}
      <section>
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-[15px] font-semibold text-text-primary">Projects</h2>
          <button className="flex items-center gap-1 text-[13px] text-accent hover:text-accent-hover font-medium transition-colors">
            All projects <ArrowRight size={13} />
          </button>
        </div>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {projects.slice(0, 6).map((project) => {
            const projectTasks = myTasks.filter((t) => t.project_id === project.id);
            const completed = projectTasks.filter((t) => t.status === "done").length;
            return (
              <ProjectCard
                key={project.id}
                project={project}
                taskCount={projectTasks.length}
                completedCount={completed}
              />
            );
          })}
        </div>
      </section>
    </div>
  );
}
