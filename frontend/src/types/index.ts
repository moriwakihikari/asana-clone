export interface User {
  id: string;
  name: string;
  email: string;
  avatar_url?: string;
}

export interface Workspace {
  id: string;
  name: string;
  description?: string;
  owner_id: string;
  created_at: string;
}

export interface WorkspaceDetail extends Workspace {
  members: WorkspaceMember[];
}

export interface WorkspaceMember {
  user_id: string;
  workspace_id: string;
  role: "owner" | "admin" | "member";
  joined_at: string;
  user?: User;
}

export interface Project {
  id: string;
  workspace_id: string;
  name: string;
  description?: string;
  color: string;
  view_type: "list" | "board" | "calendar" | "timeline";
  is_archived: boolean;
  created_at: string;
  updated_at: string;
}

export interface Section {
  id: string;
  project_id: string;
  name: string;
  position: number;
  created_at: string;
}

export type TaskStatus = "todo" | "in_progress" | "in_review" | "done" | "cancelled";
export type TaskPriority = "none" | "low" | "medium" | "high" | "urgent";

export interface Task {
  id: string;
  project_id: string;
  section_id?: string;
  assignee_id?: string;
  assignee?: User;
  title: string;
  description?: string;
  status: TaskStatus;
  priority: TaskPriority;
  due_date?: string;
  position: number;
  labels?: Label[];
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: string;
  task_id: string;
  user_id: string;
  author?: User;
  content: string;
  is_edited: boolean;
  created_at: string;
  updated_at: string;
}

export interface Label {
  id: string;
  workspace_id: string;
  name: string;
  color: string;
  created_at: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export interface AuthResponse {
  tokens: TokenPair;
  user: User;
}

export interface ApiError {
  code: string;
  message: string;
  field?: string;
}
