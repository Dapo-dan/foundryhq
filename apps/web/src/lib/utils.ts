import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

// `cn` (short for "class names") is a tiny helper used everywhere in
// shadcn/ui components to combine Tailwind classes conditionally and safely.
// - clsx(inputs) lets you pass classes that depend on conditions, e.g.
//   cn("px-2", isActive && "bg-primary") — falsy values are just skipped.
// - twMerge(...) then resolves conflicts between Tailwind classes, e.g. if
//   both "px-2" and "px-4" end up in the list, it keeps only "px-4" instead
//   of letting both apply (which is normally what plain string concatenation
//   would do, and CSS would just let the last one in the stylesheet win).
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
