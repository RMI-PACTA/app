export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return 'Empty'
  const sizes = ['Bytes', 'kB', 'MB', 'GB', 'TB']

  const k = 1000
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
