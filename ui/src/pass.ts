const userPassKey = "user-password";
const adminPassKey = "password";

export function saveAdminPass(pass: string) {
  localStorage.setItem(adminPassKey, pass);
}

export function getAdminPass(): string {
  return localStorage.getItem(adminPassKey) ?? "";
}

export function saveUserPass(pass: string) {
  localStorage.setItem(userPassKey, pass);
}

export function getUserPass(): string {
  return localStorage.getItem(userPassKey) ?? "";
}
