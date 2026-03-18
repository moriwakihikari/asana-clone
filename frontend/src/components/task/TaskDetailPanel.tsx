"use client";

import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { X, Trash2, Send, Calendar, User as UserIcon } from "lucide-react";
import clsx from "clsx";
import { format, parseISO } from "date-fns";
import type { Task, TaskStatus, TaskPriority, Comment } from "@/types";
import { taskApi, commentApi } from "@/lib/api";
import { useAuth } from "@/lib/providers";
import Avatar from "@/components/common/Avatar";
import StatusBadge from "@/components/common/StatusBadge";
import PriorityBadge from "@/components/common/PriorityBadge";

interface TaskDetailPanelProps {
  taskId: string | null;
  onClose: () => void;
  onDeleted?: () => void;
}

const statuses: TaskStatus[] = ["todo", "in_progress", "in_review", "done", "cancelled"];
const priorities: TaskPriority[] = ["none", "low", "medium", "high", "urgent"];

export default function TaskDetailPanel({ taskId, onClose, onDeleted }: TaskDetailPanelProps) {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState<TaskStatus>("todo");
  const [priority, setPriority] = useState<TaskPriority>("none");
  const [dueDate, setDueDate] = useState("");
  const [commentText, setCommentText] = useState("");
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const { data: task } = useQuery({
    queryKey: ["task", taskId],
    queryFn: () => taskApi.get(taskId!),
    enabled: !!taskId,
  });

  const { data: comments = [] } = useQuery({
    queryKey: ["comments", taskId],
    queryFn: () => commentApi.list(taskId!),
    enabled: !!taskId,
  });

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      setDescription(task.description || "");
      setStatus(task.status);
      setPriority(task.priority);
      setDueDate(task.due_date ? task.due_date.slice(0, 10) : "");
    }
  }, [task]);

  const updateMutation = useMutation({
    mutationFn: (data: Partial<Task>) => taskApi.update(taskId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["task", taskId] });
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      queryClient.invalidateQueries({ queryKey: ["my-tasks"] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: () => taskApi.delete(taskId!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      queryClient.invalidateQueries({ queryKey: ["my-tasks"] });
      onDeleted?.();
      onClose();
    },
  });

  const addCommentMutation = useMutation({
    mutationFn: (content: string) => commentApi.create(taskId!, content),
    onSuccess: () => {
      setCommentText("");
      queryClient.invalidateQueries({ queryKey: ["comments", taskId] });
    },
  });

  const saveField = (field: string, value: string) => {
    const data: Partial<Task> = {};
    if (field === "title" && value !== task?.title) data.title = value;
    if (field === "description" && value !== (task?.description || "")) data.description = value;
    if (field === "due_date") data.due_date = value || undefined;
    if (Object.keys(data).length > 0) {
      updateMutation.mutate(data);
    }
  };

  if (!taskId) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-40 bg-black/40 transition-opacity"
        onClick={onClose}
      />

      {/* Panel */}
      <div className="fixed inset-y-0 right-0 z-50 w-full max-w-lg bg-surface shadow-2xl border-l border-border flex flex-col animate-slide-in-right">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <div className="flex items-center gap-3">
            <StatusBadge status={status} />
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setShowDeleteConfirm(true)}
              className="rounded-lg p-1.5 text-text-tertiary hover:bg-surface-hover hover:text-danger transition-colors"
            >
              <Trash2 className="h-4 w-4" />
            </button>
            <button
              onClick={onClose}
              className="rounded-lg p-1.5 text-text-tertiary hover:bg-surface-hover hover:text-text-secondary transition-colors"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
        </div>

        {/* Body */}
        <div className="flex-1 overflow-y-auto px-6 py-5 space-y-6">
          {/* Title */}
          <input
            className="w-full text-xl font-semibold text-text-primary outline-none placeholder:text-text-tertiary bg-transparent"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            onBlur={() => saveField("title", title)}
            placeholder="Task title"
          />

          {/* Fields */}
          <div className="grid grid-cols-2 gap-4 text-[13px]">
            <div>
              <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">Status</label>
              <select
                className="w-full rounded-lg border border-border bg-bg px-3 py-2 text-[13px] text-text-primary outline-none focus:border-accent"
                value={status}
                onChange={(e) => {
                  const s = e.target.value as TaskStatus;
                  setStatus(s);
                  taskApi.changeStatus(taskId!, s).then(() => {
                    queryClient.invalidateQueries({ queryKey: ["task", taskId] });
                    queryClient.invalidateQueries({ queryKey: ["tasks"] });
                    queryClient.invalidateQueries({ queryKey: ["my-tasks"] });
                  });
                }}
              >
                {statuses.map((s) => (
                  <option key={s} value={s}>
                    {s.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase())}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">Priority</label>
              <select
                className="w-full rounded-lg border border-border bg-bg px-3 py-2 text-[13px] text-text-primary outline-none focus:border-accent"
                value={priority}
                onChange={(e) => {
                  const p = e.target.value as TaskPriority;
                  setPriority(p);
                  updateMutation.mutate({ priority: p });
                }}
              >
                {priorities.map((p) => (
                  <option key={p} value={p}>
                    {p.charAt(0).toUpperCase() + p.slice(1)}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">
                <span className="flex items-center gap-1">
                  <Calendar className="h-3 w-3" /> Due Date
                </span>
              </label>
              <input
                type="date"
                className="w-full rounded-lg border border-border bg-bg px-3 py-2 text-[13px] text-text-primary outline-none focus:border-accent"
                value={dueDate}
                onChange={(e) => {
                  setDueDate(e.target.value);
                  saveField("due_date", e.target.value);
                }}
              />
            </div>

            <div>
              <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">
                <span className="flex items-center gap-1">
                  <UserIcon className="h-3 w-3" /> Assignee
                </span>
              </label>
              <div className="flex items-center gap-2 rounded-lg border border-border bg-bg px-3 py-2">
                {task?.assignee ? (
                  <>
                    <Avatar name={task.assignee.name} size="sm" />
                    <span className="text-[13px] text-text-primary">{task.assignee.name}</span>
                  </>
                ) : (
                  <span className="text-[13px] text-text-tertiary">Unassigned</span>
                )}
              </div>
            </div>
          </div>

          {/* Description */}
          <div>
            <label className="block text-[12px] font-medium text-text-tertiary mb-1.5">Description</label>
            <textarea
              className="w-full rounded-lg border border-border bg-bg px-3 py-2 text-[13px] text-text-primary outline-none focus:border-accent min-h-[120px] resize-y placeholder:text-text-tertiary"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              onBlur={() => saveField("description", description)}
              placeholder="Add a description..."
              rows={5}
            />
          </div>

          {/* Comments */}
          <div>
            <h3 className="text-[12px] font-medium text-text-tertiary mb-3">
              Activity ({comments.length})
            </h3>
            <div className="space-y-3">
              {comments.map((comment: Comment) => (
                <div key={comment.id} className="flex gap-2.5">
                  <Avatar
                    name={comment.author?.name || "User"}
                    avatarUrl={comment.author?.avatar_url}
                    size="sm"
                  />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-0.5">
                      <span className="text-[12px] font-medium text-text-primary">
                        {comment.author?.name || "User"}
                      </span>
                      <span className="text-[10px] text-text-tertiary">
                        {format(parseISO(comment.created_at), "MMM d, h:mm a")}
                      </span>
                    </div>
                    <p className="text-[13px] text-text-secondary leading-relaxed">{comment.content}</p>
                  </div>
                </div>
              ))}
              {comments.length === 0 && (
                <p className="text-[12px] text-text-tertiary text-center py-4">No comments yet</p>
              )}
            </div>
          </div>
        </div>

        {/* Comment Input */}
        <div className="border-t border-border px-6 py-4">
          <div className="flex items-center gap-2">
            {user && <Avatar name={user.name} size="sm" />}
            <div className="flex-1 flex items-center gap-2 rounded-lg border border-border bg-bg px-3 py-2">
              <input
                className="flex-1 bg-transparent text-[13px] text-text-primary outline-none placeholder:text-text-tertiary"
                placeholder="Add a comment..."
                value={commentText}
                onChange={(e) => setCommentText(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter" && !e.shiftKey && commentText.trim()) {
                    addCommentMutation.mutate(commentText.trim());
                  }
                }}
              />
              <button
                onClick={() => commentText.trim() && addCommentMutation.mutate(commentText.trim())}
                disabled={!commentText.trim()}
                className="rounded-md p-1 text-accent hover:bg-surface-hover disabled:opacity-30 disabled:hover:bg-transparent transition-colors"
              >
                <Send className="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>

        {/* Delete Confirmation */}
        {showDeleteConfirm && (
          <div className="absolute inset-0 z-50 flex items-center justify-center bg-black/50">
            <div className="mx-6 w-full max-w-sm rounded-xl bg-surface p-6 shadow-xl border border-border">
              <h3 className="text-lg font-semibold text-text-primary">Delete task?</h3>
              <p className="mt-1 text-[13px] text-text-secondary">
                This action cannot be undone. The task and all its comments will be permanently deleted.
              </p>
              <div className="mt-5 flex justify-end gap-2">
                <button
                  onClick={() => setShowDeleteConfirm(false)}
                  className="rounded-lg border border-border px-4 py-2 text-[13px] font-medium text-text-secondary hover:bg-surface-hover transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={() => deleteMutation.mutate()}
                  className="rounded-lg bg-danger px-4 py-2 text-[13px] font-medium text-white hover:opacity-90 transition-opacity"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </>
  );
}
