"use client";

import clsx from "clsx";

interface AvatarProps {
  name: string;
  avatarUrl?: string;
  size?: "sm" | "md" | "lg";
}

const sizes = {
  sm: "h-5 w-5 text-[9px]",
  md: "h-7 w-7 text-[11px]",
  lg: "h-9 w-9 text-[13px]",
};

const colors = [
  "#4573d2", "#e8615a", "#5da283", "#f1914e",
  "#a276d4", "#d4748a", "#6caad9", "#d9a76c",
];

function getColor(name: string) {
  let hash = 0;
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash);
  return colors[Math.abs(hash) % colors.length];
}

export default function Avatar({ name, avatarUrl, size = "sm" }: AvatarProps) {
  if (avatarUrl) {
    return (
      <img
        src={avatarUrl}
        alt={name}
        className={clsx("rounded-full object-cover", sizes[size])}
      />
    );
  }

  return (
    <div
      className={clsx(
        "flex items-center justify-center rounded-full font-medium text-white shrink-0",
        sizes[size]
      )}
      style={{ backgroundColor: getColor(name) }}
    >
      {name.charAt(0).toUpperCase()}
    </div>
  );
}
