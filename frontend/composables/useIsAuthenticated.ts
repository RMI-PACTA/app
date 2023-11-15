export const useIsAuthenticated = () => {
  const prefix = 'useIsAuthenticated'
  const isAuthenticated = useState<boolean>(`${prefix}.isAuthenticated`, () => true)
  return isAuthenticated
}
