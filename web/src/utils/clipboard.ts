import { copyLink as apiCopyLink } from '@/api'

export function copyLink(
  path: string,
  httpAuthOn: boolean,
  httpAuthUser: string,
  httpAuthPass: string
): string {
  return apiCopyLink(path, httpAuthOn, httpAuthUser, httpAuthPass)
}
