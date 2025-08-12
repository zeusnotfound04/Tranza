'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { WalletService, PaymentService } from '@/lib/services';
import { TokenManager } from '@/lib/api-client';
import { Button } from '@tranza/ui/components/ui/button';
import { Input } from '@tranza/ui/components/ui/input';
import { Card } from '@tranza/ui/components/ui/card-ui';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { Badge } from '@tranza/ui/components/ui/badge';
import { LoadMoneyRequest, LoadMoneyResponse } from '@/types/api';

declare global {
  interface Window {
    Razorpay: any;
  }
}

interface LoadMoneyModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: (amount: number, newBalance: number) => void;
}

export default function LoadMoneyModal({ isOpen, onClose, onSuccess }: LoadMoneyModalProps) {
  const { user } = useAuth();
  const [amount, setAmount] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [razorpayLoaded, setRazorpayLoaded] = useState(false);

  // Preset amounts for quick selection
  const presetAmounts = [100, 500, 1000, 2000, 5000, 10000];

  useEffect(() => {
    // Load Razorpay script
    const script = document.createElement('script');
    script.src = 'https://checkout.razorpay.com/v1/checkout.js';
    script.onload = () => setRazorpayLoaded(true);
    script.onerror = () => setError('Failed to load payment gateway');
    document.body.appendChild(script);

    return () => {
      document.body.removeChild(script);
    };
  }, []);

  const handleAmountChange = (value: string) => {
    // Remove non-numeric characters except decimal point
    const cleanValue = value.replace(/[^0-9.]/g, '');
    setAmount(cleanValue);
    setError('');
  };

  const handlePresetAmount = (presetAmount: number) => {
    setAmount(presetAmount.toString());
    setError('');
  };

  const validateAmount = (): boolean => {
    const numAmount = parseFloat(amount);
    
    if (!amount || isNaN(numAmount)) {
      setError('Please enter a valid amount');
      return false;
    }
    
    if (numAmount < 10) {
      setError('Minimum amount is ₹10');
      return false;
    }
    
    if (numAmount > 50000) {
      setError('Maximum amount is ₹50,000');
      return false;
    }
    
    return true;
  };

  const handleLoadMoney = async () => {
    if (!validateAmount() || !razorpayLoaded) return;
    
    try {
      setLoading(true);
      setError('');
      
      const numAmount = parseFloat(amount);
      
      // Extract token from cookies and store in localStorage if not already there
      if (!TokenManager.getAccessToken()) {
        const cookies = document.cookie.split(';');
        const accessTokenCookie = cookies.find(cookie => cookie.trim().startsWith('access_token='));
        if (accessTokenCookie) {
          const token = accessTokenCookie.split('=')[1];
          console.log('DEBUG Frontend: Found access_token in cookies, storing in localStorage');
          TokenManager.setTokens(token);
        } else {
          console.log('DEBUG Frontend: No access_token found in cookies!');
          console.log('DEBUG Frontend: All cookies:', document.cookie);
          setError('Please log out and log back in to refresh your session');
          setLoading(false);
          return;
        }
      }
      
      // Debug logs
      console.log('DEBUG Frontend: Current user:', user);
      console.log('DEBUG Frontend: User ID:', user?.id);
      console.log('DEBUG Frontend: Amount being sent:', numAmount);
      console.log('DEBUG Frontend: All cookies:', document.cookie);
      console.log('DEBUG Frontend: Access token from localStorage:', TokenManager.getAccessToken());
      
      // Test authentication by calling getCurrentUser endpoint
      try {
        console.log('DEBUG Frontend: Testing authentication...');
        
        const headers: Record<string, string> = {
          'Content-Type': 'application/json'
        };
        
        // Add Bearer token if available
        const accessToken = TokenManager.getAccessToken();
        if (accessToken) {
          headers.Authorization = `Bearer ${accessToken}`;
          console.log('DEBUG Frontend: Added Authorization header with token');
        }
        
        const authResponse = await fetch('http://localhost:8080/api/v1/auth/me', {
          method: 'GET',
          credentials: 'include',
          headers
        });
        console.log('DEBUG Frontend: Auth test response status:', authResponse.status);
        if (authResponse.ok) {
          const authData = await authResponse.json();
          console.log('DEBUG Frontend: Auth test data:', authData);
        } else {
          console.log('DEBUG Frontend: Auth test failed');
        }
      } catch (authErr) {
        console.log('DEBUG Frontend: Auth test error:', authErr);
      }
      
      // Create Razorpay order
      console.log('DEBUG Frontend: About to call WalletService.createLoadMoneyOrder');
      const response = await WalletService.createLoadMoneyOrder({ amount: numAmount });
      
      console.log('Load money response:', response);
      
      if (!response.data) {
        throw new Error('Failed to create payment order');
      }
      
      const orderData = response.data;
      console.log('Order data:', orderData);
      
      // Initialize Razorpay checkout
      const options = {
        key: orderData.razorpay_key_id, // Use key from backend response
        amount: orderData.amount * 100, // Convert to paise
        currency: orderData.currency,
        order_id: orderData.order_id,
        name: 'Tranza',
        description: `Load ₹${numAmount} to wallet`,
        image: '/logo.png', // Add your logo
        handler: async (response: any) => {
          await handlePaymentSuccess(response, orderData);
        },
        prefill: {
          name: user?.username || '',
          email: user?.email || '',
        },
        theme: {
          color: '#3B82F6'
        },
        modal: {
          ondismiss: () => {
            setLoading(false);
            setError('Payment cancelled');
          }
        }
      };
      
      const razorpay = new window.Razorpay(options);
      razorpay.open();
      
    } catch (err: any) {
      setError(err.message || 'Failed to initiate payment');
      setLoading(false);
    }
  };

  const handlePaymentSuccess = async (paymentResponse: any, orderData: LoadMoneyResponse) => {
    try {
      // Verify payment with backend
      const verificationResponse = await PaymentService.verifyPayment({
        razorpay_payment_id: paymentResponse.razorpay_payment_id,
        razorpay_order_id: paymentResponse.razorpay_order_id,
        razorpay_signature: paymentResponse.razorpay_signature
      });
      
      if (verificationResponse.data?.status === 'success') {
        // Payment successful - the amount was already added by the backend
        const amountAdded = parseFloat(amount);
        
        setLoading(false);
        onSuccess?.(amountAdded, 0); // Backend will handle balance update
        onClose();
        
        // Reset form
        setAmount('');
        setError('');
      } else {
        throw new Error('Payment verification failed');
      }
      
    } catch (err: any) {
      setError(err.message || 'Payment verification failed');
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <Card className="w-full max-w-md max-h-[90vh] rounded-xl overflow-y-auto">
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-black ">Load Money</h2>
            <button
              onClick={onClose}
              className="text-black hover:text-gray-600"
              disabled={loading}
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {!razorpayLoaded && (
            <Alert className="mb-4">
              <AlertDescription>Loading payment gateway...</AlertDescription>
            </Alert>
          )}

          <div className="space-y-4">
            {/* Amount Input */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Enter Amount (₹)
              </label>
              <Input
                type="text"
                placeholder="Enter amount between ₹10 - ₹50,000"
                value={amount}
                onChange={(e) => handleAmountChange(e.target.value)}
                disabled={loading}
                className="text-lg"
              />
              <p className="text-xs text-gray-500 mt-1">
                Min: ₹10 | Max: ₹50,000
              </p>
            </div>

            {/* Preset Amounts */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Quick Select
              </label>
              <div className="grid grid-cols-3 gap-2">
                {presetAmounts.map((preset) => (
                  <Button
                    key={preset}
                    variant="outline"
                    size="sm"
                    onClick={() => handlePresetAmount(preset)}
                    disabled={loading}
                    className={amount === preset.toString() ? 'bg-blue-50 border-blue-200' : ''}
                  >
                    ₹{preset.toLocaleString('en-IN')}
                  </Button>
                ))}
              </div>
            </div>

            {/* Amount Preview */}
            {amount && !isNaN(parseFloat(amount)) && (
              <Card className="p-3 bg-blue-50">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Amount to be added:</span>
                  <span className="font-semibold text-blue-600">
                    ₹{parseFloat(amount).toLocaleString('en-IN', { 
                      minimumFractionDigits: 2, 
                      maximumFractionDigits: 2 
                    })}
                  </span>
                </div>
              </Card>
            )}

            {/* Payment Info */}
            <div className="space-y-2">
              <div className="flex items-center text-xs text-gray-600">
                <svg className="w-4 h-4 text-green-500 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Instant credit to wallet
              </div>
              <div className="flex items-center text-xs text-gray-600">
                <svg className="w-4 h-4 text-green-500 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Secure payment via Razorpay
              </div>
              <div className="flex items-center text-xs text-gray-600">
                <svg className="w-4 h-4 text-green-500 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                All major cards and UPI supported
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex space-x-3 pt-4">
              <Button
                variant="outline"
                onClick={onClose}
                disabled={loading}
                className="flex-1"
              >
                Cancel
              </Button>
              <Button
                onClick={handleLoadMoney}
                disabled={loading || !razorpayLoaded || !amount}
                className="flex-1 text-black"
              >
                {loading ? (
                  <div className="flex items-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    Processing...
                  </div>
                ) : (
                  `Load ₹${amount ? parseFloat(amount).toLocaleString('en-IN') : '0'}`
                )}
              </Button>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}
