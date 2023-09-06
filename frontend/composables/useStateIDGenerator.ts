export const useStateIDGenerator = () => {
  const nextId = useState<number>('useStateIDGenerator', () => 1)
  return {
    id: (): string => {
      nextId.value = nextId.value + 1
      return `${nextId.value}`
    }
  }
}
