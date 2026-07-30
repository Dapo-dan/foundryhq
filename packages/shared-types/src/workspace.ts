export interface Workspace {
  id: string
  name: string
  slug: string
  logoUrl: string
  createdAt: string
  updatedAt: string
}

export interface WorkspaceMember {
  id: string
  userId: string
  email: string
  role: string
  invitedAt: string
  joinedAt: string | null
}
