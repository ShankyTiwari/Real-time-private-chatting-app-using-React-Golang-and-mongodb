export function setItemInLS(key, value) {
    localStorage.setItem(
        key,
        JSON.stringify(value)
    )
}

export function getItemInLS(key) {
    const lSValue = localStorage.getItem(
        key,
    )
    if(lSValue) {
        return JSON.parse(lSValue)
    } else {
        return null;
    }
}

export function removeItemInLS(key) {
    localStorage.removeItem(
        key,
    )
}