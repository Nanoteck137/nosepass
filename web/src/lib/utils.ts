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
