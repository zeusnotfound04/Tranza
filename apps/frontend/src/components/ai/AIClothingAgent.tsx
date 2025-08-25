'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Send, ShoppingBag, MapPin, CreditCard, Loader2, CheckCircle, XCircle } from 'lucide-react';
import { aeonikPro } from '@/lib/fonts';
import { 
  AIClothingService, 
  AddressService, 
  AIClothingOrderRequest, 
  AIClothingOrderResponse,
  Address,
  SuggestedProduct,
  SelectedProductForOrder,
  ConfirmAIClothingOrderRequest
} from '@/services/AIClothingService';

interface Message {
  id: string;
  type: 'user' | 'ai';
  content: string;
  timestamp: Date;
  suggestions?: SuggestedProduct[];
  selectedProducts?: SelectedProductForOrder[];
  requiresConfirmation?: boolean;
}

interface AddressModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAddressCreated: (address: Address) => void;
}

const AddressModal: React.FC<AddressModalProps> = ({ isOpen, onClose, onAddressCreated }) => {
  const [formData, setFormData] = useState({
    name: '',
    phone: '',
    address_line: '',
    city: '',
    state: '',
    pin_code: '',
    landmark: '',
    address_type: 'home'
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const address = await AddressService.createAddress(formData);
      onAddressCreated(address);
      onClose();
      setFormData({
        name: '',
        phone: '',
        address_line: '',
        city: '',
        state: '',
        pin_code: '',
        landmark: '',
        address_type: 'home'
      });
    } catch (error) {
      console.error('Error creating address:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
        <h3 className="text-lg font-semibold mb-4">Add Delivery Address</h3>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Full Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({...formData, name: e.target.value})}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-black focus:border-transparent"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Phone Number</label>
            <input
              type="tel"
              value={formData.phone}
              onChange={(e) => setFormData({...formData, phone: e.target.value})}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-black focus:border-transparent"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Address Line</label>
            <textarea
              value={formData.address_line}
              onChange={(e) => setFormData({...formData, address_line: e.target.value})}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-black focus:border-transparent"
              rows={2}
              required
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">City</label>
              <input
                type="text"
                value={formData.city}
                onChange={(e) => setFormData({...formData, city: e.target.value})}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-black focus:border-transparent"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">State</label>
              <input
                type="text"
                value={formData.state}
                onChange={(e) => setFormData({...formData, state: e.target.value})}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-black focus:border-transparent"
                required
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">PIN Code</label>
            <input
              type="text"
              value={formData.pin_code}
              onChange={(e) => setFormData({...formData, pin_code: e.target.value})}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-black focus:border-transparent"
              required
              maxLength={6}
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Landmark (Optional)</label>
            <input
              type="text"
              value={formData.landmark}
              onChange={(e) => setFormData({...formData, landmark: e.target.value})}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-black focus:border-transparent"
            />
          </div>
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 py-2 px-4 border border-gray-300 rounded-lg hover:bg-gray-50"
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 py-2 px-4 bg-black text-white rounded-lg hover:bg-gray-800 disabled:opacity-50"
              disabled={isSubmitting}
            >
              {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin mx-auto" /> : 'Save Address'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const ProductSuggestionCard: React.FC<{ suggestion: SuggestedProduct }> = ({ suggestion }) => {
  return (
    <div className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
      <div className="flex gap-4">
        <img
          src={suggestion.image_url || '/placeholder-product.jpg'}
          alt={suggestion.name}
          className="w-16 h-16 object-cover rounded-lg"
        />
        <div className="flex-1">
          <h4 className="font-semibold text-sm">{suggestion.name}</h4>
          <p className="text-xs text-gray-600">{suggestion.brand}</p>
          <p className="text-xs text-gray-500 mt-1">{suggestion.description}</p>
          <div className="flex items-center justify-between mt-2">
            <span className="font-semibold">{suggestion.currency} {suggestion.price}</span>
            <div className="flex items-center text-xs text-yellow-600">
              <span>★ {suggestion.rating.toFixed(1)}</span>
            </div>
          </div>
          <div className="flex justify-between items-center mt-2">
            <span className="text-xs text-blue-600">{suggestion.website}</span>
            <a 
              href={suggestion.url} 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-xs bg-blue-500 text-white px-2 py-1 rounded hover:bg-blue-600"
            >
              View Product
            </a>
          </div>
        </div>
      </div>
    </div>
  );
};

interface ProductSelectionInterfaceProps {
  products: SuggestedProduct[];
  onConfirm: (selectedProducts: SelectedProductForOrder[]) => void;
}

const ProductSelectionInterface: React.FC<ProductSelectionInterfaceProps> = ({ products, onConfirm }) => {
  const [selectedProducts, setSelectedProducts] = useState<Map<string, SelectedProductForOrder>>(new Map());

  const handleProductToggle = (product: SuggestedProduct) => {
    const newSelected = new Map(selectedProducts);
    
    if (newSelected.has(product.id)) {
      newSelected.delete(product.id);
    } else {
      newSelected.set(product.id, {
        product_id: product.id,
        name: product.name,
        price: product.price,
        quantity: 1,
        url: product.url,
        website: product.website,
      });
    }
    
    setSelectedProducts(newSelected);
  };

  const handleQuantityChange = (productId: string, quantity: number) => {
    const newSelected = new Map(selectedProducts);
    const product = newSelected.get(productId);
    
    if (product && quantity > 0) {
      product.quantity = quantity;
      newSelected.set(productId, product);
      setSelectedProducts(newSelected);
    }
  };

  const getTotalAmount = () => {
    return Array.from(selectedProducts.values()).reduce(
      (total, product) => total + (product.price * product.quantity), 
      0
    );
  };

  const handleConfirm = () => {
    if (selectedProducts.size > 0) {
      onConfirm(Array.from(selectedProducts.values()));
    }
  };

  return (
    <div className="space-y-3">
      {products.map((product) => {
        const isSelected = selectedProducts.has(product.id);
        const selectedProduct = selectedProducts.get(product.id);
        
        return (
          <div
            key={product.id}
            className={`border rounded-lg p-3 cursor-pointer transition-colors ${
              isSelected ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-gray-300'
            }`}
            onClick={() => handleProductToggle(product)}
          >
            <div className="flex items-center space-x-3">
              <input
                type="checkbox"
                checked={isSelected}
                onChange={() => handleProductToggle(product)}
                className="w-4 h-4 text-blue-600"
              />
              <img
                src={product.image_url || '/placeholder-product.jpg'}
                alt={product.name}
                className="w-12 h-12 object-cover rounded"
              />
              <div className="flex-1">
                <h4 className="font-medium text-sm">{product.name}</h4>
                <p className="text-xs text-gray-600">{product.brand}</p>
                <p className="text-sm font-semibold">{product.currency} {product.price}</p>
              </div>
              
              {isSelected && (
                <div className="flex items-center space-x-2" onClick={(e) => e.stopPropagation()}>
                  <label className="text-xs">Qty:</label>
                  <input
                    type="number"
                    min="1"
                    max="10"
                    value={selectedProduct?.quantity || 1}
                    onChange={(e) => handleQuantityChange(product.id, parseInt(e.target.value))}
                    className="w-16 px-2 py-1 text-xs border rounded"
                  />
                </div>
              )}
            </div>
          </div>
        );
      })}
      
      {selectedProducts.size > 0 && (
        <div className="border-t pt-3">
          <div className="flex justify-between items-center mb-3">
            <span className="font-medium">
              Total: ₹{getTotalAmount().toFixed(2)} ({selectedProducts.size} items)
            </span>
          </div>
          <div className="flex gap-2">
            <button
              onClick={handleConfirm}
              className="flex-1 bg-green-600 text-white py-2 px-4 rounded hover:bg-green-700 text-sm"
            >
              Confirm Purchase
            </button>
            <button
              onClick={() => setSelectedProducts(new Map())}
              className="bg-gray-500 text-white py-2 px-4 rounded hover:bg-gray-600 text-sm"
            >
              Clear All
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

const AIClothingAgent: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: '1',
      type: 'ai',
      content: '👋 Hi! I\'m your AI shopping assistant. I can help you find and order clothes using your Tranza wallet. Just tell me what you\'re looking for!\n\nFor example:\n• "I need a formal shirt for office, size M, budget ₹2000"\n• "Buy me casual jeans, size 32, blue color"\n• "I want a dress for a party, budget ₹3000"',
      timestamp: new Date(),
    },
  ]);
  const [inputMessage, setInputMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [defaultAddress, setDefaultAddress] = useState<Address | null>(null);
  const [showAddressModal, setShowAddressModal] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    loadDefaultAddress();
  }, []);

  const loadDefaultAddress = async () => {
    try {
      const address = await AddressService.getDefaultAddress();
      setDefaultAddress(address);
    } catch (error) {
      console.log('No default address found');
    }
  };

  const handleSendMessage = async () => {
    if (!inputMessage.trim() || isLoading) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      type: 'user',
      content: inputMessage,
      timestamp: new Date(),
    };

    setMessages(prev => [...prev, userMessage]);
    setInputMessage('');
    setIsLoading(true);

    try {
      const request: AIClothingOrderRequest = {
        prompt: inputMessage,
        address_id: defaultAddress?.id,
      };

      const response = await AIClothingService.processAIClothingOrder(request);

      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: 'ai',
        content: response.message,
        timestamp: new Date(),
        suggestions: response.suggested_products,
        requiresConfirmation: response.requires_confirmation,
      };

      setMessages(prev => [...prev, aiMessage]);

      // Handle missing address
      if (response.required_info?.includes('delivery_address')) {
        setShowAddressModal(true);
      }
    } catch (error) {
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: 'ai',
        content: '❌ Sorry, I encountered an error processing your request. Please try again.',
        timestamp: new Date(),
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleConfirmOrder = async (selectedProducts: SelectedProductForOrder[], addressId: string) => {
    setIsLoading(true);
    
    try {
      const response = await AIClothingService.confirmAIClothingOrder({
        selected_products: selectedProducts,
        address_id: addressId,
      });

      const message: Message = {
        id: Date.now().toString(),
        type: 'ai',
        content: response.success 
          ? `✅ Great! Your order has been confirmed and payment of ${response.total_amount || 0} has been deducted from your wallet. Order ID: ${response.order_id || 'N/A'}. You will be redirected to the external websites to complete your purchases.`
          : `❌ ${response.message}${response.required_amount ? ` You need ₹${response.required_amount} but only have ₹${response.current_balance}.` : ''}`,
        timestamp: new Date(),
        selectedProducts: response.success ? response.selected_products : undefined,
      };

      setMessages(prev => [...prev, message]);
    } catch (error) {
      const errorMessage: Message = {
        id: Date.now().toString(),
        type: 'ai',
        content: '❌ Error processing order confirmation. Please try again.',
        timestamp: new Date(),
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={`flex flex-col h-full bg-white ${aeonikPro.className}`}>
      {/* Header */}
      <div className="flex items-center gap-3 p-4 border-b border-gray-200 bg-white">
        <div className="w-10 h-10 bg-black rounded-full flex items-center justify-center">
          <ShoppingBag className="w-5 h-5 text-white" />
        </div>
        <div>
          <h2 className="font-semibold text-lg">AI Shopping Assistant</h2>
          <p className="text-sm text-gray-600">Find and order clothes with AI</p>
        </div>
        <div className="ml-auto flex items-center gap-2">
          {defaultAddress ? (
            <div className="flex items-center gap-1 text-xs text-green-600 bg-green-50 px-2 py-1 rounded">
              <MapPin className="w-3 h-3" />
              <span>Address set</span>
            </div>
          ) : (
            <button
              onClick={() => setShowAddressModal(true)}
              className="flex items-center gap-1 text-xs text-orange-600 bg-orange-50 px-2 py-1 rounded hover:bg-orange-100"
            >
              <MapPin className="w-3 h-3" />
              <span>Add address</span>
            </button>
          )}
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message) => (
          <div
            key={message.id}
            className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}
          >
            <div
              className={`max-w-[80%] rounded-lg p-3 ${
                message.type === 'user'
                  ? 'bg-black text-white'
                  : 'bg-gray-100 text-black'
              }`}
            >
              <p className="whitespace-pre-line text-sm">{message.content}</p>
              
              {/* Product suggestions */}
              {message.suggestions && message.suggestions.length > 0 && (
                <div className="mt-3 space-y-2">
                  <p className="text-xs font-medium">Suggested items:</p>
                  {message.suggestions.map((suggestion, index) => (
                    <ProductSuggestionCard key={index} suggestion={suggestion} />
                  ))}
                </div>
              )}

              {/* Product selection for confirmation */}
              {message.requiresConfirmation && message.suggestions && (
                <div className="mt-3">
                  <p className="text-sm font-medium text-gray-700 mb-2">
                    Select products you want to purchase:
                  </p>
                  <ProductSelectionInterface 
                    products={message.suggestions}
                    onConfirm={(selectedProducts: SelectedProductForOrder[]) => {
                      if (defaultAddress) {
                        handleConfirmOrder(selectedProducts, defaultAddress.id);
                      } else {
                        setShowAddressModal(true);
                      }
                    }}
                  />
                </div>
              )}

              {/* Show selected products after confirmation */}
              {message.selectedProducts && (
                <div className="mt-3">
                  <p className="text-sm font-medium text-gray-700 mb-2">
                    Purchased Items:
                  </p>
                  <div className="space-y-2">
                    {message.selectedProducts.map((product, index) => (
                      <div key={index} className="bg-green-50 border border-green-200 rounded-lg p-3">
                        <div className="flex justify-between items-center">
                          <div>
                            <p className="font-medium text-sm">{product.name}</p>
                            <p className="text-xs text-gray-600">
                              Qty: {product.quantity} × ₹{product.price}
                            </p>
                          </div>
                          <a 
                            href={product.url} 
                            target="_blank" 
                            rel="noopener noreferrer"
                            className="text-xs bg-blue-500 text-white px-2 py-1 rounded hover:bg-blue-600"
                          >
                            Complete Purchase on {product.website}
                          </a>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              <p className="text-xs opacity-70 mt-2">
                {message.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </p>
            </div>
          </div>
        ))}
        
        {isLoading && (
          <div className="flex justify-start">
            <div className="bg-gray-100 rounded-lg p-3">
              <div className="flex items-center gap-2">
                <Loader2 className="w-4 h-4 animate-spin" />
                <span className="text-sm">AI is thinking...</span>
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <div className="border-t border-gray-200 p-4">
        <div className="flex gap-2">
          <input
            type="text"
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
            placeholder="Tell me what clothes you want to buy..."
            className="flex-1 border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-black focus:border-transparent"
            disabled={isLoading}
          />
          <button
            onClick={handleSendMessage}
            disabled={isLoading || !inputMessage.trim()}
            className="bg-black text-white p-2 rounded-lg hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Send className="w-5 h-5" />
          </button>
        </div>
      </div>

      {/* Address Modal */}
      <AddressModal
        isOpen={showAddressModal}
        onClose={() => setShowAddressModal(false)}
        onAddressCreated={(address) => {
          setDefaultAddress(address);
          AddressService.setDefaultAddress(address.id);
        }}
      />
    </div>
  );
};

export default AIClothingAgent;
