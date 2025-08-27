import { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import { 
  apiClient, 
  TransferValidationRequest, 
  CreateTransferRequest,
  validateAmount, 
  validateUPI, 
  validatePhone, 
  formatCurrency 
} from '../services/api';
import { Send, ArrowRight, Check, X, AlertCircle, Smartphone, CreditCard } from 'lucide-react';

type RecipientType = 'upi' | 'phone';

interface TransferFormData {
  amount: string;
  recipientType: RecipientType;
  recipientValue: string;
  recipientName: string;
  description: string;
}

export default function MoneyTransfer() {
  const { user } = useAuth();
  const [formData, setFormData] = useState<TransferFormData>({
    amount: '',
    recipientType: 'upi',
    recipientValue: '',
    recipientName: '',
    description: '',
  });
  
  const [validation, setValidation] = useState<any>(null);
  const [validating, setValidating] = useState(false);
  const [transferring, setTransferring] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [step, setStep] = useState<'form' | 'confirm' | 'success'>('form');
  const [transferResult, setTransferResult] = useState<any>(null);

  // Validate form inputs
  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.amount) {
      newErrors.amount = 'Amount is required';
    } else if (!validateAmount(formData.amount)) {
      newErrors.amount = 'Amount must be between ₹1 and ₹1,00,000';
    }

    if (!formData.recipientValue) {
      newErrors.recipientValue = `${formData.recipientType === 'upi' ? 'UPI ID' : 'Phone number'} is required`;
    } else if (formData.recipientType === 'upi' && !validateUPI(formData.recipientValue)) {
      newErrors.recipientValue = 'Please enter a valid UPI ID (e.g., user@paytm)';
    } else if (formData.recipientType === 'phone' && !validatePhone(formData.recipientValue)) {
      newErrors.recipientValue = 'Please enter a valid 10-digit mobile number';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Validate transfer with backend
  const handleValidateTransfer = async () => {
    if (!validateForm()) return;

    setValidating(true);
    setValidation(null);

    try {
      const request: TransferValidationRequest = {
        amount: formData.amount,
        recipient_type: formData.recipientType,
        recipient_value: formData.recipientValue,
      };

      const response = await apiClient.validateTransfer(request);
      
      if (response.success && response.data) {
        setValidation(response.data);
        if (response.data.valid) {
          setStep('confirm');
        }
      } else {
        setErrors({ general: response.error || 'Validation failed' });
      }
    } catch (error) {
      setErrors({ general: 'Network error. Please try again.' });
    } finally {
      setValidating(false);
    }
  };

  // Create transfer
  const handleCreateTransfer = async () => {
    setTransferring(true);

    try {
      const request: CreateTransferRequest = {
        amount: formData.amount,
        recipient_type: formData.recipientType,
        recipient_value: formData.recipientValue,
        recipient_name: formData.recipientName,
        description: formData.description,
      };

      const response = await apiClient.createTransfer(request);
      
      if (response.success && response.data) {
        setTransferResult(response.data);
        setStep('success');
      } else {
        setErrors({ general: response.error || 'Transfer failed' });
        setStep('form');
      }
    } catch (error) {
      setErrors({ general: 'Network error. Please try again.' });
      setStep('form');
    } finally {
      setTransferring(false);
    }
  };

  // Reset form
  const resetForm = () => {
    setFormData({
      amount: '',
      recipientType: 'upi',
      recipientValue: '',
      recipientName: '',
      description: '',
    });
    setValidation(null);
    setErrors({});
    setStep('form');
    setTransferResult(null);
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Send className="w-8 h-8 text-blue-600" />
          Send Money
        </h1>
        <p className="text-gray-600">Transfer money to UPI ID or phone number</p>
      </div>

      {/* Step Indicator */}
      <div className="flex items-center justify-center mb-8">
        <StepIndicator 
          steps={['Details', 'Confirm', 'Success']} 
          currentStep={step === 'form' ? 0 : step === 'confirm' ? 1 : 2} 
        />
      </div>

      {/* Form Step */}
      {step === 'form' && (
        <div className="bg-black rounded-lg border border-gray-200 p-6">
          <form onSubmit={(e) => { e.preventDefault(); handleValidateTransfer(); }}>
            {/* Amount Input */}
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Amount to Send
              </label>
              <div className="relative">
                <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">₹</span>
                <input
                  type="number"
                  value={formData.amount}
                  onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                  className={`w-full pl-8 pr-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 ${
                    errors.amount ? 'border-red-500' : 'border-gray-300'
                  }`}
                  placeholder="Enter amount"
                  min="1"
                  max="100000"
                />
              </div>
              {errors.amount && <p className="text-red-500 text-sm mt-1">{errors.amount}</p>}
            </div>

            {/* Recipient Type Selection */}
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-3">
                Send To
              </label>
              <div className="grid grid-cols-2 gap-4">
                <button
                  type="button"
                  onClick={() => setFormData({ ...formData, recipientType: 'upi', recipientValue: '' })}
                  className={`p-4 border rounded-lg flex items-center gap-3 transition-colors ${
                    formData.recipientType === 'upi'
                      ? 'border-blue-500 bg-blue-50 text-blue-700'
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                >
                  <CreditCard className="w-5 h-5" />
                  <div className="text-left">
                    <div className="font-medium">UPI ID</div>
                    <div className="text-sm opacity-70">user@paytm</div>
                  </div>
                </button>
                
                <button
                  type="button"
                  onClick={() => setFormData({ ...formData, recipientType: 'phone', recipientValue: '' })}
                  className={`p-4 border rounded-lg flex items-center gap-3 transition-colors ${
                    formData.recipientType === 'phone'
                      ? 'border-blue-500 bg-blue-50 text-blue-700'
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                >
                  <Smartphone className="w-5 h-5" />
                  <div className="text-left">
                    <div className="font-medium">Phone Number</div>
                    <div className="text-sm opacity-70">9876543210</div>
                  </div>
                </button>
              </div>
            </div>

            {/* Recipient Value Input */}
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {formData.recipientType === 'upi' ? 'UPI ID' : 'Phone Number'}
              </label>
              <input
                type="text"
                value={formData.recipientValue}
                onChange={(e) => setFormData({ ...formData, recipientValue: e.target.value })}
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 ${
                  errors.recipientValue ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder={formData.recipientType === 'upi' ? 'Enter UPI ID (e.g., user@paytm)' : 'Enter 10-digit mobile number'}
              />
              {errors.recipientValue && <p className="text-red-500 text-sm mt-1">{errors.recipientValue}</p>}
            </div>

            {/* Recipient Name (Optional) */}
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Recipient Name <span className="text-gray-400">(Optional)</span>
              </label>
              <input
                type="text"
                value={formData.recipientName}
                onChange={(e) => setFormData({ ...formData, recipientName: e.target.value })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter recipient name"
              />
            </div>

            {/* Description (Optional) */}
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Description <span className="text-gray-400">(Optional)</span>
              </label>
              <input
                type="text"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="What's this for?"
              />
            </div>

            {/* Validation Errors */}
            {validation && !validation.valid && (
              <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
                <div className="flex items-center gap-2 text-red-700 mb-2">
                  <AlertCircle className="w-5 h-5" />
                  <span className="font-medium">Validation Failed</span>
                </div>
                <ul className="list-disc list-inside text-red-600 text-sm space-y-1">
                  {validation.errors.map((error: string, index: number) => (
                    <li key={index}>{error}</li>
                  ))}
                </ul>
              </div>
            )}

            {/* General Errors */}
            {errors.general && (
              <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
                <p className="text-red-700">{errors.general}</p>
              </div>
            )}

            {/* Submit Button */}
            <button
              type="submit"
              disabled={validating}
              className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white py-3 px-4 rounded-lg font-medium transition-colors flex items-center justify-center gap-2"
            >
              {validating ? 'Validating...' : 'Continue'}
              {!validating && <ArrowRight className="w-5 h-5" />}
            </button>
          </form>
        </div>
      )}

      {/* Confirmation Step */}
      {step === 'confirm' && validation && (
        <div className="bg-black rounded-lg border border-gray-200 p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Confirm Transfer</h2>
          
          <div className="space-y-4 mb-6">
            <div className="flex justify-between py-3 border-b border-gray-100">
              <span className="text-gray-600">Amount</span>
              <span className="font-semibold">{formatCurrency(formData.amount)}</span>
            </div>
            
            <div className="flex justify-between py-3 border-b border-gray-100">
              <span className="text-gray-600">Transfer Fee</span>
              <span className="font-semibold">{formatCurrency(validation.transfer_fee)}</span>
            </div>
            
            <div className="flex justify-between py-3 border-b border-gray-100">
              <span className="text-gray-600">Total Amount</span>
              <span className="font-semibold text-lg">{formatCurrency(validation.total_amount)}</span>
            </div>
            
            <div className="flex justify-between py-3 border-b border-gray-100">
              <span className="text-gray-600">Recipient</span>
              <span className="font-semibold">{formData.recipientValue}</span>
            </div>
            
            {formData.recipientName && (
              <div className="flex justify-between py-3 border-b border-gray-100">
                <span className="text-gray-600">Recipient Name</span>
                <span className="font-semibold">{formData.recipientName}</span>
              </div>
            )}
            
            <div className="flex justify-between py-3 border-b border-gray-100">
              <span className="text-gray-600">Estimated Time</span>
              <span className="font-semibold">{validation.estimated_time}</span>
            </div>
          </div>

          {validation.warnings && validation.warnings.length > 0 && (
            <div className="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
              <div className="flex items-center gap-2 text-yellow-700 mb-2">
                <AlertCircle className="w-5 h-5" />
                <span className="font-medium">Please Note</span>
              </div>
              <ul className="list-disc list-inside text-yellow-600 text-sm space-y-1">
                {validation.warnings.map((warning: string, index: number) => (
                  <li key={index}>{warning}</li>
                ))}
              </ul>
            </div>
          )}

          <div className="flex gap-4">
            <button
              onClick={() => setStep('form')}
              className="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 py-3 px-4 rounded-lg font-medium transition-colors"
            >
              Back
            </button>
            <button
              onClick={handleCreateTransfer}
              disabled={transferring}
              className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white py-3 px-4 rounded-lg font-medium transition-colors flex items-center justify-center gap-2"
            >
              {transferring ? 'Processing...' : 'Confirm Transfer'}
              {!transferring && <Check className="w-5 h-5" />}
            </button>
          </div>
        </div>
      )}

      {/* Success Step */}
      {step === 'success' && transferResult && (
        <div className="bg-black rounded-lg border border-gray-200 p-6 text-center">
          <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <Check className="w-8 h-8 text-green-600" />
          </div>
          
          <h2 className="text-2xl font-semibold text-gray-900 mb-2">Transfer Initiated!</h2>
          <p className="text-gray-600 mb-6">Your money transfer has been successfully initiated.</p>
          
          <div className="bg-gray-50 rounded-lg p-4 mb-6 text-left">
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-gray-600">Transfer ID</span>
                <span className="font-mono text-sm">{transferResult.transfer_id}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Reference ID</span>
                <span className="font-mono text-sm">{transferResult.reference_id}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Status</span>
                <span className="font-semibold capitalize">{transferResult.status}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Amount</span>
                <span className="font-semibold">{formatCurrency(transferResult.total_amount)}</span>
              </div>
            </div>
          </div>
          
          <div className="flex gap-4">
            <button
              onClick={resetForm}
              className="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-3 px-4 rounded-lg font-medium transition-colors"
            >
              Send Another
            </button>
            <a
              href="/transactions"
              className="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 py-3 px-4 rounded-lg font-medium transition-colors text-center"
            >
              View History
            </a>
          </div>
        </div>
      )}
    </div>
  );
}

function StepIndicator({ steps, currentStep }: { steps: string[]; currentStep: number }) {
  return (
    <div className="flex items-center">
      {steps.map((step, index) => (
        <div key={index} className="flex items-center">
          <div
            className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium ${
              index <= currentStep
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-500'
            }`}
          >
            {index < currentStep ? <Check className="w-4 h-4" /> : index + 1}
          </div>
          {index < steps.length - 1 && (
            <div
              className={`w-12 h-0.5 mx-2 ${
                index < currentStep ? 'bg-blue-600' : 'bg-gray-200'
              }`}
            />
          )}
        </div>
      ))}
    </div>
  );
}
