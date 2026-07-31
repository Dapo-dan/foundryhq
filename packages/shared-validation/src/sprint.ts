import { z } from 'zod'

export const sprintSchema = z
  .object({
    name: z.string().min(1, 'Sprint name is required'),
    startDate: z.string().min(1, 'Start date is required'),
    endDate: z.string().min(1, 'End date is required'),
  })
  .refine((data) => data.endDate >= data.startDate, {
    message: 'End date must be on or after the start date',
    path: ['endDate'],
  })

// Intersected with the required fields because z.infer alone doesn't mark
// them as required in this zod version — same workaround as
// project.ts's ProjectFormValues.
export type SprintFormValues = z.infer<typeof sprintSchema> & {
  name: string
  startDate: string
  endDate: string
}
