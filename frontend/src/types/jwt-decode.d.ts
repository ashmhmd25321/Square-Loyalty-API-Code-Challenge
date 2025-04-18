declare module 'jwt-decode' {
  export function jwtDecode<T = any>(token: string): T;
  export class InvalidTokenError extends Error {
    constructor(message: string);
  }
} 