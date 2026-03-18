"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { X } from "lucide-react";
import clsx from "clsx";
import type { Section, TaskPriority } from "@/types";
import { taskApi } from "@/lib/api";

interface CreateTaskModalProps {
  projectId: string;
  sections: Section[];
  defaultSectionId?: string;
  onClose: () => void;
  onCreated?: () => void;
}

const priorities: { value: TaskPriority; label: string; color: string }[] = [
  { value: "none", label: "None", color: "bg-text-tertiary" },
  { value: "low", label: "Low", color: "bg-[#4573d2]" },
  { value: "medium", label: "Medium", color: "bg-[#f1bd6c]" },
  { value: "high", label: "High", color: "bg-[#f1914e]" },
  { value: "urgent", label: "Urgent", color: "bg-[#e8615a]" },
];

export default function CreateTaskModal({
  projectId,
  sections,
  defaultSectionId,
  onClose,
  onCreated,
}: CreateTaskModalProps) {
  const queryClient = useQueryClient();
  const [title, setTitle] = useState("");
  const [sectionId, setSectionId] = useState(defaultSectionId || sections[0]?.id || "");
  const [priority, setPriority] = useState<TaskPriority>("none");
  const [dueDate, setDueDate] = useState("");

  const createMutation = useMutation({
    mutationFn: () =>
      taskApi.create(projectId, {
        title,
        section_id: sectionId || undefined,
        priority: priority !== "none" ? priority : undefined,
        due_date: dueDate || undefined,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks", projectId] });
      queryClient.invalidateQueries({ queryKey: ["my-tasks"] });
      onCreated?.();
      onClose();
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    createMutation.mutate();
  };

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-40 bg-black/50"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div
          className="w-full max-w-md rounded-xl bg-surface shadow-2xl border border-border animate-scale-in"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="flex items-center justify-between border-b border-border px-5 py-4">
            <h2 className="text-[15px] font-semibold text-text-primary">Create Task</h2>
            <button
              onClick={onClose}
              className="rounded-lg p-1 text-text-tertiary hover:bg-surface-hover hover:text-text-secondary transition-colors"
            >
              <X className="h-4 w-4" />
            </button>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="p-5 space-y-4">
            {/* Title */}
            <div>
              <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">
                Title <span className="text-danger">*</span>
              </label>
              <input
                autoFocus
                className="w-full rounded-lg border border-border bg-bg px-3 py-2.5 text-[13px] text-text-primary outline-none focus:border-accent placeholder:text-text-tertiary transition-colors"
                placeholder="What needs to be done?"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>

            {/* Section */}
            {sections.length > 0 && (
              <div>
                <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">Section</label>
                <select
                  className="w-full rounded-lg border border-border bg-bg px-3 py-2.5 text-[13px] text-text-primary outline-none focus:border-accent transition-colors"
                  value={sectionId}
                  onChange={(e) => setSectionId(e.target.value)}
                >
                  {sections.map((s) => (
                    <option key={s.id} value={s.id}>
                      {s.name}
                    </option>
                  ))}
                </select>
              </div>
            )}

            {/* Priority */}
            <div>
              <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">Priority</label>
              <div className="flex gap-2">
                {priorities.map((p) => (
                  <button
                    key={p.value}
                    type="button"
                    onClick={() => setPriority(p.value)}
                    className={clsx(
                      "flex-1 rounded-lg border px-2 py-2 text-[12px] font-medium transition-all",
                      priority === p.value
                        ? "border-accent bg-accent/10 text-accent"
                        : "border-border bg-bg text-text-secondary hover:bg-surface-hover"
                    )}
                  >
                    <div className="flex items-center justify-center gap-1.5">
                      <span className={clsx("h-2 w-2 rounded-full", p.color)} />
                      {p.label}
                    </div>
                  </button>
                ))}
              </div>
            </div>

            {/* Due Date */}
            <div>
              <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">Due Date</label>
              <input
                type="date"
                className="w-full rounded-lg border border-border bg-bg px-3 py-2.5 text-[13px] text-text-primary outline-none focus:border-accent transition-colors"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
              />
            </div>

            {/* Actions */}
            <div className="flex justify-end gap-2 pt-2">
              <button
                type="button"
                onClick={onClose}
                className="rounded-lg border border-border px-4 py-2 text-[13px] font-medium text-text-secondary hover:bg-surface-hover transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={!title.trim() || createMutation.isPending}
                className="rounded-lg bg-accent px-4 py-2 text-[13px] font-medium text-white hover:bg-accent-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {createMutation.isPending ? "Creating..." : "Create Task"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  );
}
