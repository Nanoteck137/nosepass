import type { CheckedValue } from "$lib/types";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function capitilize(s: string) {
  if (s.length === 0) return "";
  return s[0].toUpperCase() + s.substring(1);
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isRoleAdmin(role: string) {
  switch (role) {
    case "super_user":
    case "admin":
      return true;
    default:
      return false;
  }
}

export function convertValue<T>(val: CheckedValue<T>): T | undefined {
  if (val.checked) {
    return val.value;
  }

  return undefined;
}

export function shouldDisplayHour(duration: number) {
  return duration >= 1 * 60 * 60;
}

export function formatTime(s: number, displayHour: boolean) {
  const hour = Math.floor(s / 60 / 60);
  const min = Math.floor(s / 60);
  const sec = Math.floor(s % 60);

  if (displayHour) {
    return `${hour}:${min.toString().padStart(2, "0")}:${sec.toString().padStart(2, "0")}`;
  } else {
    return `${min}:${sec.toString().padStart(2, "0")}`;
  }
}
