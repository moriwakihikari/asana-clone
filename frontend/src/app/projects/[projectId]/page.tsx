"use client";

import { useState, useMemo, use } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  LayoutGrid,
  List,
  Plus,
  Settings,
  Loader2,
  FolderKanban,
} from "lucide-react";
import clsx from "clsx";
import { useAuth } from "@/lib/providers";
import { projectApi, sectionApi, taskApi } from "@/lib/api";
import BoardView from "@/components/project/BoardView";
import ListView from "@/components/project/ListView";
import TaskDetailPanel from "@/components/task/TaskDetailPanel";
import CreateTaskModal from "@/components/task/CreateTaskModal";
import type { Project, Section, Task } from "@/types";

// Mock data for development
const mockProject: Project = {
  id: "p1",
  workspace_id: "w1",
  name: "Website Redesign",
  description: "Complete website overhaul with new branding and improved UX",
  color: "#6366f1",
  view_type: "board",
  is_archived: false,
  created_at: "2026-02-01T00:00:00Z",
  updated_at: "2026-03-17T00:00:00Z",
};

const mockSections: Section[] = [
  { id: "s1", project_id: "p1", name: "Backlog", position: 0, created_at: "2026-02-01T00:00:00Z" },
  { id: "s2", project_id: "p1", name: "To Do", position: 1, created_at: "2026-02-01T00:00:00Z" },
  { id: "s3", project_id: "p1", name: "In Progress", position: 2, created_at: "2026-02-01T00:00:00Z" },
  { id: "s4", project_id: "p1", name: "Review", position: 3, created_at: "2026-02-01T00:00:00Z" },
  { id: "s5", project_id: "p1", name: "Done", position: 4, created_at: "2026-02-01T00:00:00Z" },
];

const mockTasks: Task[] = [
  {
    id: "t1", project_id: "p1", section_id: "s2", title: "Design homepage mockup", status: "todo",
    priority: "high", due_date: "2026-03-20T00:00:00Z", position: 0,
    assignee: { id: "u1", name: "Alice Chen", email: "alice@example.com" },
    created_at: "2026-03-10T00:00:00Z", updated_at: "2026-03-17T00:00:00Z",
  },
  {
    id: "t2", project_id: "p1", section_id: "s3", title: "Implement navigation component", status: "in_progress",
    priority: "urgent", due_date: "2026-03-19T00:00:00Z", position: 0,
    assignee: { id: "u2", name: "Bob Smith", email: "bob@example.com" },
    created_at: "2026-03-12T00:00:00Z", updated_at: "2026-03-16T00:00:00Z",
  },
  {
    id: "t3", project_id: "p1", section_id: "s2", title: "Set up color palette and design tokens", status: "todo",
    priority: "medium", position: 1,
    created_at: "2026-03-14T00:00:00Z", updated_at: "2026-03-17T00:00:00Z",
  },
  {
    id: "t4", project_id: "p1", section_id: "s4", title: "Review PR #42 - Auth module", status: "in_review",
    priority: "high", due_date: "2026-03-18T00:00:00Z", position: 0,
    assignee: { id: "u1", name: "Alice Chen", email: "alice@example.com" },
    created_at: "2026-03-08T00:00:00Z", updated_at: "2026-03-15T00:00:00Z",
  },
  {
    id: "t5", project_id: "p1", section_id: "s5", title: "Project setup and boilerplate", status: "done",
    priority: "none", position: 0,
    created_at: "2026-02-01T00:00:00Z", updated_at: "2026-02-05T00:00:00Z",
  },
  {
    id: "t6", project_id: "p1", section_id: "s1", title: "Implement dark mode toggle", status: "todo",
    priority: "low", position: 0,
    created_at: "2026-03-15T00:00:00Z", updated_at: "2026-03-15T00:00:00Z",
  },
  {
    id: "t7", project_id: "p1", section_id: "s3", title: "Build dashboard widgets", status: "in_progress",
    priority: "medium", due_date: "2026-03-22T00:00:00Z", position: 1,
    assignee: { id: "u3", name: "Carol Davis", email: "carol@example.com" },
    created_at: "2026-03-13T00:00:00Z", updated_at: "2026-03-17T00:00:00Z",
  },
  {
    id: "t8", project_id: "p1", section_id: "s1", title: "Accessibility audit", status: "todo",
    priority: "medium", position: 1,
    created_at: "2026-03-16T00:00:00Z", updated_at: "2026-03-16T00:00:00Z",
  },
];

export default function ProjectPage({ params }: { params: Promise<{ projectId: string }> }) {
  const { projectId } = use(params);
  const { workspace } = useAuth();
  const queryClient = useQueryClient();

  const [viewMode, setViewMode] = useState<"board" | "list">("board");
  const [selectedTaskId, setSelectedTaskId] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createSectionId, setCreateSectionId] = useState<string | undefined>();

  // Fetch project
  const { data: project = mockProject } = useQuery({
    queryKey: ["project", workspace?.id, projectId],
    queryFn: () => projectApi.get(workspace!.id, projectId),
    enabled: !!workspace?.id,
    placeholderData: mockProject,
  });

  // Fetch sections
  const { data: sections = mockSections } = useQuery({
    queryKey: ["sections", projectId],
    queryFn: () => sectionApi.list(projectId),
    enabled: !!projectId,
    placeholderData: mockSections,
  });

  // Fetch tasks
  const { data: tasks = mockTasks, isLoading } = useQuery({
    queryKey: ["tasks", projectId],
    queryFn: () => taskApi.list(projectId),
    enabled: !!projectId,
    placeholderData: mockTasks,
  });

  // Group tasks by section
  const tasksBySection = useMemo(() => {
    const map: Record<string, Task[]> = {};
    for (const section of sections) {
      map[section.id] = [];
    }
    for (const task of tasks) {
      const sid = task.section_id || sections[0]?.id;
      if (sid && map[sid]) {
        map[sid].push(task);
      }
    }
    for (const key in map) {
      map[key].sort((a, b) => a.position - b.position);
    }
    return map;
  }, [tasks, sections]);

  // Mutations
  const moveTaskMutation = useMutation({
    mutationFn: ({ taskId, sectionId, position }: { taskId: string; sectionId: string; position: number }) =>
      taskApi.move(taskId, { section_id: sectionId, position }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks", projectId] });
    },
  });

  const changeStatusMutation = useMutation({
    mutationFn: (task: Task) =>
      taskApi.changeStatus(task.id, task.status === "done" ? "todo" : "done"),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks", projectId] });
    },
  });

  const createSectionMutation = useMutation({
    mutationFn: () => sectionApi.create(projectId, { name: "New Section" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sections", projectId] });
    },
  });

  const renameSectionMutation = useMutation({
    mutationFn: ({ sectionId, name }: { sectionId: string; name: string }) =>
      sectionApi.rename(projectId, sectionId, name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sections", projectId] });
    },
  });

  const deleteSectionMutation = useMutation({
    mutationFn: (sectionId: string) => sectionApi.delete(projectId, sectionId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sections", projectId] });
      queryClient.invalidateQueries({ queryKey: ["tasks", projectId] });
    },
  });

  const handleAddTask = (sectionId: string) => {
    setCreateSectionId(sectionId);
    setShowCreateModal(true);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-6 w-6 animate-spin text-accent" />
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Project Header */}
      <div className="border-b border-border bg-bg px-6 py-4">
        <div className="flex items-center gap-3 mb-1">
          <div
            className="flex h-8 w-8 items-center justify-center rounded-lg"
            style={{ backgroundColor: (project.color || "#4573d2") + "20" }}
          >
            <FolderKanban className="h-4 w-4" style={{ color: project.color || "#4573d2" }} />
          </div>
          <div>
            <h1 className="text-xl font-bold text-text-primary">{project.name}</h1>
            {project.description && (
              <p className="text-[13px] text-text-secondary mt-0.5">{project.description}</p>
            )}
          </div>
        </div>

        {/* Toolbar */}
        <div className="flex items-center justify-between mt-4">
          <div className="flex items-center gap-1 rounded-lg border border-border p-0.5 bg-surface">
            <button
              onClick={() => setViewMode("board")}
              className={clsx(
                "inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-[13px] font-medium transition-all",
                viewMode === "board"
                  ? "bg-surface-hover text-text-primary"
                  : "text-text-secondary hover:text-text-primary"
              )}
            >
              <LayoutGrid className="h-3.5 w-3.5" />
              Board
            </button>
            <button
              onClick={() => setViewMode("list")}
              className={clsx(
                "inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-[13px] font-medium transition-all",
                viewMode === "list"
                  ? "bg-surface-hover text-text-primary"
                  : "text-text-secondary hover:text-text-primary"
              )}
            >
              <List className="h-3.5 w-3.5" />
              List
            </button>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => {
                setCreateSectionId(sections[0]?.id);
                setShowCreateModal(true);
              }}
              className="inline-flex items-center gap-1.5 rounded-lg bg-accent px-3 py-1.5 text-[13px] font-medium text-white hover:bg-accent-hover transition-colors"
            >
              <Plus className="h-3.5 w-3.5" />
              Add Task
            </button>
            <button className="rounded-lg border border-border p-1.5 text-text-tertiary hover:bg-surface-hover hover:text-text-secondary transition-colors">
              <Settings className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto p-6 bg-bg">
        {viewMode === "board" ? (
          <BoardView
            sections={sections}
            tasksBySection={tasksBySection}
            onMoveTask={(taskId, sectionId, position) =>
              moveTaskMutation.mutate({ taskId, sectionId, position })
            }
            onClickTask={(task) => setSelectedTaskId(task.id)}
            onAddTask={handleAddTask}
            onAddSection={() => createSectionMutation.mutate()}
            onRenameSection={(sectionId, name) =>
              renameSectionMutation.mutate({ sectionId, name })
            }
            onDeleteSection={(sectionId) => deleteSectionMutation.mutate(sectionId)}
          />
        ) : (
          <ListView
            sections={sections}
            tasksBySection={tasksBySection}
            onClickTask={(task) => setSelectedTaskId(task.id)}
            onToggleStatus={(task) => changeStatusMutation.mutate(task)}
          />
        )}
      </div>

      {/* Task Detail Panel */}
      {selectedTaskId && (
        <TaskDetailPanel
          taskId={selectedTaskId}
          onClose={() => setSelectedTaskId(null)}
          onDeleted={() => setSelectedTaskId(null)}
        />
      )}

      {/* Create Task Modal */}
      {showCreateModal && (
        <CreateTaskModal
          projectId={projectId}
          sections={sections}
          defaultSectionId={createSectionId}
          onClose={() => setShowCreateModal(false)}
          onCreated={() => {
            queryClient.invalidateQueries({ queryKey: ["tasks", projectId] });
          }}
        />
      )}
    </div>
  );
}
