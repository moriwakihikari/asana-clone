"use client";

import { Inbox as InboxIcon } from "lucide-react";

export default function InboxPage() {
  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <h1 className="text-xl font-bold text-text-primary mb-6">Inbox</h1>

      <div className="rounded-lg border border-border bg-surface py-16 text-center">
        <InboxIcon size={32} className="mx-auto text-text-tertiary mb-3" />
        <p className="text-[14px] text-text-secondary font-medium">You&apos;re all caught up!</p>
        <p className="text-[13px] text-text-tertiary mt-1">
          Notifications about your tasks and projects will appear here.
        </p>
      </div>
    </div>
  );
}
