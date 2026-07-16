import { Users, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function ContactsPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Contacts</h1>
          <p className="text-sm text-muted-foreground">
            Keep track of the people you're building relationships with.
          </p>
        </div>
        <Button>
          <Plus size={20} />
          New contact
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <Users size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">No contacts yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Add a contact to start logging activity and deals with them.
        </p>
        <Button variant="outline">Add contact</Button>
      </div>
    </div>
  )
}
