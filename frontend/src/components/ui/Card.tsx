import React from 'react'

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode
  hover?: boolean
  glass?: boolean
}

export const Card: React.FC<CardProps> = ({ 
  children, 
  hover = false, 
  glass = false, 
  className = '', 
  ...props 
}) => {
  const baseClasses = "rounded-xl border shadow-sm p-6"
  const glassClasses = glass 
    ? "glass-card" 
    : "bg-white border-primary-200/30"
  const hoverClasses = hover ? "card-hover" : ""
  
  return (
    <div 
      className={`${baseClasses} ${glassClasses} ${hoverClasses} ${className}`}
      {...props}
    >
      {children}
    </div>
  )
}

interface CardHeaderProps {
  children: React.ReactNode
  className?: string
}

export const CardHeader: React.FC<CardHeaderProps> = ({ children, className = '' }) => (
  <div className={`mb-4 ${className}`}>{children}</div>
)

interface CardTitleProps {
  children: React.ReactNode
  className?: string
}

export const CardTitle: React.FC<CardTitleProps> = ({ children, className = '' }) => (
  <h3 className={`text-lg font-semibold text-neutral-800 ${className}`}>
    {children}
  </h3>
)

interface CardContentProps {
  children: React.ReactNode
  className?: string
}

export const CardContent: React.FC<CardContentProps> = ({ children, className = '' }) => (
  <div className={`${className}`}>{children}</div>
)