const navLinks = document.querySelectorAll("nav a");
for (let i = 0; i < navLinks.length; i++) {
    const link = navLinks[i];
    if (link.getAttribute('href') === window.location.pathname) {
        link.classList.add("live");
        break;
    }
}

function confirmDelete(feedID) {
    if (confirm("Are you sure you want to delete this feed?")) {
        window.location.href = '/feed/delete?id=' + feedID;
    }
}