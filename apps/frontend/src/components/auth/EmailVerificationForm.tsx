'use client';

import React, { useState, useEffect } from 'react';
import { authService, AuthError } from '@/lib/auth';
import { EmailVerificationRequest, ResendVerificationRequest } from '@/types/auth';

interface EmailVerificationFormProps {
  email: string;
  expiresAt: string;
  onSuccess: () => void;
  onBack: () => void;
}

export default function EmailVerificationForm({ 
  email, 
  expiresAt, 
  onSuccess, 
  onBack 
}: EmailVerificationFormProps) {
  const [code, setCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isResending, setIsResending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [timeLeft, setTimeLeft] = useState<number>(0);

  // Calculate time left
  useEffect(() => {
    const updateTimeLeft = () => {
      const now = new Date().getTime();
      const expiry = new Date(expiresAt).getTime();
      const difference = expiry - now;
      
      if (difference > 0) {
        setTimeLeft(Math.floor(difference / 1000));
      } else {
        setTimeLeft(0);
      }
    };

    updateTimeLeft();
    const interval = setInterval(updateTimeLeft, 1000);

    return () => clearInterval(interval);
  }, [expiresAt]);

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setIsLoading(true);

    try {
      const verificationData: EmailVerificationRequest = {
        email,
        code: code.trim(),
      };

      const response = await authService.verifyEmail(verificationData);
      setSuccess(response.message);
      
      // Wait a moment before calling onSuccess to show the success message
      setTimeout(() => {
        onSuccess();
      }, 1500);
    } catch (error) {
      const message = error instanceof AuthError ? error.message : 'Verification failed';
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleResend = async () => {
    setError(null);
    setSuccess(null);
    setIsResending(true);

    try {
      const resendData: ResendVerificationRequest = { email };
      const response = await authService.resendVerificationCode(resendData);
      setSuccess('New verification code sent to your email');
    } catch (error) {
      const message = error instanceof AuthError ? error.message : 'Failed to resend code';
      setError(message);
    } finally {
      setIsResending(false);
    }
  };

  const handleCodeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 6);
    setCode(value);
  };

  return (
    <div className="w-full max-w-md mx-auto">
      <div className="bg-white shadow-md rounded-lg px-8 pt-6 pb-8">
        <div className="text-center mb-6">
          <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">📧</span>
          </div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">
            Verify Your Email
          </h2>
          <p className="text-gray-600 text-sm">
            We've sent a 6-digit code to
          </p>
          <p className="font-semibold text-gray-800">{email}</p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {success && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-4">
            {success}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-6">
            <label 
              className="block text-gray-700 text-sm font-bold mb-2 text-center"
              htmlFor="code"
            >
              Enter Verification Code
            </label>
            <input
              className="text-center text-2xl font-mono letter-spacing-wide shadow appearance-none border rounded w-full py-3 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500"
              id="code"
              type="text"
              placeholder="000000"
              value={code}
              onChange={handleCodeChange}
              maxLength={6}
              required
              disabled={isLoading}
              style={{ letterSpacing: '0.5em' }}
            />
            <div className="text-center mt-2">
              {timeLeft > 0 ? (
                <p className="text-sm text-gray-500">
                  Code expires in {formatTime(timeLeft)}
                </p>
              ) : (
                <p className="text-sm text-red-500">
                  Code has expired
                </p>
              )}
            </div>
          </div>

          <div className="flex flex-col space-y-3">
            <button
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50 disabled:cursor-not-allowed"
              type="submit"
              disabled={isLoading || code.length !== 6}
            >
              {isLoading ? (
                <div className="flex items-center justify-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Verifying...
                </div>
              ) : (
                'Verify Email'
              )}
            </button>

            <button
              type="button"
              onClick={handleResend}
              disabled={isResending || timeLeft > 0}
              className="text-blue-500 hover:text-blue-700 font-medium py-2 disabled:text-gray-400 disabled:cursor-not-allowed"
            >
              {isResending ? (
                <div className="flex items-center justify-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-500 mr-2"></div>
                  Resending...
                </div>
              ) : timeLeft > 0 ? (
                `Resend code in ${formatTime(timeLeft)}`
              ) : (
                'Resend verification code'
              )}
            </button>

            <button
              type="button"
              onClick={onBack}
              className="text-gray-500 hover:text-gray-700 font-medium py-2"
              disabled={isLoading}
            >
              ← Back to registration
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}