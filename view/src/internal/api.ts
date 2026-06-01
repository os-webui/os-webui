import { useLocaleStore } from "@/stores/locale";
import { HttpClient } from "./core/http_client";
import { Second } from "@own-js-org/time";

export const DefaultHttpClient = new HttpClient({
  timeout: Second * 30,
  interceptor: [
    (req, next) => {
      const choose = useLocaleStore().choose
      if (!choose || choose === '' || choose === 'auto') {
        return next(req)
      }
      const languages: string[] = []
      languages.push(choose)
      if (navigator.languages && navigator.languages.length > 0) {
        languages.push(...navigator.languages)
      }
      if (languages.length > 0) {
        const uniqueLanguages = Array.from(new Set(languages))
        req.headers.set('Accept-Language', uniqueLanguages.join(','))

        // const acceptLanguageVal = uniqueLanguages
        //   .map((lang, index) => {
        //     if (index === 0) return lang;
        //     const q = (1 - index * 0.1).toFixed(1);
        //     return `${lang};q=${Math.max(0.1, parseFloat(q))}`;
        //   })
        //   .join(',');
        // req.headers.set('Accept-Language', acceptLanguageVal);
      }
      return next(req)
    },
  ],
})
