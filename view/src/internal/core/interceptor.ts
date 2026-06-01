/**
 * 爲攔截器鏈定義最終處理函數
 */
export type Done<Request, Response> = (req: Request) => Response

/**
 * 定義攔截器
 */
export type Interceptor<Request, Response> = (req: Request, next: (req: Request) => Response) => Response

/**
 * 執行攔截器鏈
 * @param {Array} interceptors - 攔截器陣列
 * @param {Function} done - 處理鏈的最終函式
 * @param {any} initialRequest - 初始的請求物件
 */
export function executeInterceptors<Request, Response>(
  interceptors: Array<Interceptor<Request, Response>>,
  done: Done<Request, Response>,
  initialRequest: Request) {
  let index = 0

  const next = (req: Request): Response => {
    if (index < interceptors.length) {
      const currentInterceptor = interceptors[index++]!
      // 遞迴呼叫下一個攔截器
      return currentInterceptor(req, next)
    } else {
      // 如果沒有更多攔截器，則呼叫最終處理函式
      return done(req)
    }
  }
  return next(initialRequest)
}
/**
 * 執行巢狀攔截器鏈
 * @param {Array} interceptors - 巢狀攔截器鏈
 * @param {Function} done - 處理鏈的最終異步函式
 * @param {any} initialRequest - 初始的請求物件
 */
export function executeInterceptorsArray<Request, Response>(
  interceptors: Array<Array<Interceptor<Request, Response>>>,
  done: Done<Request, Response>,
  initialRequest: Request): Response {

  let outerIndex = 0
  let innerIndex = 0

  const next = (req: Request): Response => {
    // 遍歷內部陣列
    if (innerIndex < interceptors[outerIndex]!.length) {
      const currentInterceptor = interceptors[outerIndex]![innerIndex++]!
      return currentInterceptor(req, next)
    }
    // 遍歷外部陣列
    else if (++outerIndex < interceptors.length) {
      // 重置內部索引，進入下一個內部陣列
      innerIndex = 0
      return next(req)
    }
    // 所有攔截器都已執行完畢
    else {
      return done(req)
    }
  }

  return next(initialRequest)
}
