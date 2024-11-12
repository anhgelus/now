const tags = document.querySelectorAll<HTMLElement>(".tag")

tags?.forEach(t => {
    t.addEventListener("click", _ => {
        t.classList.toggle("active")
    })
})
