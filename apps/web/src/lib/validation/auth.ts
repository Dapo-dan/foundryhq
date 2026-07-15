import { z } from 'zod'

// This project's tsconfig doesn't enable `strict`/`strictNullChecks`, which
// breaks zod's usual required-vs-optional inference (every key ends up
// optional in `z.infer`, regardless of `.optional()`). `Required<...>`
// restores the correct shape without depending on that project-wide setting.

export const signInSchema = z.object({
  email: z.string().min(1, 'Email is required').email('Enter a valid email address'),
  password: z.string().min(1, 'Password is required'),
})

export type SignInFormValues = Required<z.infer<typeof signInSchema>>

export const signUpSchema = z.object({
  email: z.string().min(1, 'Email is required').email('Enter a valid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
})

export type SignUpFormValues = Required<z.infer<typeof signUpSchema>>
