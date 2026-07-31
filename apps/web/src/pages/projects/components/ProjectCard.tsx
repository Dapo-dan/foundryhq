import type { Project } from '@foundryhq/shared-types'

interface ProjectCardProps {
  project: Project
}

// Extracted to its own file to match the pattern SprintCard already
// established, rather than the list item being inlined in the page.
export function ProjectCard({ project }: ProjectCardProps) {
  return (
    <div className="rounded-lg border border-border p-4">
      <h3 className="font-medium">{project.name}</h3>
      {project.description && (
        <p className="text-sm text-muted-foreground">{project.description}</p>
      )}
    </div>
  )
}
