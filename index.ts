function setupEvents() {
    document.querySelectorAll<HTMLElement>(".tag")?.forEach(t => {
        t.addEventListener("click", _ => {
            t.classList.toggle("active")
        })
    })

    document.querySelectorAll<HTMLAnchorElement>("a")?.forEach(a => {
        a.addEventListener("click", e => {
            if (!a.href.startsWith(window.location.origin)) return
            e.preventDefault()
            changePage(a.href)
        })
    })
}

function changePage(href: string) {
    fetch(href)
        .then(resp => resp.text())
        .then(html => {
            const doc = new DOMParser().parseFromString(html, "text/html")
            window.history.pushState({}, "", href)
            document.title = doc.title
            const distMain = doc.querySelector("main")
            const currentMain = document.querySelector("main")
            if (distMain === null || currentMain === null) document.body = doc.body
            else currentMain.innerHTML = distMain.innerHTML
            setupEvents()
        })
}

window.addEventListener("popstate", e => window.location.reload())
setupEvents()
