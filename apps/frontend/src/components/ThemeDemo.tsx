"use client"

import { ThemeSwitcher, ThemeSwitcherWithLabel, ThemeButton } from '@/components/theme-switcher'
import { Card } from '@tranza/ui/components/ui/card-ui'

export default function ThemeDemo() {
  return (
    <div className="max-w-4xl mx-auto p-6 space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">
          Theme Switcher Demo
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          Test the different theme switcher components below
        </p>
      </div>

      {/* Basic Theme Switcher */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
          Basic Theme Switcher
        </h2>
        <div className="flex items-center justify-between">
          <p className="text-gray-600 dark:text-gray-400">
            Simple switch with sun and moon icons
          </p>
          <ThemeSwitcher />
        </div>
      </Card>

      {/* Theme Switcher with Labels */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
          Theme Switcher with Labels
        </h2>
        <div className="flex items-center justify-between">
          <p className="text-gray-600 dark:text-gray-400">
            Switch with "Light" and "Dark" labels
          </p>
          <ThemeSwitcherWithLabel />
        </div>
      </Card>

      {/* Theme Button */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
          Theme Button
        </h2>
        <div className="flex items-center justify-between">
          <p className="text-gray-600 dark:text-gray-400">
            Button style theme toggle
          </p>
          <ThemeButton />
        </div>
      </Card>

      {/* Demo Content */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
          Sample Content
        </h2>
        <div className="space-y-4">
          <p className="text-gray-600 dark:text-gray-400">
            This is sample content to demonstrate how the theme switching affects the appearance.
            Toggle between light and dark themes using any of the switchers above.
          </p>
          
          <div className="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg">
            <h3 className="font-medium text-gray-900 dark:text-gray-100 mb-2">
              Card Example
            </h3>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              This card demonstrates how backgrounds and text colors adapt to the theme.
            </p>
          </div>

          <div className="flex space-x-4">
            <button className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors">
              Primary Button
            </button>
            <button className="bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 px-4 py-2 rounded-lg transition-colors">
              Secondary Button
            </button>
          </div>
        </div>
      </Card>
    </div>
  )
}
