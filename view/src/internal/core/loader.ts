/* eslint-disable @typescript-eslint/no-explicit-any */

import { sleep } from "@own-js-org/time"
import { Completer } from "./completer"

export enum PageStatus {
  None,
  Loading,
  Error,
  Ok,
}

/**
 * Stores and manages asynchronously loaded values, typically used for Vue page initialization data.
 * Caches successful results to prevent duplicate API calls, but allows retries on subsequent fetches if an error occurs.
 */
export class PageValue<T> {
  value?: T
  error?: unknown
  private completer_?: Completer<T>
  private fetch_: () => T | Promise<T>

  constructor(fetch: () => T | Promise<T>) {
    this.fetch_ = fetch
  }

  /**
   * Fetches the data. Returns the cached promise if already loading or succeeded.
   * If the previous attempt failed, it resets and retries the fetch.
   */
  async fetch(): Promise<T> {
    let c = this.completer_
    if (c) {
      return c.promise
    }

    c = new Completer<T>()
    this.completer_ = c

    try {
      const fetch = this.fetch_
      const val = await fetch()
      this.value = val
      c.resolve(val)
    } catch (e) {
      // Reset the completer on error so that the next fetch() call will trigger a retry.
      this.completer_ = undefined
      this.error = e
      c.reject(e)
    }

    return c.promise
  }
}
export interface IPageLoader {
  status: PageStatus
  error?: unknown
  fetch(): Promise<void>
}
/**
 * A parallel loader utility to initialize multiple page dependencies simultaneously.
 * Automatically manages the aggregate PageStatus for Vue UI states.
 */
export class PageLoader implements IPageLoader {
  // Expose the current status and the unified error for the Vue component to bind to.
  status: PageStatus = PageStatus.None
  error?: unknown

  constructor(public values: PageValue<any>[], public minLoadingDuration = 400) { }

  /**
   * Loads all values in parallel.
   * Automatically manages status transitions (Loading -> Ok / Error) and enables retry capability.
   */
  async fetch(): Promise<void> {
    switch (this.status) {
      case PageStatus.None:
      case PageStatus.Error:
        break
      default:
        return
    }
    if (this.values.length === 0) {
      this.status = PageStatus.Ok
      return
    }
    this.status = PageStatus.Loading

    const at = Date.now()
    let ret: any
    let hasError = false
    try {
      // Leverages Promise.all for standard fail-fast concurrent execution.
      await Promise.all(this.values.map(v => v.fetch()))
    } catch (e) {
      hasError = true
      ret = e
    }

    // Anti-flicker: Ensure the loading indicator stays visible for at least minLoadingDuration
    // to prevent the UI from blinking on ultra-fast networks.
    const diff = Date.now() - at
    if (diff < this.minLoadingDuration) {
      await sleep(this.minLoadingDuration - diff)
    }

    if (hasError) {
      this.status = PageStatus.Error
      this.error = ret
      throw ret
    } else {
      this.status = PageStatus.Ok
    }
  }
}
