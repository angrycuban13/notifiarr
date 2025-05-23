import { toast } from '@zerodevx/svelte-toast'

/** Remove a prefix from a string. */
export const ltrim = (str: string, prefix: string) =>
  str.slice(str.startsWith(prefix) ? prefix.length : 0)

/** Remove a suffix from a string. */
export const rtrim = (str: string, suffix: string) =>
  str.endsWith(suffix) ? str.slice(0, str.length - suffix.length) : str

/** Show a success toast. */
export const success = (m: string) =>
  toast.push(m, {
    theme: {
      '--toastBackground': 'green',
      '--toastColor': 'white',
      '--toastBarBackground': 'olive',
    },
  })

/** Show a warning toast. */
export const warning = (m: string) =>
  toast.push(m, {
    theme: {
      '--toastBackground': 'orange',
      '--toastColor': 'white',
      '--toastBarBackground': 'black',
    },
  })

/** Show a failure toast. */
export const failure = (m: string) =>
  toast.push(m, {
    theme: {
      '--toastBackground': 'red',
      '--toastColor': 'white',
      '--toastBarBackground': 'royalblue',
    },
  })

/** age converts a milliseconds counter into human readable: 13h 5m 45s */
export function age(milliseconds: number, includeSeconds = false): string {
  const seconds = Math.floor(milliseconds / 1000)
  if (!seconds) return '0s'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds - days * 86400) / 3600)
  const minutes = Math.floor((seconds - days * 86400 - hours * 3600) / 60)
  const secs =
    !includeSeconds && seconds > 60
      ? 0
      : Math.floor(seconds - days * 86400 - hours * 3600 - minutes * 60)
  return (
    (days > 0 ? days + 'd ' : '') +
    (hours > 0 ? hours + 'h ' : '') +
    (minutes > 0 ? minutes + 'm ' : '') +
    (secs > 0 ? secs + 's ' : '')
  ).trim()
}

/** Add a delay anywhere in any async function. */
export const delay = (ms: number): Promise<void> =>
  new Promise(resolve => setTimeout(resolve, ms))

/** Check if two strings are equal, case insensitive. */
export const iequals = (a: string, b: string) => a.toLowerCase() === b.toLowerCase()

/** Format bytes number into a human readable string. */
export function formatBytes(bytes: number): string {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let unitIndex = 0

  while (bytes >= 1024 && unitIndex < units.length - 1) {
    bytes /= 1024
    unitIndex++
  }

  return `${bytes.toFixed(2)} ${units[unitIndex]}`
}

/** Check if two objects are equal. */
export const deepEqual = (obj1: any, obj2: any): boolean => {
  if (
    typeof obj1 !== 'object' ||
    typeof obj2 !== 'object' ||
    obj1 === null ||
    obj2 === null
  ) {
    return obj1 === obj2
  }

  const keys1 = Object.keys(obj1)
  const keys2 = Object.keys(obj2)
  if (keys1.length !== keys2.length) {
    return false
  }

  for (let key of keys1) {
    if (!obj2.hasOwnProperty(key) || !deepEqual(obj1[key], obj2[key])) {
      return false
    }
  }

  return true
}

/** Deep copy an object. */
export const deepCopy = <T>(obj: T): T => {
  if (typeof obj !== 'object' || obj === null) {
    return obj
  }

  if (Array.isArray(obj)) {
    return obj.map(item => deepCopy(item)) as any
  }

  const copiedObj: { [key in keyof T]?: T[key] } = {}
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      copiedObj[key] = deepCopy(obj[key])
    }
  }

  return copiedObj as T
}

export const maxLength = (str: string, max: number) =>
  str.length > max ? str.slice(0, max) + ' ....' : str
