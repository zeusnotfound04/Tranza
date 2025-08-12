'use client';

import { useEffect } from 'react';

interface RazorpayCheckoutProps {
  orderId: string;
  amount: number;
  currency: string;
  onSuccess: (response: any) => void;
  onFailure: (error: any) => void;
  onDismiss?: () => void;
  userDetails?: {
    name?: string;
    email?: string;
    contact?: string;
  };
}

// Define Razorpay types
declare global {
  interface Window {
    Razorpay: any;
  }
}

export default function RazorpayCheckout({
  orderId,
  amount,
  currency,
  onSuccess,
  onFailure,
  onDismiss,
  userDetails = {}
}: RazorpayCheckoutProps) {

  useEffect(() => {
    const loadRazorpayScript = async () => {
      // Load Razorpay script if not already loaded
      if (!window.Razorpay) {
        const script = document.createElement('script');
        script.src = 'https://checkout.razorpay.com/v1/checkout.js';
        script.async = true;
        document.body.appendChild(script);
        
        await new Promise((resolve, reject) => {
          script.onload = resolve;
          script.onerror = reject;
        });
      }

      // Configure Razorpay checkout options
      const options = {
        key: process.env.NEXT_PUBLIC_RAZORPAY_KEY_ID,
        amount: amount * 100, // Amount in paise
        currency: currency,
        name: 'Tranza',
        description: 'Wallet Load Payment',
        order_id: orderId,
        handler: function (response: any) {
          console.log('Razorpay payment successful:', response);
          onSuccess(response);
        },
        prefill: {
          name: userDetails.name || '',
          email: userDetails.email || '',
          contact: userDetails.contact || ''
        },
        notes: {
          address: 'Tranza Payment'
        },
        theme: {
          color: '#3B82F6'
        },
        modal: {
          ondismiss: function() {
            console.log('Razorpay modal dismissed');
            if (onDismiss) {
              onDismiss();
            }
          }
        },
        retry: {
          enabled: true,
          max_count: 3
        },
        remember_customer: false,
        readonly: {
          email: false,
          contact: false,
          name: false
        },
        hidden: {
          email: false,
          contact: false,
          name: false
        }
      };

      try {
        const rzp = new window.Razorpay(options);
        
        // Handle payment failure
        rzp.on('payment.failed', function (response: any) {
          console.error('Razorpay payment failed:', response);
          onFailure(response.error);
        });

        // Open the checkout
        rzp.open();
        
      } catch (error) {
        console.error('Error initializing Razorpay:', error);
        onFailure(error);
      }
    };

    loadRazorpayScript();
  }, [orderId, amount, currency, onSuccess, onFailure, onDismiss, userDetails]);

  // This component doesn't render any UI - it just handles the Razorpay integration
  return null;
}

// Utility function to load Razorpay script
export const loadRazorpayScript = (): Promise<boolean> => {
  return new Promise((resolve) => {
    // Check if Razorpay is already loaded
    if (window.Razorpay) {
      resolve(true);
      return;
    }

    const script = document.createElement('script');
    script.src = 'https://checkout.razorpay.com/v1/checkout.js';
    script.async = true;
    
    script.onload = () => {
      resolve(true);
    };
    
    script.onerror = () => {
      resolve(false);
    };
    
    document.body.appendChild(script);
  });
};

// Utility function for manual Razorpay checkout
export const openRazorpayCheckout = async (options: {
  orderId: string;
  amount: number;
  currency: string;
  onSuccess: (response: any) => void;
  onFailure: (error: any) => void;
  onDismiss?: () => void;
  userDetails?: {
    name?: string;
    email?: string;
    contact?: string;
  };
}) => {
  const isLoaded = await loadRazorpayScript();
  
  if (!isLoaded) {
    options.onFailure(new Error('Failed to load Razorpay script'));
    return;
  }

  const rzpOptions = {
    key: process.env.NEXT_PUBLIC_RAZORPAY_KEY_ID,
    amount: options.amount * 100,
    currency: options.currency,
    name: 'Tranza',
    description: 'Wallet Payment',
    order_id: options.orderId,
    handler: options.onSuccess,
    prefill: {
      name: options.userDetails?.name || '',
      email: options.userDetails?.email || '',
      contact: options.userDetails?.contact || ''
    },
    theme: {
      color: '#3B82F6'
    },
    modal: {
      ondismiss: options.onDismiss
    }
  };

  const rzp = new window.Razorpay(rzpOptions);
  
  rzp.on('payment.failed', function (response: any) {
    options.onFailure(response.error);
  });

  rzp.open();
};
