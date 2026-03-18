"use client";

import { useState } from "react";
import { DragDropContext, Droppable, Draggable, type DropResult } from "@hello-pangea/dnd";
import { Plus, MoreHorizontal } from "lucide-react";
import clsx from "clsx";
import type { Section, Task } from "@/types";
import TaskCard from "./TaskCard";

interface BoardViewProps {
  sections: Section[];
  tasksBySection: Record<string, Task[]>;
  onMoveTask: (taskId: string, sectionId: string, position: number) => void;
  onClickTask: (task: Task) => void;
  onAddTask: (sectionId: string) => void;
  onAddSection: () => void;
  onRenameSection?: (sectionId: string, name: string) => void;
  onDeleteSection?: (sectionId: string) => void;
}

export default function BoardView({
  sections,
  tasksBySection,
  onMoveTask,
  onClickTask,
  onAddTask,
  onAddSection,
  onRenameSection,
  onDeleteSection,
}: BoardViewProps) {
  const [editingSectionId, setEditingSectionId] = useState<string | null>(null);
  const [editingSectionName, setEditingSectionName] = useState("");
  const [menuSectionId, setMenuSectionId] = useState<string | null>(null);
  const [inlineAddSectionId, setInlineAddSectionId] = useState<string | null>(null);
  const [inlineTitle, setInlineTitle] = useState("");

  const handleDragEnd = (result: DropResult) => {
    if (!result.destination) return;
    const { draggableId, destination } = result;
    onMoveTask(draggableId, destination.droppableId, destination.index);
  };

  const handleSectionRename = (sectionId: string) => {
    if (editingSectionName.trim() && onRenameSection) {
      onRenameSection(sectionId, editingSectionName.trim());
    }
    setEditingSectionId(null);
  };

  return (
    <DragDropContext onDragEnd={handleDragEnd}>
      <div className="flex gap-3 overflow-x-auto pb-4 px-1 pt-1">
        {sections.map((section) => {
          const tasks = tasksBySection[section.id] || [];
          return (
            <div
              key={section.id}
              className="flex w-[280px] shrink-0 flex-col rounded-lg bg-bg-raised"
            >
              {/* Column Header */}
              <div className="flex items-center justify-between px-3 py-2.5 rounded-t-lg bg-surface-hover">
                <div className="flex items-center gap-2 min-w-0">
                  {editingSectionId === section.id ? (
                    <input
                      autoFocus
                      className="text-[13px] font-bold text-text-primary bg-surface rounded px-1.5 py-0.5 outline-none border border-accent"
                      value={editingSectionName}
                      onChange={(e) => setEditingSectionName(e.target.value)}
                      onBlur={() => handleSectionRename(section.id)}
                      onKeyDown={(e) => {
                        if (e.key === "Enter") handleSectionRename(section.id);
                        if (e.key === "Escape") setEditingSectionId(null);
                      }}
                    />
                  ) : (
                    <h3
                      className="text-[13px] font-bold text-text-primary truncate cursor-pointer"
                      onDoubleClick={() => {
                        setEditingSectionId(section.id);
                        setEditingSectionName(section.name);
                      }}
                    >
                      {section.name}
                    </h3>
                  )}
                  <span className="text-[12px] font-medium text-text-tertiary">
                    {tasks.length}
                  </span>
                </div>
                <div className="flex items-center gap-0.5">
                  <button
                    onClick={() => onAddTask(section.id)}
                    className="rounded p-1 text-text-tertiary hover:bg-surface-active hover:text-text-secondary transition-colors"
                    title="Add task"
                  >
                    <Plus className="h-3.5 w-3.5" />
                  </button>
                  <div className="relative">
                    <button
                      onClick={() =>
                        setMenuSectionId(menuSectionId === section.id ? null : section.id)
                      }
                      className="rounded p-1 text-text-tertiary hover:bg-surface-active hover:text-text-secondary transition-colors"
                    >
                      <MoreHorizontal className="h-3.5 w-3.5" />
                    </button>
                    {menuSectionId === section.id && (
                      <>
                        <div
                          className="fixed inset-0 z-10"
                          onClick={() => setMenuSectionId(null)}
                        />
                        <div className="absolute right-0 top-full z-20 mt-1 w-36 rounded-lg border border-border-strong bg-surface py-1 shadow-xl">
                          <button
                            className="w-full px-3 py-1.5 text-left text-[13px] text-text-secondary hover:bg-surface-hover transition-colors"
                            onClick={() => {
                              setEditingSectionId(section.id);
                              setEditingSectionName(section.name);
                              setMenuSectionId(null);
                            }}
                          >
                            Rename
                          </button>
                          <button
                            className="w-full px-3 py-1.5 text-left text-[13px] text-danger hover:bg-surface-hover transition-colors"
                            onClick={() => {
                              onDeleteSection?.(section.id);
                              setMenuSectionId(null);
                            }}
                          >
                            Delete
                          </button>
                        </div>
                      </>
                    )}
                  </div>
                </div>
              </div>

              {/* Cards */}
              <Droppable droppableId={section.id}>
                {(provided, snapshot) => (
                  <div
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                    className={clsx(
                      "flex-1 space-y-1.5 p-1.5 min-h-[60px] transition-colors duration-200",
                      snapshot.isDraggingOver && "bg-surface/40 rounded-b-lg"
                    )}
                  >
                    {tasks.map((task, index) => (
                      <Draggable key={task.id} draggableId={task.id} index={index}>
                        {(provided, snapshot) => (
                          <div
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                            {...provided.dragHandleProps}
                          >
                            <TaskCard
                              task={task}
                              onClick={onClickTask}
                              isDragging={snapshot.isDragging}
                            />
                          </div>
                        )}
                      </Draggable>
                    ))}
                    {provided.placeholder}

                    {/* Inline add task */}
                    {inlineAddSectionId === section.id ? (
                      <div className="rounded-lg border border-accent bg-surface p-2">
                        <input
                          autoFocus
                          className="w-full text-[13px] text-text-primary bg-transparent outline-none placeholder:text-text-tertiary"
                          placeholder="Task name..."
                          value={inlineTitle}
                          onChange={(e) => setInlineTitle(e.target.value)}
                          onKeyDown={(e) => {
                            if (e.key === "Enter" && inlineTitle.trim()) {
                              onAddTask(section.id);
                              setInlineTitle("");
                              setInlineAddSectionId(null);
                            }
                            if (e.key === "Escape") {
                              setInlineAddSectionId(null);
                              setInlineTitle("");
                            }
                          }}
                          onBlur={() => {
                            setInlineAddSectionId(null);
                            setInlineTitle("");
                          }}
                        />
                      </div>
                    ) : (
                      <button
                        onClick={() => setInlineAddSectionId(section.id)}
                        className="w-full flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-[12px] text-text-tertiary hover:bg-surface-hover hover:text-text-secondary transition-colors"
                      >
                        <Plus className="h-3 w-3" />
                        Add task
                      </button>
                    )}
                  </div>
                )}
              </Droppable>
            </div>
          );
        })}

        {/* Add Section */}
        <button
          onClick={onAddSection}
          className="flex h-10 w-[280px] shrink-0 items-center justify-center gap-1.5 rounded-lg border-2 border-dashed border-border text-[13px] text-text-tertiary hover:border-border-strong hover:text-text-secondary hover:bg-bg-raised transition-colors"
        >
          <Plus className="h-4 w-4" />
          Add Section
        </button>
      </div>
    </DragDropContext>
  );
}
