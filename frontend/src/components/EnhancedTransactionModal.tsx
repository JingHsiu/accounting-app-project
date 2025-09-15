import React, { useState, useEffect } from 'react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui'
import { Card, CardContent } from '@/components/ui'
// Simple Label component since shadcn/ui not fully configured
const Label: React.FC<{ htmlFor?: string; children: React.ReactNode }> = ({ htmlFor, children }) => (
  <label htmlFor={htmlFor} className="block text-sm font-medium text-neutral-700 mb-2">
    {children}
  </label>
)
import { 
  ArrowUpRight, 
  ArrowDownLeft 
} from 'lucide-react'

// Enhanced Input Component with better styling
const EnhancedInput = React.forwardRef<HTMLInputElement, React.ComponentProps<"input">>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          "flex h-9 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground shadow-sm shadow-black/5 transition-shadow placeholder:text-muted-foreground/70 focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20 disabled:cursor-not-allowed disabled:opacity-50",
          className,
        )}
        ref={ref}
        {...props}
      />
    );
  },
);
EnhancedInput.displayName = "EnhancedInput";

// Currency Input Component with TWD step=1 support
interface CurrencyInputProps {
  value: string;
  onChange: (value: string) => void;
  currency: 'TWD' | 'USD';
  placeholder?: string;
  className?: string;
}

const CurrencyInput: React.FC<CurrencyInputProps> = ({
  value,
  onChange,
  currency,
  placeholder,
  className
}) => {
  const step = currency === 'TWD' ? '1' : '0.01';
  const symbol = currency === 'TWD' ? 'NT$' : '$';

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value;
    
    // Allow empty string
    if (inputValue === '') {
      onChange('');
      return;
    }

    // Validate based on currency
    if (currency === 'TWD') {
      // TWD: only integers
      if (/^\d+$/.test(inputValue)) {
        onChange(inputValue);
      }
    } else {
      // USD: allow decimals up to 2 places
      if (/^\d*\.?\d{0,2}$/.test(inputValue)) {
        onChange(inputValue);
      }
    }
  };

  return (
    <div className="relative">
      <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">
        {symbol}
      </span>
      <EnhancedInput
        type="text"
        value={value}
        onChange={handleChange}
        placeholder={placeholder}
        className={cn("pl-10", className)}
        step={step}
        min="0"
      />
    </div>
  );
};

// Simple Select Component (fallback if shadcn not available)
interface SimpleSelectProps {
  value: string;
  onValueChange: (value: string) => void;
  placeholder: string;
  children: React.ReactNode;
  className?: string;
}

const SimpleSelect: React.FC<SimpleSelectProps> = ({
  value,
  onValueChange,
  placeholder,
  children,
  className
}) => {
  return (
    <select
      value={value}
      onChange={(e) => onValueChange(e.target.value)}
      className={cn(
        "flex h-9 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
        className
      )}
    >
      <option value="" disabled>{placeholder}</option>
      {children}
    </select>
  );
};

// Transaction Form Data Interface
interface TransactionFormData {
  walletID: string;
  categoryID: string;
  amount: string;
  currency: 'TWD' | 'USD';
  description: string;
  type: 'income' | 'expense';
}

// Enhanced Transaction Modal Props
interface EnhancedTransactionModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: {
    wallet_id: string;
    subcategory_id: string;
    amount: number;
    currency: string;
    description?: string;
    date: string;
  }) => void;
  wallets?: Array<{ id: string; name: string; }>;
  categories?: Array<{ id: string; name: string; }>;
  type: 'income' | 'expense';
}

export const EnhancedTransactionModal: React.FC<EnhancedTransactionModalProps> = ({
  open,
  onOpenChange,
  onSubmit,
  wallets = [],
  categories = [],
  type
}) => {
  const [step, setStep] = useState(1);
  const [formData, setFormData] = useState<TransactionFormData>({
    walletID: '',
    categoryID: '',
    amount: '',
    currency: 'TWD',
    description: '',
    type: type
  });

  const resetForm = () => {
    setStep(1);
    setFormData({
      walletID: '',
      categoryID: '',
      amount: '',
      currency: 'TWD',
      description: '',
      type: type
    });
  };

  const handleNext = () => {
    if (step < 3) setStep(step + 1);
  };

  const handleBack = () => {
    if (step > 1) setStep(step - 1);
  };

  const handleSubmit = () => {
    const amount = parseFloat(formData.amount);
    if (isNaN(amount) || amount <= 0) return;

    // Convert to backend format based on currency
    const backendAmount = formData.currency === 'TWD' ? amount : Math.round(amount * 100);

    onSubmit({
      wallet_id: formData.walletID,
      subcategory_id: formData.categoryID,
      amount: backendAmount,
      currency: formData.currency,
      description: formData.description.trim() || undefined,
      date: new Date().toISOString()
    });

    resetForm();
    onOpenChange(false);
  };

  const canProceedStep1 = formData.walletID && formData.categoryID;
  const canProceedStep2 = formData.amount && parseFloat(formData.amount) > 0;

  useEffect(() => {
    if (!open) {
      resetForm();
    }
  }, [open]);

  useEffect(() => {
    setFormData(prev => ({ ...prev, type: type }));
  }, [type]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-xl max-w-md w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex items-center gap-2 mb-6">
            <div className="w-8 h-8 bg-primary/10 rounded-lg flex items-center justify-center">
              {type === 'expense' ? (
                <ArrowDownLeft className="w-4 h-4 text-red-600" />
              ) : (
                <ArrowUpRight className="w-4 h-4 text-green-600" />
              )}
            </div>
            <h2 className="text-xl font-bold">Add {type === 'expense' ? 'Expense' : 'Income'}</h2>
          </div>

        <div className="space-y-6">
          {/* Progress Indicator */}
          <div className="flex items-center justify-center space-x-2">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className={cn(
                  "w-2 h-2 rounded-full transition-colors",
                  i <= step ? "bg-primary" : "bg-muted"
                )}
              />
            ))}
          </div>

          {/* Step 1: Wallet & Category */}
          {step === 1 && (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="wallet">Wallet</Label>
                <SimpleSelect
                  value={formData.walletID}
                  onValueChange={(value) => setFormData(prev => ({ ...prev, walletID: value }))}
                  placeholder="Select wallet"
                >
                  {wallets.map((wallet) => (
                    <option key={wallet.id} value={wallet.id}>
                      {wallet.name}
                    </option>
                  ))}
                </SimpleSelect>
              </div>

              <div className="space-y-2">
                <Label htmlFor="category">Category</Label>
                <SimpleSelect
                  value={formData.categoryID}
                  onValueChange={(value) => setFormData(prev => ({ ...prev, categoryID: value }))}
                  placeholder="Select category"
                >
                  {categories.map((category) => (
                    <option key={category.id} value={category.id}>
                      {category.name}
                    </option>
                  ))}
                </SimpleSelect>
              </div>
            </div>
          )}

          {/* Step 2: Amount & Currency */}
          {step === 2 && (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label>Currency</Label>
                <div className="flex gap-2">
                  <Button
                    type="button"
                    variant={formData.currency === 'TWD' ? 'primary' : 'outline'}
                    onClick={() => setFormData(prev => ({ ...prev, currency: 'TWD', amount: '' }))}
                    className="flex-1"
                  >
                    TWD (1 unit)
                  </Button>
                  <Button
                    type="button"
                    variant={formData.currency === 'USD' ? 'primary' : 'outline'}
                    onClick={() => setFormData(prev => ({ ...prev, currency: 'USD', amount: '' }))}
                    className="flex-1"
                  >
                    USD (0.01)
                  </Button>
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="amount">Amount</Label>
                <CurrencyInput
                  value={formData.amount}
                  onChange={(value) => setFormData(prev => ({ ...prev, amount: value }))}
                  currency={formData.currency}
                  placeholder={formData.currency === 'TWD' ? 'Enter whole number' : 'Enter amount (e.g., 10.50)'}
                />
                <p className="text-xs text-muted-foreground">
                  {formData.currency === 'TWD' 
                    ? 'TWD amounts must be whole numbers (1 unit increments)'
                    : 'USD amounts can include cents (0.01 increments)'
                  }
                </p>
              </div>
            </div>
          )}

          {/* Step 3: Description (Optional) */}
          {step === 3 && (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="description">Description (Optional)</Label>
                <textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                  placeholder="Add a note about this transaction..."
                  className="flex min-h-[80px] w-full rounded-lg border border-input bg-background px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                />
              </div>

              {/* Summary */}
              <Card>
                <CardContent className="p-4">
                  <h4 className="font-medium mb-2">Transaction Summary</h4>
                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Wallet:</span>
                      <span>{wallets.find(w => w.id === formData.walletID)?.name}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Category:</span>
                      <span>{categories.find(c => c.id === formData.categoryID)?.name}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Amount:</span>
                      <span className={type === 'expense' ? 'text-red-600' : 'text-green-600'}>
                        {formData.currency === 'TWD' ? 'NT$' : '$'}{formData.amount}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Type:</span>
                      <span className={cn(
                        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
                        type === 'expense' 
                          ? 'bg-red-100 text-red-800'
                          : 'bg-green-100 text-green-800'
                      )}>
                        {type}
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </div>

          <div className="flex justify-between pt-6">
            <div className="flex gap-2">
              {step > 1 && (
                <Button type="button" variant="outline" onClick={handleBack}>
                  Back
                </Button>
              )}
            </div>
            <div className="flex gap-2">
              <Button type="button" variant="ghost" onClick={() => onOpenChange(false)}>
                Cancel
              </Button>
              {step < 3 ? (
                <Button 
                  type="button" 
                  onClick={handleNext}
                  disabled={step === 1 ? !canProceedStep1 : !canProceedStep2}
                >
                  Next
                </Button>
              ) : (
                <Button type="button" onClick={handleSubmit}>
                  Add {type === 'expense' ? 'Expense' : 'Income'}
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default EnhancedTransactionModal;