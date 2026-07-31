import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Task } from '@foundryhq/shared-types'
import { updateTask, type UpdateTaskInput } from '@/services/tasks'
import { useWorkspaceStore } from '@/store/slices/workspace'

type UpdateTaskVariables = UpdateTaskInput & { taskId: string }

function applyPatch(task: Task, variables: UpdateTaskVariables): Task {
  return {
    ...task,
    ...(variables.title !== undefined && { title: variables.title }),
    ...(variables.status !== undefined && { status: variables.status }),
    ...(variables.projectId !== undefined && { projectId: variables.projectId }),
    ...(variables.assigneeId !== undefined && { assigneeId: variables.assigneeId }),
    ...(variables.sprintId !== undefined && { sprintId: variables.sprintId }),
    ...(variables.priority !== undefined && { priority: variables.priority }),
    ...(variables.storyPoints !== undefined && { storyPoints: variables.storyPoints }),
    ...(variables.dueDate !== undefined && { dueDate: variables.dueDate }),
    // Clear* flags are checked last so they win over their corresponding
    // field above if both were somehow set on the same call.
    ...(variables.clearAssignee && { assigneeId: null }),
    ...(variables.clearSprint && { sprintId: null }),
    ...(variables.clearStoryPoints && { storyPoints: null }),
    ...(variables.clearDueDate && { dueDate: null }),
  }
}

// The first mutation in the app implementing ADR-0005 (optimistic UI with
// TanStack Query) — this is the pattern later drag/status interactions
// (deal-pipeline stages, OKR check-ins) are expected to reuse. Used both by
// the Kanban board's drag-and-drop (status changes) and by regular task
// edits, so every task mutation gets the same instant-feedback/rollback
// behavior for free.
export function useUpdateTask() {
  const queryClient = useQueryClient()
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  const listKey = ['tasks', workspaceId]

  return useMutation({
    mutationFn: ({ taskId, ...input }: UpdateTaskVariables) => updateTask(taskId, input),
    onMutate: async (variables) => {
      // Cancel any in-flight refetch of the tasks list so it can't land
      // after our optimistic write and clobber it — matters most for rapid
      // successive drags on the same task (ADR-0005).
      await queryClient.cancelQueries({ queryKey: listKey })

      const previous = queryClient.getQueriesData<Task[]>({ queryKey: listKey })
      queryClient.setQueriesData<Task[]>({ queryKey: listKey }, (tasks) =>
        tasks?.map((task) => (task.id === variables.taskId ? applyPatch(task, variables) : task))
      )
      return { previous }
    },
    onError: (_err, _variables, context) => {
      context?.previous?.forEach(([key, data]) => {
        queryClient.setQueryData(key, data)
      })
      toast.error('Could not update the task — your change was reverted.')
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: listKey })
    },
  })
}
