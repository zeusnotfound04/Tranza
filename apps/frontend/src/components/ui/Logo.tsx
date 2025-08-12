'use client';

import Link from 'next/link';
import Image from 'next/image';

interface LogoProps {
  href?: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export default function Logo({ 
  href = '/', 
  size = 'md', 
  className = '' 
}: LogoProps) {
  const sizeClasses = {
    sm: {
      width: 24,
      height: 24
    },
    md: {
      width: 32,
      height: 32
    },
    lg: {
      width: 48,
      height: 48
    }
  };

  const LogoContent = () => (
    <div className={`flex items-center ${className}`}>
      <Image
        src="/logo.png"
        alt="Tranza Logo"
        width={sizeClasses[size].width}
        height={sizeClasses[size].height}
        className="object-contain"
        priority
      />
    </div>
  );

  if (href) {
    return (
      <Link href={href} className="inline-flex">
        <LogoContent />
      </Link>
    );
  }

  return <LogoContent />;
}
