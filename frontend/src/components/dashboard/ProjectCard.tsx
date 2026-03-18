"use client";

import { FolderKanban, Clock } from "lucide-react";
import { formatDistanceToNow, parseISO } from "date-fns";
import { useRouter } from "next/navigation";
import type { Project } from "@/types";

interface ProjectCardProps {
  project: Project;
  taskCount?: number;
  completedCount?: number;
}

export default function ProjectCard({ project, taskCount = 0, completedCount = 0 }: ProjectCardProps) {
  const router = useRouter();
  const progress = taskCount > 0 ? Math.round((completedCount / taskCount) * 100) : 0;

  return (
    <div
      className="group rounded-lg border border-border bg-surface p-4 hover:bg-surface-hover cursor-pointer transition-colors"
      onClick={() => router.push(`/projects/${project.id}`)}
    >
      <div className="flex items-start gap-3">
        <div
          className="flex h-8 w-8 shrink-0 items-center justify-center rounded"
          style={{ backgroundColor: project.color + "20" }}
        >
          <FolderKanban size={14} style={{ color: project.color }} />
        </div>
        <div className="min-w-0 flex-1">
          <h3 className="text-[13px] font-semibold text-text-primary truncate">
            {project.name}
          </h3>
          {project.description && (
            <p className="mt-0.5 text-[12px] text-text-tertiary truncate">{project.description}</p>
          )}
        </div>
      </div>

      {taskCount > 0 && (
        <div className="mt-3">
          <div className="flex items-center justify-between text-[11px] text-text-tertiary mb-1">
            <span>{completedCount} of {taskCount} tasks</span>
            <span>{progress}%</span>
          </div>
          <div className="h-1 w-full rounded-full bg-bg">
            <div
              className="h-full rounded-full transition-all"
              style={{ width: `${progress}%`, backgroundColor: project.color }}
            />
          </div>
        </div>
      )}

      <div className="mt-3 flex items-center gap-1 text-[11px] text-text-tertiary">
        <Clock size={10} />
        Updated {formatDistanceToNow(parseISO(project.updated_at), { addSuffix: true })}
      </div>
    </div>
  );
}
