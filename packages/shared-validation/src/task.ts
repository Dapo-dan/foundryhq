import { z } from 'zod'

export const taskSchema = z.object({
  title: z.string().min(1, 'Task title is required'),
  projectId: z.string().min(1, 'Project is required'),
  assigneeId: z.string().optional(),
  sprintId: z.string().optional(),
  priority: z.enum(['urgent', 'high', 'medium', 'low']).optional(),
  // Raw form value from a number input — react-hook-form gives back a
  // string; the dialog converts it to a number before calling the API.
  storyPoints: z.string().optional(),
  dueDate: z.string().optional(),
})

// Intersected with the required fields because z.infer alone doesn't mark
// them as required in this zod version — same workaround as
// project.ts's ProjectFormValues.
export type TaskFormValues = z.infer<typeof taskSchema> & { title: string; projectId: string }
