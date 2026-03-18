"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/providers";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await login(email, password);
      router.push("/dashboard");
    } catch {
      setError("Invalid email or password");
    } finally {
      setLoading(false);
    }
  };

  const fillDemo = () => {
    setEmail("demo@example.com");
    setPassword("password123");
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-bg px-4">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-semibold text-text-primary">Sign in</h1>
          <p className="mt-1 text-[13px] text-text-secondary">Welcome back to your workspace</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="rounded-md bg-danger/10 border border-danger/20 px-3 py-2 text-[13px] text-danger">
              {error}
            </div>
          )}

          <div>
            <label className="block text-[13px] font-medium text-text-primary mb-1.5">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="demo@example.com"
              className="h-9 w-full rounded-md border border-border-strong bg-surface px-3 text-[13px] text-text-primary placeholder:text-text-tertiary focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent/30"
              required
            />
          </div>

          <div>
            <label className="block text-[13px] font-medium text-text-primary mb-1.5">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="password123"
              className="h-9 w-full rounded-md border border-border-strong bg-surface px-3 text-[13px] text-text-primary placeholder:text-text-tertiary focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent/30"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="h-9 w-full rounded-md bg-accent text-[13px] font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
          >
            {loading ? "Signing in..." : "Sign in"}
          </button>
        </form>

        <div className="mt-4 text-center">
          <button
            type="button"
            onClick={fillDemo}
            className="text-[12px] text-accent hover:text-accent-hover font-medium"
          >
            Use demo account
          </button>
        </div>

        <p className="mt-6 text-center text-[12px] text-text-tertiary">
          Don&apos;t have an account?{" "}
          <Link href="/register" className="text-accent hover:text-accent-hover font-medium">
            Sign up
          </Link>
        </p>
      </div>
    </div>
  );
}
