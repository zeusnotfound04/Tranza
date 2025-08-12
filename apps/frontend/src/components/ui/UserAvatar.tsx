'use client';

import * as React from "react"
import { Avatar, AvatarImage, AvatarFallback } from '@tranza/ui/components/ui/avatar'
import { User } from '@/types/api'
import { cn } from '@/lib/utils'

interface UserAvatarProps {
  user?: User | null
  src?: string | null
  size?: 'sm' | 'md' | 'lg'
  className?: string
  showFallback?: boolean
}

export function UserAvatar({ 
  user, 
  src, 
  size = 'md',
  className,
  showFallback = true
}: UserAvatarProps) {
  const sizeClasses = {
    sm: 'h-6 w-6',
    md: 'h-8 w-8', 
    lg: 'h-12 w-12'
  }

  const textSizeClasses = {
    sm: 'text-xs',
    md: 'text-sm',
    lg: 'text-lg'
  }

  // Get avatar source from props or user object
  const avatarSrc = src || user?.avatar

  // Generate fallback text from user data
  const generateFallback = () => {
    if (user?.username) {
      // Use first letter of username
      return user.username.charAt(0).toUpperCase()
    }
    if (user?.email) {
      // Use first letter of email if no username
      return user.email.charAt(0).toUpperCase()
    }
    // Default fallback
    return 'U'
  }

  const fallbackText = generateFallback()

  return (
    <Avatar className={cn(sizeClasses[size], className)}>
      {avatarSrc && (
        <AvatarImage 
          src={avatarSrc} 
          alt={user?.username || 'User avatar'} 
        />
      )}
      {showFallback && (
        <AvatarFallback 
          className={cn(
            "bg-blue-500 text-white font-medium",
            textSizeClasses[size]
          )}
        >
          {fallbackText}
        </AvatarFallback>
      )}
    </Avatar>
  )
}
