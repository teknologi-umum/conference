<template>
    <footer class="py-8">
        <div class="container mx-auto px-5">
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-5 lg:gap-0">
                <div class="footer-menu footer-logo">
                    <img src="@/assets/images/logo_copy.png" alt="TeknumConf logo" class="h-15">
                </div>
                <div class="footer-menu">
                    <h4 class="footer-menu-header">Menu</h4>
                    <ul class="list-none">
                        <li>
                            <router-link class="white" to="code-of-conduct">Code of Conduct</router-link>
                        </li>
                        <li>
                            <router-link class="white" to="cookie-policy">Cookie Policy</router-link>
                        </li>
                        <li>
                            <router-link class="white" to="privacy-policy">Privacy Policy</router-link>
                        </li>
                    </ul>
                </div>
                <div class="footer-menu social-media">
                    <h4 class="footer-menu-header">Our Social Media</h4>
                    <ul class="list-none">
                        <li>
                            <a data-umami-event="social-media-telegram" class="white" href="https://t.me/teknologi_umum_v2" target="_blank">
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M9.78 18.65l.28-4.23l7.68-6.92c.34-.31-.07-.46-.52-.19L7.74 13.3L3.64 12c-.88-.25-.89-.86.2-1.3l15.97-6.16c.73-.33 1.43.18 1.15 1.3l-2.72 12.81c-.19.91-.74 1.13-1.5.71L12.6 16.3l-1.99 1.93c-.23.23-.42.42-.83.42z" fill="currentColor"/></svg>
                                <span>@teknologi_umum_v2</span>
                            </a>
                        </li>
                        <li>
                            <a data-umami-event="social-media-github" class="white" href="https://github.com/teknologi-umum/" target="_blank">
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M12 2A10 10 0 0 0 2 12c0 4.42 2.87 8.17 6.84 9.5c.5.08.66-.23.66-.5v-1.69c-2.77.6-3.36-1.34-3.36-1.34c-.46-1.16-1.11-1.47-1.11-1.47c-.91-.62.07-.6.07-.6c1 .07 1.53 1.03 1.53 1.03c.87 1.52 2.34 1.07 2.91.83c.09-.65.35-1.09.63-1.34c-2.22-.25-4.55-1.11-4.55-4.92c0-1.11.38-2 1.03-2.71c-.1-.25-.45-1.29.1-2.64c0 0 .84-.27 2.75 1.02c.79-.22 1.65-.33 2.5-.33c.85 0 1.71.11 2.5.33c1.91-1.29 2.75-1.02 2.75-1.02c.55 1.35.2 2.39.1 2.64c.65.71 1.03 1.6 1.03 2.71c0 3.82-2.34 4.66-4.57 4.91c.36.31.69.92.69 1.85V21c0 .27.16.59.67.5C19.14 20.16 22 16.42 22 12A10 10 0 0 0 12 2Z"/></svg>
                                <span>teknologi-umum</span>
                            </a>
                        </li>
                        <li>
                            <a data-umami-event="social-media-twitter" class="white" href="https://twitter.com/teknumconf/" target="_blank">
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                  <line x1="18" y1="6" x2="6" y2="18"></line>
                                  <line x1="6" y1="6" x2="18" y2="18"></line>
                                </svg>
                                <span>@teknumconf</span>
                            </a>
                        </li>
                        <li>
                            <a data-umami-event="social-media-youtube" class="white" href="https://www.youtube.com/@teknologiumum1227" target="_blank">
                              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M22.54 6.42a2.78 2.78 0 0 0-1.94-2C18.88 4 12 4 12 4s-6.88 0-8.6.46a2.78 2.78 0 0 0-1.94 2A29 29 0 0 0 1 11.75a29 29 0 0 0 .46 5.33A2.78 2.78 0 0 0 3.4 19a53.5 53.5 0 0 0 8.6.53 53.5 53.5 0 0 0 8.6-.53 2.78 2.78 0 0 0 1.94-2 29 29 0 0 0 .46-5.33 29 29 0 0 0-.46-5.33z"></path>
                                <polygon points="9.75 15.02 15.5 11.75 9.75 8.48 9.75 15.02"></polygon>
                              </svg>
                              <span>@teknologiumum</span>
                            </a>
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