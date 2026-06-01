/* eslint-disable @typescript-eslint/no-explicit-any */

export class CancelError extends Error { }
export class TimeoutError extends Error { }
export interface AbortContextOptions {
  /**
   * 關聯的 signal, 一旦 signal.aborted 變爲真將自動爲 Context 執行 abort(signal.reason)
   */
  signal?: AbortSignal | null

  /**
   * 超時毫秒數，當達到超時時間將自動爲 Context 執行 abort(new TimeoutError("timeouted"))
   */
  timeout?: number | null
}
export class AbortContext implements AbortController {
  private readonly ctrl: AbortController
  private listener?: () => void
  private parent?: AbortSignal
  private timer?: number
  constructor(opts?: AbortContextOptions) {
    if (opts) {
      const ctrl = new AbortController()
      const signal = opts.signal
      const timeout = opts.timeout
      let timer: any
      let listener: undefined | (() => void)
      try {
        if (signal) {
          if (signal.aborted) {
            ctrl.abort(signal.reason)
            this.ctrl = ctrl
            return
          } else {
            listener = () => {
              this._abort(signal.reason, true)
            }
            this.parent = signal
            this.listener = listener
            signal.addEventListener('abort', listener)
          }
        }
        if (timeout && timeout > 0) {
          timer = setTimeout(() => {
            this._abort(new TimeoutError("timeouted"), true)
          }, timeout)
          this.timer = timer
        }
        this.ctrl = ctrl
      } catch (e) {
        if (listener) {
          signal!.removeEventListener('abort', listener)
        }
        throw e
      }
    } else {
      this.ctrl = new AbortController()
    }
  }
  get signal(): AbortSignal {
    return this.ctrl.signal
  }
  abort(reason?: any): void {
    this._abort(reason)
  }
  private _abort(reason?: any, internal?: boolean): void {
    const timer = this.timer
    if (timer) {
      this.timer = undefined
      clearTimeout(timer)
    }
    const listener = this.listener
    if (listener) {
      const parent = this.parent
      if (parent) {
        this.listener = undefined
        this.parent = undefined
        parent.removeEventListener('abort', listener)
      } else {
        this.listener = undefined
      }
    }
    if (!this.ctrl.signal.aborted) {
      if (internal) {
        this.ctrl.abort(reason)
      }
      else if (reason === undefined || reason === null) {
        this.ctrl.abort(new CancelError("canceled"))
      } else {
        this.ctrl.abort(reason)
      }
    }
  }
}
