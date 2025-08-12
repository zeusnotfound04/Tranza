'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import withAuth from '@/hooks/useAuth';
import { WalletService } from '@/lib/services';
import { Card } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@tranza/ui/components/ui/button';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import LoadMoneyModal from '@/components/wallet/LoadMoneyModal';
import PaymentHistory from '@/components/wallet/PaymentHistory';
import { 
  CreditCard, 
  Smartphone, 
  Building2, 
  Wallet2, 
  CheckCircle, 
  ArrowLeft
} from 'lucide-react';
import { Wallet } from '@/types/api';

function LoadMoneyPage() {
  const router = useRouter();
  const [wallet, setWallet] = useState<Wallet | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showLoadModal, setShowLoadModal] = useState(false);
  const [loadSuccess, setLoadSuccess] = useState<{ amount: number; newBalance: number } | null>(null);

  useEffect(() => {
    loadWallet();
  }, []);

  const loadWallet = async () => {
    try {
      setLoading(true);
      setError('');
      
      const response = await WalletService.getWallet();
      if (response.data) {
        setWallet(response.data);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load wallet');
    } finally {
      setLoading(false);
    }
  };

  const handleLoadSuccess = (amount: number, newBalance: number) => {
    setLoadSuccess({ amount, newBalance });
    setShowLoadModal(false);
    
    // Reload wallet to get updated balance
    loadWallet();
    
    // Clear success message after 5 seconds
    setTimeout(() => {
      setLoadSuccess(null);
    }, 5000);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-4 border-gray-200 border-t-black mx-auto"></div>
          <p className="mt-4 text-lg font-medium text-gray-900">Loading your wallet...</p>
          <p className="text-sm text-gray-600">Please wait</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button
                variant="outline"
                onClick={() => router.back()}
                className="flex items-center space-x-2"
              >
                <ArrowLeft className="w-4 h-4" />
                <span>Back</span>
              </Button>
              <div>
                <h1 className="text-3xl font-bold text-gray-900">
                  Load Money
                </h1>
                <p className="text-gray-600">Add funds to your Tranza wallet</p>
              </div>
            </div>
          </div>
        </div>

        {/* Success Message */}
        {loadSuccess && (
          <div className="mb-6">
            <Alert className="bg-green-50 border-green-200">
              <div className="flex items-center space-x-3">
                <CheckCircle className="w-5 h-5 text-green-600" />
                <AlertDescription className="text-green-800 font-medium">
                  Successfully added ₹{loadSuccess.amount.toLocaleString('en-IN')} to your wallet!
                </AlertDescription>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setLoadSuccess(null)}
                  className="ml-auto h-8 w-8 p-0"
                >
                  ✕
                </Button>
              </div>
            </Alert>
          </div>
        )}

        {/* Error */}
        {error && (
          <div className="mb-6">
            <Alert variant="destructive" className="bg-red-50 border-red-200">
              <AlertDescription className="flex items-center justify-between">
                <span>{error}</span>
                <Button onClick={loadWallet} size="sm" variant="outline" className="ml-4">
                  Retry
                </Button>
              </AlertDescription>
            </Alert>
          </div>
        )}

        {/* Current Balance */}
        {wallet && (
          <Card className="p-6 mb-8 border border-gray-200">
            <div className="text-center">
              <h2 className="text-lg font-semibold text-gray-900 mb-2">Current Balance</h2>
              <p className="text-4xl font-bold text-gray-900 mb-4">
                ₹{wallet.balance.toLocaleString('en-IN', { 
                  minimumFractionDigits: 2, 
                  maximumFractionDigits: 2 
                })}
              </p>
              <div className="flex items-center justify-center space-x-4">
                <div className="flex items-center space-x-2">
                  <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                  <span className="text-sm text-gray-600">Available</span>
                </div>
              </div>
            </div>
          </Card>
        )}

        {/* Load Money Button */}
        <Card className="p-6 mb-8 border border-gray-200">
          <div className="text-center">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Add Money to Wallet</h3>
            <Button
              onClick={() => setShowLoadModal(true)}
              size="lg"
              className="px-8 py-3 bg-black text-white hover:bg-gray-800"
            >
              Load Money
            </Button>
            <p className="text-sm text-gray-500 mt-3">
              Secure payments • ₹10 - ₹50,000 range
            </p>
          </div>
        </Card>

        {/* Payment Methods */}
        <Card className="p-6 mb-8 border border-gray-200">
          <h3 className="text-lg font-semibold text-gray-900 mb-6">Supported Payment Methods</h3>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex items-center space-x-3 p-4 bg-gray-50 rounded-lg">
              <CreditCard className="w-6 h-6 text-gray-600" />
              <div>
                <p className="font-medium text-gray-900">Credit/Debit Cards</p>
                <p className="text-sm text-gray-600">Visa, Mastercard, RuPay</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3 p-4 bg-gray-50 rounded-lg">
              <Smartphone className="w-6 h-6 text-gray-600" />
              <div>
                <p className="font-medium text-gray-900">UPI</p>
                <p className="text-sm text-gray-600">All UPI apps supported</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3 p-4 bg-gray-50 rounded-lg">
              <Building2 className="w-6 h-6 text-gray-600" />
              <div>
                <p className="font-medium text-gray-900">Net Banking</p>
                <p className="text-sm text-gray-600">All major banks</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3 p-4 bg-gray-50 rounded-lg">
              <Wallet2 className="w-6 h-6 text-gray-600" />
              <div>
                <p className="font-medium text-gray-900">Digital Wallets</p>
                <p className="text-sm text-gray-600">Paytm, PhonePe, etc.</p>
              </div>
            </div>
          </div>
          
          <div className="mt-6 p-4 bg-gray-50 rounded-lg">
            <div className="flex items-center space-x-3">
              <CheckCircle className="w-5 h-5 text-green-600" />
              <div>
                <p className="font-medium text-gray-900">Instant Credit</p>
                <p className="text-sm text-gray-600">Money is added to your wallet immediately after successful payment</p>
              </div>
            </div>
          </div>
        </Card>

        {/* Recent Load Transactions */}
        <Card className="p-6 border border-gray-200">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-lg font-semibold text-gray-900">Recent Transactions</h3>
            <Button
              variant="outline"
              onClick={() => router.push('/transactions')}
              className="text-sm"
            >
              View All
            </Button>
          </div>
          <PaymentHistory 
            limit={5} 
            showFilters={false}
          />
        </Card>

        {/* Load Money Modal */}
        <LoadMoneyModal
          isOpen={showLoadModal}
          onClose={() => setShowLoadModal(false)}
          onSuccess={handleLoadSuccess}
        />
      </div>
    </div>
  );
}

export default withAuth(LoadMoneyPage);
