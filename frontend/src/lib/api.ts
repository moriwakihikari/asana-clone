import axios from "axios";
import type {
  AuthResponse,
  User,
  Workspace,
  WorkspaceDetail,
  Project,
  Section,
  Task,
  Comment,
  Label,
  TaskStatus,
  TaskPriority,
} from "@/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

const api = axios.create({
  baseURL: API_URL,
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use((config) => {
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("access_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

// Auth
export const authApi = {
  register: (data: { name: string; email: string; password: string }) =>
    api.post<AuthResponse>("/auth/register", data).then((r) => r.data),
  login: (data: { email: string; password: string }) =>
    api.post<AuthResponse>("/auth/login", data).then((r) => r.data),
};

// Users
export const userApi = {
  getMe: () => api.get<User>("/users/me").then((r) => r.data),
  updateMe: (data: { name?: string; avatar_url?: string }) =>
    api.patch<User>("/users/me", data).then((r) => r.data),
  search: (q: string) =>
    api.get<User[]>("/users/search", { params: { q } }).then((r) => r.data),
};

// Workspaces
export const workspaceApi = {
  list: () => api.get<Workspace[]>("/workspaces").then((r) => r.data),
  get: (id: string) =>
    api.get<WorkspaceDetail>(`/workspaces/${id}`).then((r) => r.data),
  create: (data: { name: string; description?: string }) =>
    api.post<Workspace>("/workspaces", data).then((r) => r.data),
  update: (id: string, data: { name?: string; description?: string }) =>
    api.patch<Workspace>(`/workspaces/${id}`, data).then((r) => r.data),
  delete: (id: string) => api.delete(`/workspaces/${id}`),
  addMember: (id: string, data: { email: string; role?: string }) =>
    api.post(`/workspaces/${id}/members`, data),
  removeMember: (workspaceId: string, userId: string) =>
    api.delete(`/workspaces/${workspaceId}/members/${userId}`),
};

// Projects
export const projectApi = {
  list: (workspaceId: string) =>
    api.get<Project[]>(`/workspaces/${workspaceId}/projects`).then((r) => r.data),
  get: (workspaceId: string, id: string) =>
    api.get<Project>(`/workspaces/${workspaceId}/projects/${id}`).then((r) => r.data),
  create: (workspaceId: string, data: { name: string; description?: string; color?: string; view_type?: string }) =>
    api.post<Project>(`/workspaces/${workspaceId}/projects`, data).then((r) => r.data),
  update: (workspaceId: string, id: string, data: Partial<Project>) =>
    api.patch<Project>(`/workspaces/${workspaceId}/projects/${id}`, data).then((r) => r.data),
  archive: (workspaceId: string, id: string) =>
    api.post(`/workspaces/${workspaceId}/projects/${id}/archive`),
  delete: (workspaceId: string, id: string) =>
    api.delete(`/workspaces/${workspaceId}/projects/${id}`),
};

// Sections
export const sectionApi = {
  list: (projectId: string) =>
    api.get<Section[]>(`/projects/${projectId}/sections`).then((r) => r.data),
  create: (projectId: string, data: { name: string }) =>
    api.post<Section>(`/projects/${projectId}/sections`, data).then((r) => r.data),
  rename: (projectId: string, id: string, name: string) =>
    api.patch(`/projects/${projectId}/sections/${id}`, { name }),
  reorder: (projectId: string, data: { section_id: string; new_position: number }) =>
    api.post(`/projects/${projectId}/sections/reorder`, data),
  delete: (projectId: string, id: string) =>
    api.delete(`/projects/${projectId}/sections/${id}`),
};

// Tasks
export const taskApi = {
  list: (projectId: string, params?: Record<string, string>) =>
    api.get<Task[]>(`/projects/${projectId}/tasks`, { params }).then((r) => r.data),
  get: (taskId: string) =>
    api.get<Task>(`/tasks/${taskId}`).then((r) => r.data),
  create: (projectId: string, data: { title: string; section_id?: string; priority?: string; assignee_id?: string; due_date?: string }) =>
    api.post<Task>(`/projects/${projectId}/tasks`, data).then((r) => r.data),
  update: (taskId: string, data: Partial<Task>) =>
    api.patch<Task>(`/tasks/${taskId}`, data).then((r) => r.data),
  changeStatus: (taskId: string, status: TaskStatus) =>
    api.post(`/tasks/${taskId}/status`, { status }),
  move: (taskId: string, data: { section_id: string; position?: number }) =>
    api.post(`/tasks/${taskId}/move`, data),
  assign: (taskId: string, assignee_id: string | null) =>
    api.post(`/tasks/${taskId}/assign`, { assignee_id }),
  delete: (taskId: string) => api.delete(`/tasks/${taskId}`),
  getMyTasks: (workspaceId: string) =>
    api.get<Task[]>(`/workspaces/${workspaceId}/my-tasks`).then((r) => r.data),
};

// Comments
export const commentApi = {
  list: (taskId: string, params?: { page?: number; page_size?: number }) =>
    api.get<Comment[]>(`/tasks/${taskId}/comments`, { params }).then((r) => r.data),
  create: (taskId: string, content: string) =>
    api.post<Comment>(`/tasks/${taskId}/comments`, { content }).then((r) => r.data),
  update: (taskId: string, commentId: string, content: string) =>
    api.patch(`/tasks/${taskId}/comments/${commentId}`, { content }),
  delete: (taskId: string, commentId: string) =>
    api.delete(`/tasks/${taskId}/comments/${commentId}`),
};

// Labels
export const labelApi = {
  list: (workspaceId: string) =>
    api.get<Label[]>(`/workspaces/${workspaceId}/labels`).then((r) => r.data),
  create: (workspaceId: string, data: { name: string; color: string }) =>
    api.post<Label>(`/workspaces/${workspaceId}/labels`, data).then((r) => r.data),
  update: (workspaceId: string, id: string, data: { name?: string; color?: string }) =>
    api.patch(`/workspaces/${workspaceId}/labels/${id}`, data),
  delete: (workspaceId: string, id: string) =>
    api.delete(`/workspaces/${workspaceId}/labels/${id}`),
};

export default api;
