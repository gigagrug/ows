class Shell extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: "open" })
    this.shadowRoot.innerHTML = `
			<link href="./assets/style.css" rel="stylesheet" />

			<div id="wrap">
				<nav class="bg-white border-b border-neutral-200 px-4 py-2.5 dark:bg-neutral-800 dark:border-neutral-700 fixed left-0 right-0 top-0 z-50">
					<div class="flex flex-wrap justify-between items-center">
						<div class="flex justify-start items-center">
							<button id="openDrawer" aria-controls="drawer-navigation" class="p-2 mr-2 text-neutral-600 rounded-lg cursor-pointer md:hidden hover:text-neutral-900 hover:bg-neutral-100 focus:bg-neutral-100 dark:focus:bg-neutral-700 focus:ring-2 focus:ring-neutral-100 dark:focus:ring-neutral-700 dark:text-neutral-400 dark:hover:bg-neutral-700 dark:hover:text-white">
								<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
									<path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h6a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd"></path>
								</svg>
								<svg class="hidden w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
									<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path>
								</svg>
								<span class="sr-only">Toggle sidebar</span>
							</button>
							<a href="/" class="c-text ml-3 self-center text-2xl font-bold whitespace-nowrap">OWS</a>
						</div>
						<div class="flex">
							<label for="darkModeToggle" class="flex items-center cursor-pointer">
								<input type="checkbox" id="darkModeToggle" />
								<p class="ms-1">Dark</p>
							</label>
						</div>
					</div>
				</nav>

				<!-- Sidebar -->
				<aside id="sideDrawer" class="fixed top-0 left-0 z-40 w-52 h-screen pt-14 transition-transform -translate-x-full bg-white border-r border-neutral-200 md:translate-x-0 dark:bg-neutral-800 dark:border-neutral-700" aria-label="Sidenav" id="drawer-navigation">
					<div class="overflow-y-auto py-5 px-3 h-full">
						<a href="/" class="flex items-center p-2 font-medium c-text rounded-lg hover:bg-neutral-100 dark:hover:bg-neutral-700">
							<span class="ml-3">Home</span>
						</a>
						<a href="/logs" class="flex items-center p-2 font-medium c-text rounded-lg hover:bg-neutral-100 dark:hover:bg-neutral-700">
							<span class="ml-3">Logs</span>
						</a>
						<hr class="c-border border-0 border-t"/>
						<a href="/account" class="flex items-center p-2 font-medium c-text rounded-lg hover:bg-neutral-100 dark:hover:bg-neutral-700">
							<span class="ml-3">Account</span>
						</a>				
					</div>
				</aside>

				<main class="p-4 md:ml-52 h-auto pt-20">
					<slot></slot>
				</main>
			</div>    
		`
  }

  connectedCallback() {
    console.log(window.location.pathname)
    const openDrawerBtn = this.shadowRoot.getElementById("openDrawer")
    const sideDrawer = this.shadowRoot.getElementById("sideDrawer")
    openDrawerBtn.addEventListener("click", function () {
      sideDrawer.classList.toggle("-translate-x-full")
    })

    const isDarkMode = localStorage.getItem("darkMode")
    const wrap = this.shadowRoot.getElementById("wrap")
    const htmlTag = document.documentElement
    const toggleSwitch = this.shadowRoot.querySelector("#darkModeToggle")
    if (isDarkMode === "true") {
      wrap.classList.add("dark")
      htmlTag.classList.add("dark")
      if (toggleSwitch) {
        toggleSwitch.checked = true
      }
    }
    if (toggleSwitch) {
      toggleSwitch.addEventListener("change", function () {
        if (this.checked) {
          wrap.classList.add("dark")
          htmlTag.classList.add("dark")
          localStorage.setItem("darkMode", "true")
        } else {
          wrap.classList.remove("dark")
          htmlTag.classList.remove("dark")
          localStorage.setItem("darkMode", "false")
        }
      })
    }
  }
}

customElements.define("app-shell", Shell)
