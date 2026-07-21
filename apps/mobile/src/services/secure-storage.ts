import * as SecureStore from 'expo-secure-store'

export function getSecureItem(key: string) {
  return SecureStore.getItemAsync(key)
}

export function setSecureItem(key: string, value: string) {
  return SecureStore.setItemAsync(key, value)
}

export function deleteSecureItem(key: string) {
  return SecureStore.deleteItemAsync(key)
}
