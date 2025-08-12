'use client';

import React, { useState } from 'react';
import { authService, AuthError } from '@/lib/auth';
import { PreRegistrationRequest, PreRegistrationResponse } from '@/types/auth';
import EmailVerificationForm from './EmailVerificationForm';

interface RegisterFormProps {
  onSuccess?: () => void;
  onSwitchToLogin?: () => void;
}

export default function RegisterForm({ onSuccess, onSwitchToLogin }: RegisterFormProps) {
  const [step, setStep] = useState<'register' | 'verify'>('register');
  const [formData, setFormData] = useState<PreRegistrationRequest>({
    email: '',
    username: '',
    password: '',
  });
  const [confirmPassword, setConfirmPassword] = useState('');
  const [verificationData, setVerificationData] = useState<PreRegistrationResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (formData.password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setIsLoading(true);

    try {
      const response = await authService.preRegister(formData);
      setVerificationData(response);
      setStep('verify');
    } catch (error) {
      const message = error instanceof AuthError ? error.message : 'Registration failed';
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleVerificationSuccess = () => {
    onSuccess?.();
  };

  const handleBackToRegister = () => {
    setStep('register');
    setVerificationData(null);
    setError(null);
  };

  // Show email verification form if we're in verify step
  if (step === 'verify' && verificationData) {
    return (
      <EmailVerificationForm
        email={verificationData.email}
        expiresAt={verificationData.expires_at}
        onSuccess={handleVerificationSuccess}
        onBack={handleBackToRegister}
      />
    );
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
    // Clear error when user starts typing
    if (error) setError(null);
  };

  const handleConfirmPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setConfirmPassword(e.target.value);
    if (error) setError(null);
  };

  const passwordsMatch = formData.password === confirmPassword;
  const showPasswordError = confirmPassword && !passwordsMatch;

  return (
    <div className="w-full max-w-md mx-auto">
      <div className="bg-white shadow-md rounded-lg px-8 pt-6 pb-8 mb-4">
        <h2 className="text-2xl font-bold text-center text-gray-800 mb-6">
          Create Your Account
        </h2>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label 
              className="block text-gray-700 text-sm font-bold mb-2"
              htmlFor="email"
            >
              Email
            </label>
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500"
              id="email"
              type="email"
              name="email"
              placeholder="Enter your email"
              value={formData.email}
              onChange={handleChange}
              required
              disabled={isLoading}
            />
          </div>

          <div className="mb-4">
            <label 
              className="block text-gray-700 text-sm font-bold mb-2"
              htmlFor="username"
            >
              Username
            </label>
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500"
              id="username"
              type="text"
              name="username"
              placeholder="Choose a username"
              value={formData.username}
              onChange={handleChange}
              required
              disabled={isLoading}
            />
          </div>

          <div className="mb-4">
            <label 
              className="block text-gray-700 text-sm font-bold mb-2"
              htmlFor="password"
            >
              Password
            </label>
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500"
              id="password"
              type="password"
              name="password"
              placeholder="Enter your password (min 8 characters)"
              value={formData.password}
              onChange={handleChange}
              required
              minLength={8}
              disabled={isLoading}
            />
          </div>

          <div className="mb-6">
            <label 
              className="block text-gray-700 text-sm font-bold mb-2"
              htmlFor="confirmPassword"
            >
              Confirm Password
            </label>
            <input
              className={`shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500 ${
                showPasswordError ? 'border-red-500' : ''
              }`}
              id="confirmPassword"
              type="password"
              name="confirmPassword"
              placeholder="Confirm your password"
              value={confirmPassword}
              onChange={handleConfirmPasswordChange}
              required
              disabled={isLoading}
            />
            {showPasswordError && (
              <p className="text-red-500 text-xs italic mt-1">
                Passwords do not match
              </p>
            )}
          </div>

          <div className="flex items-center justify-between">
            <button
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50 disabled:cursor-not-allowed w-full"
              type="submit"
              disabled={isLoading || !passwordsMatch || (error === 'Passwords do not match')}
            >
              {isLoading ? (
                <div className="flex items-center justify-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Sending Verification...
                </div>
              ) : (
                'Send Verification Code'
              )}
            </button>
          </div>
        </form>

        {onSwitchToLogin && (
          <div className="text-center mt-4">
            <p className="text-gray-600 text-sm">
              Already have an account?{' '}
              <button
                type="button"
                onClick={onSwitchToLogin}
                className="text-blue-500 hover:text-blue-700 font-medium"
                disabled={isLoading}
              >
                Sign in
              </button>
            </p>
          </div>
        )}
      </div>
    </div>
  );
}