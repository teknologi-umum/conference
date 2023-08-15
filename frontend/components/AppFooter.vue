<template>
    <footer class="py-8">
        <div class="container mx-auto px-5">
            <div class="grid grid-cols-1 lg:grid-cols-3">
                <div class="footer-menu footer-logo">
                    <img src="@/assets/images/logo_copy.png" alt="TeknumConf logo" class="h-15">
                </div>
                <div class="footer-menu">
                    <h4 class="footer-menu-header">Menu</h4>
                    <ul>
                        <li>
                            <router-link to="code-of-conduct">Code of Conduct</router-link>
                        </li>
                        <li>
                            <router-link to="code-of-conduct">Cookie Policy</router-link>
                        </li>
                        <li>
                            <router-link to="code-of-conduct">Privacy Policy</router-link>
                        </li>
                    </ul>
                </div>
                <div class="footer-menu social-media">
                    <h4 class="footer-menu-header">Our Social Media</h4>
                    <ul>
                        <li>
                            <router-link to="https://github.com/teknologi-umum/" >
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M7.8 2h8.4C19.4 2 22 4.6 22 7.8v8.4a5.8 5.8 0 0 1-5.8 5.8H7.8C4.6 22 2 19.4 2 16.2V7.8A5.8 5.8 0 0 1 7.8 2m-.2 2A3.6 3.6 0 0 0 4 7.6v8.8C4 18.39 5.61 20 7.6 20h8.8a3.6 3.6 0 0 0 3.6-3.6V7.6C20 5.61 18.39 4 16.4 4H7.6m9.65 1.5a1.25 1.25 0 0 1 1.25 1.25A1.25 1.25 0 0 1 17.25 8A1.25 1.25 0 0 1 16 6.75a1.25 1.25 0 0 1 1.25-1.25M12 7a5 5 0 0 1 5 5a5 5 0 0 1-5 5a5 5 0 0 1-5-5a5 5 0 0 1 5-5m0 2a3 3 0 0 0-3 3a3 3 0 0 0 3 3a3 3 0 0 0 3-3a3 3 0 0 0-3-3Z"/></svg>
                                <span>@teknum</span>
                            </router-link>
                        </li>
                        <li>
                            <router-link to="https://github.com/teknologi-umum/" >
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M9.78 18.65l.28-4.23l7.68-6.92c.34-.31-.07-.46-.52-.19L7.74 13.3L3.64 12c-.88-.25-.89-.86.2-1.3l15.97-6.16c.73-.33 1.43.18 1.15 1.3l-2.72 12.81c-.19.91-.74 1.13-1.5.71L12.6 16.3l-1.99 1.93c-.23.23-.42.42-.83.42z" fill="currentColor"/></svg>
                                <span>@teknologiumum</span>
                            </router-link>
                        </li>
                        <li>
                            <router-link to="https://github.com/teknologi-umum/" >
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M12 2A10 10 0 0 0 2 12c0 4.42 2.87 8.17 6.84 9.5c.5.08.66-.23.66-.5v-1.69c-2.77.6-3.36-1.34-3.36-1.34c-.46-1.16-1.11-1.47-1.11-1.47c-.91-.62.07-.6.07-.6c1 .07 1.53 1.03 1.53 1.03c.87 1.52 2.34 1.07 2.91.83c.09-.65.35-1.09.63-1.34c-2.22-.25-4.55-1.11-4.55-4.92c0-1.11.38-2 1.03-2.71c-.1-.25-.45-1.29.1-2.64c0 0 .84-.27 2.75 1.02c.79-.22 1.65-.33 2.5-.33c.85 0 1.71.11 2.5.33c1.91-1.29 2.75-1.02 2.75-1.02c.55 1.35.2 2.39.1 2.64c.65.71 1.03 1.6 1.03 2.71c0 3.82-2.34 4.66-4.57 4.91c.36.31.69.92.69 1.85V21c0 .27.16.59.67.5C19.14 20.16 22 16.42 22 12A10 10 0 0 0 12 2Z"/></svg>
                                <span>teknologiumum</span>
                            </router-link>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </footer>
</template>
<style>
footer {
    border-top: 1px solid white;
}
.footer-menu {
    margin-bottom: 1rem;
}
.footer-menu-header {
    font-size: 1.3rem;
    font-weight: bold;
    margin-bottom: 2rem;
}
.footer-menu ul li {
    margin-bottom: 1rem;
}
.footer-menu ul li a {
    display: block;
}
.footer-menu.social-media a {
    display: flex;
    align-items: center;
    gap: .5rem;
}
</style>