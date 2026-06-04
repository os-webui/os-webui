export function errorString(e: unknown) {
  return e instanceof Error ? e.message : String(e)
}
