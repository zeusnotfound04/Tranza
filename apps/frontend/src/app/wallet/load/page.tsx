'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function LoadRedirect() {
  const router = useRouter();

  useEffect(() => {
    // Redirect to the correct load money page
    router.replace('/wallet/load-money');
  }, [router]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <p className="mt-2 text-gray-600">Redirecting to load money...</p>
      </div>
    </div>
  );
}
