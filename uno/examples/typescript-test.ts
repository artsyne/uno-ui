/**
 * TypeScript Test File
 *
 * This file demonstrates and validates TypeScript type definitions.
 * You can use this to verify that types work correctly in your IDE.
 */

import type {
  UnoOptions,
  BoxOptions,
  BoxPosition,
  NotifyOptions,
  NotifyPosition,
  TooltipOptions,
  TooltipPlacement,
  TabOptions
} from 'uno-ui'

// ============================================================
// Test 1: Core Configuration
// ============================================================

const config: UnoOptions = {
  debug: true,
  csp: {
    nonce: 'abc123'
  },
  plugins: {
    box: {
      defaultWidth: 500,
      closeOnEsc: true
    },
    notify: {
      duration: 5000,
      position: 'top-right'
    }
  },
  ui: {
    base: 'rounded-lg shadow-md',
    layers: {
      box: 1400,
      notify: 1600,
      tooltip: 1700
    }
  }
}

// ============================================================
// Test 2: Box Component Types
// ============================================================

const boxPosition: BoxPosition = 'center' // ✅ Valid
// const invalidPos: BoxPosition = 'invalid' // ❌ Type error

const boxOptions: BoxOptions = {
  type: 'confirm',
  title: 'Confirmation',
  content: 'Are you sure?',
  position: 'top-center',
  width: 500,
  showOkButton: true,
  showCancelButton: true,
  closeOnEsc: true,
  onOk: () => {
    console.log('OK clicked')
    return false // Prevent close
  },
  onCancel: () => {
    console.log('Cancel clicked')
  },
  hxUrl: '/api/content',
  hxMethod: 'GET',
  hxHeaders: {
    'X-CSRF-Token': 'token'
  }
}

// ============================================================
// Test 3: Notify Component Types
// ============================================================

const notifyPosition: NotifyPosition = 'bottom-right' // ✅ Valid

const notifyOptions: NotifyOptions = {
  type: 'success',
  position: 'top-right',
  duration: 3000,
  closable: true,
  icon: 'fas fa-check',
  hxUrl: '/api/notification'
}

// ============================================================
// Test 4: Tooltip Component Types
// ============================================================

const tooltipPlacement: TooltipPlacement = 'top-start' // ✅ Valid

const tooltipOptions: TooltipOptions = {
  content: 'This is a tooltip',
  placement: 'bottom-end',
  trigger: 'hover',
  offset: 10,
  arrow: true,
  clickOutside: true
}

// ============================================================
// Test 5: Tab Component Types
// ============================================================

const tabOptions: TabOptions = {
  name: 'profile',
  group: 'settings',
  hash: true,
  updateHash: false
}

// ============================================================
// Test 6: Store Usage (requires Alpine.js)
// ============================================================

// This would work in a real application with Alpine.js loaded
/*
declare const Alpine: any

// Box store
Alpine.store('uno').box.open(boxOptions)
Alpine.store('uno').box.alert('Hello World')
Alpine.store('uno').box.confirm('Are you sure?').then((confirmed: boolean) => {
  console.log('Confirmed:', confirmed)
})

// Notify store
const id: number = Alpine.store('uno').notify.success('Success!')
Alpine.store('uno').notify.error('Error message')
Alpine.store('uno').notify.remove(id)

// Tooltip store
Alpine.store('uno').tooltip.show(tooltipOptions)
Alpine.store('uno').tooltip.hide()

// Tab store
Alpine.store('uno').tab.activate(tabOptions)
const isActive: boolean = Alpine.store('uno').tab.isActive('profile', 'settings')
*/

// ============================================================
// Test 7: Callback Function Types
// ============================================================

const boxWithCallbacks: BoxOptions = {
  title: 'Test',
  onOk: (): void | false => {
    console.log('OK')
    return false // Type-safe return
  },
  onCancel: (): void => {
    console.log('Cancel')
  },
  onContentLoad: (el: HTMLElement): void => {
    console.log('Loaded:', el)
  }
}

// ============================================================
// Test 8: Type Guards
// ============================================================

function isValidBoxPosition(pos: string): pos is BoxPosition {
  const validPositions: BoxPosition[] = [
    'center',
    'top-left',
    'top-center',
    'top-right',
    'center-left',
    'center-right',
    'bottom-left',
    'bottom-center',
    'bottom-right'
  ]
  return validPositions.includes(pos as BoxPosition)
}

// ============================================================
// Test 9: Generic Helper Functions
// ============================================================

function createBoxDialog<T extends BoxOptions>(options: T): T {
  return options
}

const dialog = createBoxDialog({
  title: 'Generic Test',
  position: 'center'
})

// ============================================================
// Test 10: Union Type Testing
// ============================================================

type AllPositions = BoxPosition | NotifyPosition | TooltipPlacement

const testPosition: AllPositions = 'center' // ✅ Valid across all

// ============================================================
// Summary
// ============================================================

console.log('✅ All TypeScript type tests passed!')
console.log('📘 Type definitions are working correctly')
console.log('🎉 IntelliSense and type checking enabled')

export {}
