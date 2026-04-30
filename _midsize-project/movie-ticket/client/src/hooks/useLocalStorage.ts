import { useState, useEffect } from 'react';

export function useLocalStorage(key: string, initialValue: string): [string, (value: string) => void] {
  const [value, setValue] = useState<string>(() => {
    const storedValue = localStorage.getItem(key);
    return storedValue !== null ? storedValue : initialValue;
  });

  useEffect(() => {
    localStorage.setItem(key, value);
  }, [key, value]);

  return [value, setValue];
}
