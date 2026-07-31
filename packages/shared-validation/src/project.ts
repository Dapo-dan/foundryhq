import { z } from 'zod'

export const projectSchema = z.object({
  name: z.string().min(1, 'Project name is required'),
  description: z.string().optional(),
})

// Intersected with `{ name: string }` because z.infer alone doesn't mark
// `name` as required in this zod version — same workaround as
// onboarding.ts's WorkspaceFormValues (there via `Required<>`, but that
// would also force `description` required, which it shouldn't be).
export type ProjectFormValues = z.infer<typeof projectSchema> & { name: string }
