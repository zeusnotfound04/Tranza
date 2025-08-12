import * as React from "react"
import { cn } from "@/lib/utils"
import { 
  Wallet, 
  ArrowUpRight, 
  ArrowDownLeft, 
  CreditCard,
  Repeat,
  Zap,
  DollarSign,
  TrendingUp,
  Bot
} from "lucide-react"

export interface TransactionBadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  type: string
  status?: string
  size?: 'sm' | 'md' | 'lg'
  showIcon?: boolean
  showStatus?: boolean
  variant?: 'default' | 'outline' | 'minimal'
}

const TransactionBadge = React.forwardRef<HTMLDivElement, TransactionBadgeProps>(
  ({ 
    className, 
    type, 
    status = 'completed', 
    size = 'md', 
    showIcon = true, 
    showStatus = true, 
    variant = 'default',
    ...props 
  }, ref) => {
    
    const getTransactionIcon = (transactionType: string) => {
      const iconSize = size === 'sm' ? 'w-4 h-4' : size === 'lg' ? 'w-6 h-6' : 'w-5 h-5';
      
      switch (transactionType.toLowerCase()) {
        case 'load_money':
        case 'deposit':
          return <Wallet className={`${iconSize} text-emerald-600`} />;
        case 'send_money':
        case 'transfer':
          return <ArrowUpRight className={`${iconSize} text-blue-600`} />;
        case 'receive_money':
        case 'received':
          return <ArrowDownLeft className={`${iconSize} text-green-600`} />;
        case 'ai_agent':
        case 'ai_payment':
          return <Bot className={`${iconSize} text-purple-600`} />;
        case 'refund':
          return <ArrowDownLeft className={`${iconSize} text-orange-600`} />;
        case 'payment':
          return <CreditCard className={`${iconSize} text-indigo-600`} />;
        case 'subscription':
          return <Repeat className={`${iconSize} text-cyan-600`} />;
        case 'reward':
        case 'bonus':
          return <Zap className={`${iconSize} text-yellow-600`} />;
        case 'fee':
        case 'charge':
          return <DollarSign className={`${iconSize} text-red-600`} />;
        case 'investment':
          return <TrendingUp className={`${iconSize} text-blue-700`} />;
        default:
          return <Wallet className={`${iconSize} text-slate-600`} />;
      }
    };

    const getStatusColor = (transactionStatus: string) => {
      switch (transactionStatus.toLowerCase()) {
        case 'completed':
        case 'success':
          return 'bg-emerald-100 text-emerald-800';
        case 'pending':
        case 'processing':
          return 'bg-amber-100 text-amber-800';
        case 'failed':
        case 'cancelled':
        case 'error':
          return 'bg-red-100 text-red-800';
        case 'refunded':
          return 'bg-orange-100 text-orange-800';
        default:
          return 'bg-slate-100 text-slate-800';
      }
    };

    const getTransactionTypeColor = (transactionType: string) => {
      switch (transactionType.toLowerCase()) {
        case 'load_money':
        case 'deposit':
        case 'receive_money':
        case 'received':
          return 'bg-emerald-50';
        case 'send_money':
        case 'transfer':
          return 'bg-blue-50';
        case 'ai_agent':
        case 'ai_payment':
          return 'bg-purple-50';
        case 'refund':
          return 'bg-orange-50';
        case 'payment':
          return 'bg-indigo-50';
        case 'subscription':
          return 'bg-cyan-50';
        case 'reward':
        case 'bonus':
          return 'bg-yellow-50';
        case 'fee':
        case 'charge':
          return 'bg-red-50';
        case 'investment':
          return 'bg-blue-50';
        default:
          return 'bg-slate-50';
      }
    };

    const sizeClasses = {
      sm: 'px-2 py-1 text-xs gap-1.5',
      md: 'px-3 py-1.5 text-sm gap-2',
      lg: 'px-4 py-2 text-base gap-2.5'
    };

    const baseClasses = "inline-flex items-center rounded-full font-medium transition-all duration-300 hover:scale-105";
    
    const variantClasses = {
      default: `${getTransactionTypeColor(type)}`,
      outline: "bg-white text-slate-700 hover:bg-slate-50",
      minimal: "bg-transparent text-slate-600 hover:bg-slate-100"
    };

    return (
      <div
        ref={ref}
        className={cn(
          baseClasses,
          sizeClasses[size],
          variantClasses[variant],
          className
        )}
        {...props}
      >
        {showIcon && (
          <span className="flex-shrink-0">
            {getTransactionIcon(type)}
          </span>
        )}
        
        {showStatus && (
          <span 
            className={cn(
              "px-2 py-0.5 rounded-full text-xs font-semibold",
              getStatusColor(status)
            )}
          >
            {status.charAt(0).toUpperCase() + status.slice(1)}
          </span>
        )}
        
        {!showStatus && (
          <span className="capitalize font-semibold">
            {type.replace('_', ' ')}
          </span>
        )}
      </div>
    )
  }
)

TransactionBadge.displayName = "TransactionBadge"

export { TransactionBadge }
