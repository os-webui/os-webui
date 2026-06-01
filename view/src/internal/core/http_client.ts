/* eslint-disable @typescript-eslint/no-explicit-any */
import { AbortContext } from "./context.ts"
import { executeInterceptors, executeInterceptorsArray, type Interceptor } from './interceptor.ts'
/**
 * 響應 body 的格式
 */
export enum ResponseType {
  json = 1,
  text = 2,
  arrayBuffer = 3,
  blob = 4,
  bytes = 5,
}
interface ResponseTypeHook {
  set?: (req: Request) => void
  get: (resp: Response) => Promise<any>
}
const responseTypeHookJSON: ResponseTypeHook = {
  set(req) {
    req.headers.set('accept', 'application/json')
  },
  get(resp) {
    return resp.json()
  },
}
const responseTypeHookText: ResponseTypeHook = {
  get(resp) {
    return resp.text()
  },
}
const responseTypeHookArrayBuffer: ResponseTypeHook = {
  get(resp) {
    return resp.arrayBuffer()
  },
}
const responseTypeHookBlob: ResponseTypeHook = {
  get(resp) {
    return resp.blob()
  },
}
const responseTypeHookBytes: ResponseTypeHook = {
  get(resp) {
    return resp.bytes()
  },
}

export interface Options<T extends ResponseType> {
  /**
   * 返回數據類型
   */
  type?: T
  /**
   * 超時時間，如果 <=0 則不設置超時
   */
  timeout?: number

  /**
   * 請求攔截器，在默認攔截器之前調用
   */
  beforeInterceptor?: Array<Interceptor<Request, Promise<Response>>>
  /**
   * 請求攔截器，會覆蓋默認的攔截器
   */
  interceptor?: Array<Interceptor<Request, Promise<Response>>>
  /**
   * 請求攔截器，在默認攔截器之後調用
   */
  afterInterceptor?: Array<Interceptor<Request, Promise<Response>>>
}
/**
 * 請求拋出的異常
 */
export class HttpError extends Error {
  /**
   * 原始響應如果爲 undefined 則表示 fetch 拋出了異常，否則表示後續解析數據出現了異常
   */
  readonly response?: Response
  /**
   * 拋出的原始異常
   */
  readonly error: unknown
  constructor(error: unknown, response?: Response) {
    if (error === undefined || error === null) {
      super(`${response?.statusText}`)
    }
    else if (error instanceof Error) {
      super(error.message)
    } else {
      super(`${error}`)
    }
    this.response = response
    this.error = error
  }
}
/**
 * http 請求成功的響應
 */
export interface HttpResponse<Data, T extends ResponseType> {
  /**
   * 錯誤代碼
   */
  status: number
  /**
   * 解析後的 body
   */
  data: T extends ResponseType.text ? string :
  (T extends ResponseType.arrayBuffer ? ArrayBuffer :
    (T extends ResponseType.blob ? Blob :
      (T extends ResponseType.bytes ? Uint8Array : Data)
    )
  )
  /**
   * 原始響應
   */
  response: Response
}
function getResponseTypeHook(opts?: Options<any>): ResponseTypeHook {
  if (opts) {
    const t = opts.type
    if (t === undefined || t === null) {
      return responseTypeHookJSON
    }
    switch (t) {
      case ResponseType.json:
        return responseTypeHookJSON
      case ResponseType.text:
        return responseTypeHookText
      case ResponseType.arrayBuffer:
        return responseTypeHookArrayBuffer
      case ResponseType.blob:
        return responseTypeHookBlob
      case ResponseType.bytes:
        return responseTypeHookBytes
      default:
        throw new Error(`unknow response type: ${t}`)
    }
  }
  return responseTypeHookJSON
}
export interface HttpClientOptions {
  /**
    * 每个请求默认的超時時間，如果 <=0 則不設置超時
    */
  timeout?: number
  /**
   * 默認的攔截器
   */
  interceptor?: Array<(req: Request, next: (req: Request) => Promise<Response>) => Promise<Response>>
}
/**
 * 使用 fetch 作爲底層的 http client，提供了一些方便的功能
 */
export class HttpClient {
  constructor(private readonly opts?: HttpClientOptions) { }
  get<Data>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<ResponseType.json>): Promise<HttpResponse<Data, ResponseType.json>>
  get(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.text }): Promise<HttpResponse<string, ResponseType.text>>
  get(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.arrayBuffer }): Promise<HttpResponse<ArrayBuffer, ResponseType.arrayBuffer>>
  get(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.blob }): Promise<HttpResponse<Blob, ResponseType.blob>>
  get(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.bytes }): Promise<HttpResponse<Uint8Array, ResponseType.bytes>>
  get<Data, T extends ResponseType>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<T>) {
    init?.signal?.throwIfAborted()
    return this._fetch<Data, T>(
      input, getResponseTypeHook(opts),
      init, 'GET',
      opts)
  }
  post<Data>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<ResponseType.json>): Promise<HttpResponse<Data, ResponseType.json>>
  post(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.text }): Promise<HttpResponse<string, ResponseType.text>>
  post(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.arrayBuffer }): Promise<HttpResponse<ArrayBuffer, ResponseType.arrayBuffer>>
  post(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.blob }): Promise<HttpResponse<Blob, ResponseType.blob>>
  post(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.bytes }): Promise<HttpResponse<Uint8Array, ResponseType.bytes>>
  post<Data, T extends ResponseType>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<T>) {
    init?.signal?.throwIfAborted()
    return this._fetch<Data, T>(
      input, getResponseTypeHook(opts),
      init, 'POST',
      opts)
  }
  put<Data>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<ResponseType.json>): Promise<HttpResponse<Data, ResponseType.json>>
  put(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.text }): Promise<HttpResponse<string, ResponseType.text>>
  put(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.arrayBuffer }): Promise<HttpResponse<ArrayBuffer, ResponseType.arrayBuffer>>
  put(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.blob }): Promise<HttpResponse<Blob, ResponseType.blob>>
  put(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.bytes }): Promise<HttpResponse<Uint8Array, ResponseType.bytes>>
  put<Data, T extends ResponseType>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<T>) {
    init?.signal?.throwIfAborted()
    return this._fetch<Data, T>(
      input, getResponseTypeHook(opts),
      init, 'PUT',
      opts)
  }
  patch<Data>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<ResponseType.json>): Promise<HttpResponse<Data, ResponseType.json>>
  patch(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.text }): Promise<HttpResponse<string, ResponseType.text>>
  patch(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.arrayBuffer }): Promise<HttpResponse<ArrayBuffer, ResponseType.arrayBuffer>>
  patch(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.blob }): Promise<HttpResponse<Blob, ResponseType.blob>>
  patch(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.bytes }): Promise<HttpResponse<Uint8Array, ResponseType.bytes>>
  patch<Data, T extends ResponseType>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<T>) {
    init?.signal?.throwIfAborted()
    return this._fetch<Data, T>(
      input, getResponseTypeHook(opts),
      init, 'PATCH',
      opts)
  }
  delete<Data>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<ResponseType.json>): Promise<HttpResponse<Data, ResponseType.json>>
  delete(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.text }): Promise<HttpResponse<string, ResponseType.text>>
  delete(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.arrayBuffer }): Promise<HttpResponse<ArrayBuffer, ResponseType.arrayBuffer>>
  delete(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.blob }): Promise<HttpResponse<Blob, ResponseType.blob>>
  delete(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.bytes }): Promise<HttpResponse<Uint8Array, ResponseType.bytes>>
  delete<Data, T extends ResponseType>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<T>) {
    init?.signal?.throwIfAborted()
    return this._fetch<Data, T>(
      input, getResponseTypeHook(opts),
      init, 'DELETE',
      opts)
  }
  do<Data>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<ResponseType.json>): Promise<HttpResponse<Data, ResponseType.json>>
  do(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.text }): Promise<HttpResponse<string, ResponseType.text>>
  do(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.arrayBuffer }): Promise<HttpResponse<ArrayBuffer, ResponseType.arrayBuffer>>
  do(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.blob }): Promise<HttpResponse<Blob, ResponseType.blob>>
  do(input: RequestInfo | URL, init?: RequestInit, opts?: Omit<Options<any>, 'type'> & { type: ResponseType.bytes }): Promise<HttpResponse<Uint8Array, ResponseType.bytes>>
  do<Data, T extends ResponseType>(input: RequestInfo | URL, init?: RequestInit, opts?: Options<T>) {
    init?.signal?.throwIfAborted()
    return this._fetch<Data, T>(
      input, getResponseTypeHook(opts),
      init, undefined,
      opts)
  }
  /**
   * 發送 fetch 請求
   */
  private async _fetch<Data, T extends ResponseType>(
    input: RequestInfo | URL,
    hook: ResponseTypeHook,
    init?: RequestInit,
    method?: string,
    opts?: Options<T>): Promise<HttpResponse<Data, T>> {
    if (method && init?.method != method) {
      if (init) {
        init = {
          ...init,
          method: method,
        }
      } else {
        init = {
          method: method,
        }
      }
    }
    // 創建請求
    const req = new Request(input, init)
    switch (req.method) {
      case "POST":
      case "PUT":
      case "PATCH":
        {
          const t = req.headers.get('content-type')
          if (!t || t.startsWith('text/plain')) {
            req.headers.set('content-type', 'application/json;charset=utf-8')
          }
        }
        break
    }

    let interceptor: Array<Array<Interceptor<Request, Promise<Response>>>> | undefined
    if (opts) {
      let i = 0
      if (opts.beforeInterceptor?.length) {
        ++i
      }
      if (this.opts?.interceptor?.length || opts.interceptor?.length) {
        ++i
      }
      if (opts.afterInterceptor?.length) {
        ++i
      }
      interceptor = new Array<Array<Interceptor<Request, Promise<Response>>>>(i)
      i = 0
      if (opts.beforeInterceptor?.length) {
        interceptor[i++] = opts.beforeInterceptor
      }
      if (opts.interceptor?.length) {
        interceptor[i++] = opts.interceptor
      } else if (this.opts?.interceptor?.length) {
        interceptor[i++] = this.opts.interceptor
      }
      if (opts.afterInterceptor?.length) {
        interceptor[i++] = opts.afterInterceptor
      }
    }
    let resp: undefined | Response
    try {
      if (interceptor && interceptor.length) {
        resp = await executeInterceptorsArray(interceptor,
          (req: Request) => {
            return this._fetchTimeout(req, hook, opts)
          },
          req)
      } else {
        const interceptor = this.opts?.interceptor
        if (interceptor && interceptor.length) {
          resp = await executeInterceptors(interceptor,
            (req: Request) => {
              return this._fetchTimeout(req, hook, opts)
            },
            req)
        } else {
          resp = await this._fetchTimeout(req, hook, opts)
        }
      }
      const status = resp.status
      if (status < 200 || status > 299) {
        const text = await resp.clone().text()
        throw new Error(`http ${status}: ${text}`)
      }
      const data = await hook.get(resp)
      return {
        status: status,
        data: data,
        response: resp,
      }
    } catch (e) {
      if (e instanceof HttpError) {
        throw e
      } else {
        throw new HttpError(e, resp)
      }
    }
  }
  /**
   * 爲 fetch 添加上超時
   */
  private _fetchTimeout(
    req: Request,
    hook: ResponseTypeHook,
    opts?: Options<any>): Promise<Response> {
    let init: RequestInit | undefined
    // 設置超時
    const timeout = opts?.timeout ?? (this.opts?.timeout) ?? 0
    if (timeout > 0) {
      init = {
        signal: new AbortContext({
          signal: req.signal,
          timeout: timeout,
        }).signal,
      }
    }
    if (hook.set) {
      hook.set(req)
    }
    return fetch(req, init)
  }
}
export const httpClient = new HttpClient({
  timeout: 1000 * 30,
})
export function getError(e: any): string {
  if (typeof e === "object") {
    if (typeof e.message === "string") {
      return e.message
    }
  }
  return `${e}`
}
