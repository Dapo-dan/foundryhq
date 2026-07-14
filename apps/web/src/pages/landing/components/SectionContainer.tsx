import type { ComponentPropsWithoutRef, ElementType } from 'react'
import { cn } from '@/lib/utils'

interface SectionContainerProps {
  as?: ElementType
  className?: string
}

export function SectionContainer({
  as: Component = 'section',
  className,
  ...props
}: SectionContainerProps & ComponentPropsWithoutRef<'section'>) {
  return <Component className={cn('px-6 sm:px-10 lg:px-[100px]', className)} {...props} />
}
