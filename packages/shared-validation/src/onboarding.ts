import { z } from 'zod'

export const workspaceSchema = z.object({
  name: z.string().min(1, 'Workspace name is required'),
})

export type WorkspaceFormValues = Required<z.infer<typeof workspaceSchema>>

export const inviteEmailSchema = z.string().email('Enter a valid email address')
