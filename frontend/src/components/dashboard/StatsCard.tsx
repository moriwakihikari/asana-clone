"use client";

interface StatsCardProps {
  icon: React.ElementType;
  label: string;
  value: number;
  color: string;
}

export default function StatsCard({ icon: Icon, label, value, color }: StatsCardProps) {
  return (
    <div className="rounded-lg border border-border bg-surface p-4">
      <div className="flex items-center justify-between">
        <p className="text-[13px] text-text-secondary">{label}</p>
        <Icon size={16} style={{ color }} />
      </div>
      <p className="mt-2 text-2xl font-semibold text-text-primary tabular-nums">
        {value}
      </p>
    </div>
  );
}
