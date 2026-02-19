const storageKey = "password"

export function savePass(pass: string) {
    localStorage.setItem(storageKey, pass)
}

export function getPass(): string {
    return localStorage.getItem(storageKey) ?? ""
}