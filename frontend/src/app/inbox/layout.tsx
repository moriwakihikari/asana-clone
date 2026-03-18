"use client";

import { useAuth } from "@/lib/providers";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import AppLayout from "@/components/layout/AppLayout";

export default function InboxLayout({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !user) router.push("/login");
  }, [user, isLoading, router]);

  if (isLoading || !user) {
    return (
      <div className="flex h-screen items-center justify-center bg-bg">
        <div className="h-6 w-6 animate-spin rounded-full border-2 border-accent border-t-transparent" />
      </div>
    );
  }

  return <AppLayout>{children}</AppLayout>;
}
