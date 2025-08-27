"use client"

import { Button } from "@/components/ui/Button"
import { 
  FiHeart, 
  FiDownload, 
  FiPlus, 
  FiTrash2, 
  FiEdit, 
  FiShare,
  FiStar,
  FiUser,
  FiSettings
} from "react-icons/fi"

export default function ButtonShowcase() {
  return (
    <div className="max-w-6xl mx-auto p-8 space-y-12">
      <div className="text-center">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Beautiful Button Components
        </h1>
        <p className="text-gray-600 dark:text-gray-400 text-lg">
          Enhanced button components with modern design and smooth animations
        </p>
      </div>

      {/* Variants Section */}
      <div className="space-y-8">
        <h2 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">Button Variants</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div className="bg-black dark:bg-gray-800 rounded-2xl p-6 shadow-lg">
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Default</h3>
            <div className="space-y-3">
              <Button>Default Button</Button>
              <Button leftIcon={<FiPlus />}>Add Item</Button>
              <Button rightIcon={<FiDownload />}>Download</Button>
            </div>
          </div>

          <div className="bg-black dark:bg-gray-800 rounded-2xl p-6 shadow-lg">
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Destructive</h3>
            <div className="space-y-3">
              <Button variant="destructive">Delete</Button>
              <Button variant="destructive" leftIcon={<FiTrash2 />}>Remove</Button>
              <Button variant="destructive" loading>Deleting...</Button>
            </div>
          </div>

          <div className="bg-black dark:bg-gray-800 rounded-2xl p-6 shadow-lg">
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Outline</h3>
            <div className="space-y-3">
              <Button variant="outline">Outline</Button>
              <Button variant="outline" leftIcon={<FiEdit />}>Edit</Button>
              <Button variant="outline" rightIcon={<FiShare />}>Share</Button>
            </div>
          </div>

          <div className="bg-black dark:bg-gray-800 rounded-2xl p-6 shadow-lg">
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Secondary</h3>
            <div className="space-y-3">
              <Button variant="secondary">Secondary</Button>
              <Button variant="secondary" leftIcon={<FiSettings />}>Settings</Button>
              <Button variant="secondary" loading>Loading...</Button>
            </div>
          </div>

          <div className="bg-black dark:bg-gray-800 rounded-2xl p-6 shadow-lg">
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Ghost & Link</h3>
            <div className="space-y-3">
              <Button variant="ghost">Ghost</Button>
              <Button variant="link">Link Button</Button>
              <Button variant="ghost" leftIcon={<FiUser />}>Profile</Button>
            </div>
          </div>

          <div className="bg-black dark:bg-gray-800 rounded-2xl p-6 shadow-lg">
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Success & Warning</h3>
            <div className="space-y-3">
              <Button variant="success">Success</Button>
              <Button variant="warning">Warning</Button>
              <Button variant="premium" leftIcon={<FiStar />}>Premium</Button>
            </div>
          </div>
        </div>
      </div>

      {/* Sizes Section */}
      <div className="space-y-6">
        <h2 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">Button Sizes</h2>
        
        <div className="bg-black dark:bg-gray-800 rounded-2xl p-8 shadow-lg">
          <div className="flex flex-wrap items-center gap-4">
            <Button size="sm">Small</Button>
            <Button size="default">Default</Button>
            <Button size="lg">Large</Button>
            <Button size="xl">Extra Large</Button>
            <Button size="icon" variant="outline">
              <FiHeart />
            </Button>
          </div>
        </div>
      </div>

      {/* Interactive Examples */}
      <div className="space-y-6">
        <h2 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">Interactive Examples</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-black dark:bg-gray-800 rounded-2xl p-6 shadow-lg">
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Loading States</h3>
            <div className="space-y-3">
              <Button loading>Processing...</Button>
              <Button variant="outline" loading>Saving...</Button>
              <Button variant="success" loading>Uploading...</Button>
            </div>
          </div>

          <div className="bg-black dark:bg-gray-800 rounded-2xl p-6 shadow-lg">
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Disabled States</h3>
            <div className="space-y-3">
              <Button disabled>Disabled</Button>
              <Button variant="outline" disabled>Disabled Outline</Button>
              <Button variant="destructive" disabled>Disabled Destructive</Button>
            </div>
          </div>
        </div>
      </div>

      {/* Usage Example */}
      <div className="bg-gradient-to-r from-blue-50 to-purple-50 dark:from-gray-800 dark:to-gray-900 rounded-2xl p-8">
        <h2 className="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-4">Usage Example</h2>
        <div className="bg-gray-900 dark:bg-gray-950 rounded-xl p-4 mb-6">
          <pre className="text-green-400 text-sm overflow-x-auto">
{`import { Button } from "@/components/ui/Button"
import { FiPlus } from "react-icons/fi"

export function MyComponent() {
  return (
    <div className="space-y-4">
      <Button>Default Button</Button>
      <Button variant="outline" size="lg">
        Large Outline
      </Button>
      <Button 
        variant="success" 
        leftIcon={<FiPlus />}
        loading={isLoading}
        onClick={handleClick}
      >
        Add Item
      </Button>
    </div>
  )
}`}
          </pre>
        </div>
        
        <div className="flex flex-wrap gap-4">
          <Button>Try me!</Button>
          <Button variant="outline" size="lg">Large Outline</Button>
          <Button variant="success" leftIcon={<FiPlus />}>Add Item</Button>
          <Button variant="premium" size="xl">✨ Premium Experience</Button>
        </div>
      </div>
    </div>
  )
}
