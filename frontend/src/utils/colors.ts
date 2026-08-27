// 分类/系列名到色点的映射
const dotColors = ['dot-blue', 'dot-pink', 'dot-purple', 'dot-green', 'dot-amber', 'dot-rose']

export function getDotColor(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = ((hash << 5) - hash + name.charCodeAt(i)) | 0
  }
  return dotColors[Math.abs(hash) % dotColors.length]
}
