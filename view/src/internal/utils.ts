/**
 * 用於淺比較兩個數組是否相等
 */
export function equal<T>(a: readonly T[], b: readonly T[]) {
    if (a === b) {
        return true
    }
    if (a.length === b.length) {
        for (let i = 0; i < a.length; i++) {
            if (a[i] !== b[i]) {
                return false
            }
        }
        return true
    }
    return false
}
/**
 * 保證函數只會被調用1次
 */
export function once(f: () => void) {
    let ok = true
    return () => {
        if (ok) {
            ok = false
            f()
        }
    }
}

/**
 * 存儲設置到 localStorage
 * @param key 要存儲的鍵
 * @param val 要存儲的值
 * @param def 默認值，如果 val 和默認值相等則會 移除 key
 */
export function saveValue(key: string, val: string, def: string) {
    try {
        if (val === def) {
            localStorage.removeItem(key)
        } else {
            localStorage.setItem(key, val)
        }
        console.info(`save "${key}" to "${val}"`)
    } catch (e) {
        console.warn(`save "${key}" to "${val}" fail`, e)
    }
}