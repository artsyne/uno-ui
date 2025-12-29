import { measure, debounce } from '../src/utils/performance.js'

console.log('--- Performance Benchmark ---')

// 1. Measure synchronous function
const expensiveFn = () => {
  let sum = 0
  for (let i = 0; i < 1e6; i++) sum += i
  return sum
}

const result = measure('Expensive Loop', expensiveFn)
console.log(`Loop Result: ${result}`)

// 2. Verify Debounce
console.log('\n--- Debounce Test ---')
let callCount = 0
const debouncedFn = debounce(() => {
  callCount++
  console.log('Debounced function executed')
}, 100)

debouncedFn()
debouncedFn()
debouncedFn()

setTimeout(() => {
  console.log(`Call Count (should be 1): ${callCount}`)
  if (callCount === 1) console.log('✅ Debounce Passed')
  else console.error('❌ Debounce Failed')
}, 200)
