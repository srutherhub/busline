export class ApiClient {
  baseUrl: string = "/api";

  async request(path: string, options = {}) {
    const res = await fetch(`${this.baseUrl}${path}`, {
      headers: { "Content-Type": "application/json" },
      ...options,
    });
    if (!res.ok) {
      let body: ApiErrorBody;
      try {
        body = await res.json();
      } catch {
        body = {
          code: "UNKNOWN",
          message: res.statusText || "Request failed",
        };
      }
      throw new ApiError(res.status, body);
    }
    return res.status === 204 ? null : res.json();
  }

  get(path: string) {
    return this.request(path);
  }

  post(path: string, body = {}) {
    return this.request(path, { method: "POST", body: JSON.stringify(body) });
  }

  put(path: string, body = {}) {
    return this.request(path, { method: "PUT", body: JSON.stringify(body) });
  }

  delete(path: string) {
    return this.request(path, { method: "DELETE" });
  }
}

type ApiErrorBody = {
  code: string;
  message: string;
  details?: string;
};

export class ApiError extends Error {
  readonly code: string;
  readonly details?: string;
  readonly status: number;

  constructor(status: number, body: ApiErrorBody) {
    super(body.message);
    this.name = "ApiError";
    this.status = status;
    this.code = body.code;
    this.details = body.details;
  }
}
