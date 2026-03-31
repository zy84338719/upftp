import { MCP_KEY_LENGTH, MCP_KEY_CHARS } from '@/constants'

export function generateMCPKey(): string {
  let key = ''
  for (let i = 0; i < MCP_KEY_LENGTH; i++) {
    key += MCP_KEY_CHARS.charAt(Math.floor(Math.random() * MCP_KEY_CHARS.length))
  }
  return key
}
