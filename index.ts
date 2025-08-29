function setupEvents() {
    document.querySelectorAll<HTMLElement>(".tag")?.forEach(t => {
        t.addEventListener("click", _ => {
            const link = t.getAttribute("data-href")
            if (link === null || link === "") return
            if (!link.startsWith(window.location.origin) && link.startsWith("https://")) window.open(link)
            else changePage(link)
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
            document.body = doc.body
            document.head.querySelectorAll("style").forEach(e => document.head.removeChild(e))
            doc.head.querySelectorAll("style").forEach(e => {
                document.head.appendChild(e)
            })
            setupEvents()
            scrollTo({
                top: 0,
                left: 0,
                behavior: "smooth"
            })
        })
}

window.addEventListener("popstate", _ => {
    window.location.reload()
    scrollTo({
        top: 0,
        left: 0,
        behavior: "smooth"
    })
})
setupEvents()
