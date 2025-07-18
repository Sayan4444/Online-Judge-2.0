import Env from "./env";

export const BASE_URL = Env.BACKEND_URL;
export const LOGIN_URL = `${Env.BACKEND_URL}/login`;
export const ADMIN_LOGIN_URL = `${Env.BACKEND_URL}/admin/login`;
export const ADMIN_URL = `${Env.BACKEND_URL}/admin`;
export const API_URL = `${Env.BACKEND_URL}/api`;
