import { HttpClient } from "./core/http_client";
import { Second } from "@own-js-org/time";

export const DefaultHttpClient = new HttpClient({
  timeout: Second * 30,
})
